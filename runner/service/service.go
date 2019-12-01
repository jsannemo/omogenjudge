package service

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/google/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
	"github.com/jsannemo/omogenjudge/runner/eval"
	"github.com/jsannemo/omogenjudge/runner/language"
	"github.com/jsannemo/omogenjudge/runner/runners"
	execpb "github.com/jsannemo/omogenjudge/sandbox/api"
	eclient "github.com/jsannemo/omogenjudge/sandbox/client"
	"github.com/jsannemo/omogenjudge/util/go/files"
	"github.com/jsannemo/omogenjudge/util/go/users"
)

var (
	address = flag.String("run_listen_addr", "127.0.0.1:61811", "The run server address to listen to in the format host:port")
)

type runServer struct {
	languages []*runpb.Language
	exec      execpb.ExecuteServiceClient
}

// TODO: rewrite
func (s *runServer) SimpleRun(ctx context.Context, req *runpb.SimpleRunRequest) (*runpb.SimpleRunResponse, error) {
	logger.Infof("RunService.SimpleRun: %v", req)
	execStream, err := s.exec.Execute(context.Background())
	defer execStream.CloseSend()
	if err != nil {
		return nil, err
	}

	lang, exists := language.GetLanguage(req.Program.LanguageId)
	if !exists {
		return nil, status.Errorf(codes.InvalidArgument, "Language %v does not exist", req.Program.LanguageId)
	}
	program, err := lang.Program(req.Program, execStream)
	if err != nil {
		return nil, err
	}
	dir, err := ioutil.TempDir("/var/lib/omogen/tmps", "simplerun")
	if err != nil {
		return nil, err
	}
	// defer os.RemoveAll(dir)
	fb := files.NewFileBase(dir)
	fb.Gid = users.OmogenClientsID()
	if err := fb.FixOwners("."); err != nil {
		return nil, err
	}
	if err := fb.FixModeExec("."); err != nil {
		return nil, err
	}
	fb.GroupWritable = true

	env, err := runners.NewFileLinker(filepath.Join(dir, "env"))
	if err != nil {
		return nil, err
	}
	program.SetArgs(&runners.ProgramArgs{
		InputPath:     env.PathFor("input", false),
		OutputPath:    env.PathFor("output", true),
		ErrorPath:     env.PathFor("error", true),
		ExtraArgs:     req.Arguments,
		TimeLimitMs:   5000,
		MemoryLimitKb: 1000 * 1000,
	})
	var results []*runpb.SimpleRunResult
	for i, input := range req.InputFiles {
		outName := fmt.Sprintf("output-%d", i)
		errName := fmt.Sprintf("error-%d", i)

		outPath := filepath.Join(dir, outName)
		if err := fb.WriteFile(outName, []byte{}); err != nil {
			return nil, err
		}
		errPath := filepath.Join(dir, errName)
		if err := fb.WriteFile(errName, []byte{}); err != nil {
			return nil, err
		}
		if err := env.LinkFile(input, "input", false); err != nil {
			return nil, err
		}
		if err := env.LinkFile(outPath, "output", true); err != nil {
			return nil, err
		}
		if err := env.LinkFile(errPath, "error", true); err != nil {
			return nil, err
		}
		res, err := program.Execute()
		if err != nil {
			return nil, err
		}
		errData, _ := ioutil.ReadFile(errPath)
		outData, _ := ioutil.ReadFile(outPath)
		results = append(results, &runpb.SimpleRunResult{
			ExitCode: res.ExitCode,
			Signal:   res.Signal,
			Timeout:  res.TimedOut(),
			Stderr:   string(errData),
			Stdout:   string(outData),
		})
		if err := env.Clear(); err != nil {
			return nil, err
		}
	}
	return &runpb.SimpleRunResponse{Results: results}, nil
}

// Implementation of RunServer.GetLanguages.
func (s *runServer) GetLanguages(ctx context.Context, _ *runpb.GetLanguagesRequest) (*runpb.GetLanguagesResponse, error) {
	return &runpb.GetLanguagesResponse{InstalledLanguages: s.languages}, nil
}

var compileCache = make(map[string]*runpb.CompileResponse)
var cacheLock sync.Mutex

func (s *runServer) CompileCached(ctx context.Context, req *runpb.CompileCachedRequest) (*runpb.CompileCachedResponse, error) {
	logger.Infof("/RunService.CompileCached: %v", req)
	if validation := ValidateCompileRequest(req.Request); validation.Code() != codes.OK {
		return nil, validation.Err()
	}
	cacheLock.Lock()
	defer cacheLock.Unlock()
	res, has := compileCache[req.Identifier]
	if has {
		return &runpb.CompileCachedResponse{Response: res}, nil
	} else {
		res, err := s.Compile(ctx, req.Request)
		if err != nil {
			return nil, err
		}
		compileCache[req.Identifier] = res
		return &runpb.CompileCachedResponse{Response: res}, nil
	}
}

// Implementation of RunServer.Compile.
func (s *runServer) Compile(ctx context.Context, req *runpb.CompileRequest) (*runpb.CompileResponse, error) {
	logger.Infof("/RunService.Compile: %v", req)
	if status := ValidateCompileRequest(req); status.Code() != codes.OK {
		return nil, status.Err()
	}

	language, exists := language.GetLanguage(req.Program.LanguageId)
	if !exists {
		return nil, status.Errorf(codes.InvalidArgument, "Language %v does not exist", req.Program.LanguageId)
	}
	// TODO: pass on context here
	compiledProgram, err := language.Compile(req.Program, req.OutputPath, s.exec)
	if err != nil {
		return nil, err
	}
	response := &runpb.CompileResponse{
		Program:           compiledProgram.Program,
		CompilationOutput: compiledProgram.Output,
		CompilationError:  compiledProgram.Errors,
	}
	return response, nil
}

// Implementation of RunServer.Evaluate.
func (s *runServer) Evaluate(req *runpb.EvaluateRequest, stream runpb.RunService_EvaluateServer) error {
	logger.Infof("/RunService.Evaluate: %v", req)
	if err := ValidateEvaluateRequest(req); err != nil {
		return err
	}

	execStream, err := s.exec.Execute(context.Background())
	defer execStream.CloseSend()
	if err != nil {
		return status.Errorf(codes.Internal, "Failed opening execution stream: %v", err)
	}

	lang, exists := language.GetLanguage(req.Program.LanguageId)
	if !exists {
		return status.Errorf(codes.InvalidArgument, "Language %v does not exist", req.Program.LanguageId)
	}
	program, err := lang.Program(req.Program, execStream)
	if err != nil {
		return status.Errorf(codes.Internal, "Invalid compiled program: %v", err)
	}
	var validator runners.Program
	if req.Validator != nil {
		valExecStream, err := s.exec.Execute(context.Background())
		defer valExecStream.CloseSend()
		if err != nil {
			return status.Errorf(codes.Internal, "Failed opening execution stream for validation: %v", err)
		}
		lang, exists := language.GetLanguage(req.Validator.Program.LanguageId)
		if !exists {
			return status.Errorf(codes.InvalidArgument, "Language %v does not exist", req.Program.LanguageId)
		}
		validator, err = lang.Program(req.Validator.Program, valExecStream)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "Invalid validator: %v", err)
		}
	}

	root := fmt.Sprintf("/var/lib/omogen/submissions/%s", req.SubmissionId)
	evaluator, err := eval.NewEvaluator(root, program, validator)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed creating evaluator: %v", err)
	}
	evaluator.EvaluateAll = req.EvaluateAll

	results := make(chan *eval.Result, 10)
	wg := sync.WaitGroup{}
	wg.Add(1)
	var evalErr error
	go func() {
		// TODO handle these errors
		for res := range results {
			if res.TestCaseVerdict != runpb.Verdict_VERDICT_UNSPECIFIED {
				if err = stream.Send(&runpb.EvaluateResponse{
					Result: &runpb.EvaluateResponse_TestCase{TestCase: &runpb.TestCaseResult{Verdict: res.TestCaseVerdict,
						TimeUsageMs: res.TimeUsageMs,
						Score:       res.Score,
					}},
				}); err != nil {
					break
				}
			} else if res.TestGroupVerdict != runpb.Verdict_VERDICT_UNSPECIFIED {
				if err = stream.Send(&runpb.EvaluateResponse{
					Result: &runpb.EvaluateResponse_TestGroup{TestGroup: &runpb.TestGroupResult{Verdict: res.TestGroupVerdict,
						TimeUsageMs: res.TimeUsageMs,
						Score:       res.Score,
					}},
				}); err != nil {
					break
				}
			} else if res.SubmissionVerdict != runpb.Verdict_VERDICT_UNSPECIFIED {
				if err = stream.Send(&runpb.EvaluateResponse{
					Result: &runpb.EvaluateResponse_Submission{Submission: &runpb.SubmissionResult{
						Verdict:     res.SubmissionVerdict,
						TimeUsageMs: res.TimeUsageMs,
						Score:       res.Score,
					}},
				}); err != nil {
					break
				}
			}
		}
		wg.Done()
	}()
	if err := evaluator.Evaluate(req.Groups, req.TimeLimitMs, req.MemLimitKb, results); err != nil {
		return status.Errorf(codes.Internal, "Failed evaluation: %v", err)
	}
	if evalErr != nil {
		return status.Errorf(codes.Internal, "Could not send back verdict: %v", evalErr)
	}
	wg.Wait()
	return nil
}

func newServer() (*runServer, error) {
	apiLanguages := make([]*runpb.Language, 0)
	for _, language := range language.GetLanguages() {
		apiLanguages = append(apiLanguages, language.ToApiLanguage())
	}
	execClient, err := eclient.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed creating ExecuteService client: %v", err)
	}
	s := &runServer{
		languages: apiLanguages,
		exec:      execClient,
	}
	return s, nil
}

// Register registers a new RunService with the given server.
func Register(grpcServer *grpc.Server) error {
	server, err := newServer()
	if err != nil {
		return err
	}
	runpb.RegisterRunServiceServer(grpcServer, server)
	return nil
}
