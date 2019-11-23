package main

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	"github.com/google/logger"
	"google.golang.org/grpc"

	filepb "github.com/jsannemo/omogenjudge/filehandler/api"
	fhclient "github.com/jsannemo/omogenjudge/filehandler/client"
	"github.com/jsannemo/omogenjudge/masterjudge/queue"
	"github.com/jsannemo/omogenjudge/masterjudge/service"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
	rclient "github.com/jsannemo/omogenjudge/runner/client"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/problems"
	"github.com/jsannemo/omogenjudge/storage/submissions"
	"github.com/jsannemo/omogenjudge/util/go/files"
	"github.com/jsannemo/omogenjudge/util/go/filestore"
	"github.com/jsannemo/omogenjudge/util/go/users"
)

var (
	listenAddr = flag.String("master_listen_addr", "127.0.0.1:61813", "The master server address to listen to in the format host:port")
)

var runner runpb.RunServiceClient
var filehandler filepb.FileHandlerServiceClient
var slotPool = make(chan int, 4)

func init() {
	slotPool <- 0
	slotPool <- 1
	slotPool <- 2
	slotPool <- 3
}

func compile(ctx context.Context, s *models.Submission, output chan<- *runpb.CompileResponse, outerr **error) {
	response, err := runner.Compile(ctx,
		&runpb.CompileRequest{
			Program:    s.ToRunnerProgram(),
			OutputPath: fmt.Sprintf("/var/lib/omogen/submissions/%d/program", s.SubmissionID),
		})
	if err != nil {
		*outerr = &err
	} else {
		output <- response
	}
	close(output)
}

func toProgram(val *models.OutputValidator) (*runpb.Program, error) {
	program := &runpb.Program{LanguageId: val.ValidatorLanguageID.String}
	contents, err := filestore.GetFile(val.ValidatorSourceZIP.Url)
	r, err := zip.NewReader(bytes.NewReader(contents), int64(len(contents)))
	if err != nil {
		return nil, err
	}
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		buf, err := ioutil.ReadAll(rc)
		if err != nil {
			return nil, err
		}
		program.Sources = append(program.Sources, &runpb.SourceFile{
			Path:     f.Name,
			Contents: string(buf),
		})
	}
	return program, nil
}

func judge(ctx context.Context, run *models.SubmissionRun) error {
	slot := <-slotPool
	defer func() { slotPool <- slot }()
	logger.Infof("Judging run %d", run.SubmissionRunID)
	subs, err := submissions.ListSubmissions(ctx, submissions.ListArgs{WithFiles:true}, submissions.ListFilter{SubmissionID: []int32{run.SubmissionID}})
	if err != nil {
		return err
	}
	if len(subs) != 1 {
		return fmt.Errorf("could not find submission %d", run.SubmissionID)
	}
	submission := subs[0]
	// TODO: this should be created in the local evaluator prior to compilation somehow
	root := files.NewFileBase(fmt.Sprintf("/var/lib/omogen/submissions/%d", run.SubmissionRunID))
	root.Gid = users.OmogenClientsID()
	root.GroupWritable = true
	if err := root.Mkdir("."); err != nil {
		return fmt.Errorf("Could not create submission directory: %v", err)
	}

	var compileErr *error = nil
	compileOutput := make(chan *runpb.CompileResponse)
	go compile(ctx, submission, compileOutput, &compileErr)
	run.Status = models.StatusCompiling
	if err := submissions.UpdateRun(ctx, run, submissions.UpdateRunArgs{Fields: []submissions.RunField{submissions.RunFieldStatus}}); err != nil {
		return err
	}
	probs, err := problems.List(ctx, problems.ListArgs{WithTests: problems.TestsAll}, problems.ListFilter{ProblemId: []int32{submission.ProblemID}})
	if err != nil {
		return err
	}
	if len(probs) == 0 {
		return fmt.Errorf("requested problem %d, got %d problems", submission.ProblemID, len(probs))
	}
	problem := probs[0]

	var testHandles []*filepb.FileHandle
	for _, tg := range problem.CurrentVersion.TestGroups {
		for _, tc := range tg.Tests {
			testHandles = append(testHandles, &filepb.FileHandle{Sha256Hash: tc.InputFile.Hash, Url: tc.InputFile.URL})
			testHandles = append(testHandles, &filepb.FileHandle{Sha256Hash: tc.OutputFile.Hash, Url: tc.OutputFile.URL})
		}
	}
	resp, err := filehandler.EnsureFile(ctx, &filepb.EnsureFileRequest{Handles: testHandles})
	if err != nil {
		return err
	}
	testPathMap := make(map[string]string)
	for i, handle := range testHandles {
		testPathMap[handle.Sha256Hash] = resp.Paths[i]
	}

	var validatorProgram *runpb.CompiledProgram
	validator := problem.CurrentVersion.OutputValidator
	if !validator.Nil() {
		outputValidator, err := toProgram(validator)
		if err != nil {
			return err
		}
		validatorId := validator.ValidatorSourceZIP.Hash.String
		resp, err := runner.CompileCached(ctx,
			&runpb.CompileCachedRequest{
				Identifier: validatorId,
				Request: &runpb.CompileRequest{
					Program:    outputValidator,
					OutputPath: fmt.Sprintf("/var/lib/omogen/tmps/val-%s", validatorId),
				},
			})
		if err != nil {
			return fmt.Errorf("failed calling compilier: %v", err)
		}
		if resp.Response.Program == nil {
			return fmt.Errorf("failed compiling output validator: %v", resp.Response.CompilationError)
		}
		validatorProgram = resp.Response.Program
	}

	compileResponse := <-compileOutput
	if compileResponse == nil {
		return fmt.Errorf("Compilation crashed: %v", *compileErr)
	}
	if compileResponse.Program == nil {
		run.Status = models.StatusSuccessful
		run.Verdict = models.VerdictCompilationError
		run.CompileError = sql.NullString{compileResponse.CompilationError, true}
		if err := submissions.UpdateRun(ctx, run, submissions.UpdateRunArgs{Fields: []submissions.RunField{submissions.RunFieldStatus, submissions.RunFieldVerdict, submissions.RunFieldCompileError}}); err != nil {
			return err
		}
		return nil
	}

	run.Status = models.StatusRunning
	if err := submissions.UpdateRun(ctx, run, submissions.UpdateRunArgs{Fields: []submissions.RunField{submissions.RunFieldStatus}}); err != nil {
		return err
	}
	var reqGroups []*runpb.TestGroup
	for _, group := range problem.CurrentVersion.TestGroups {
		var reqTests []*runpb.TestCase
		for _, test := range group.Tests {
			reqTests = append(reqTests, &runpb.TestCase{
				Name:       fmt.Sprintf("test-%d", test.TestCaseID),
				InputPath:  testPathMap[test.InputFile.Hash],
				OutputPath: testPathMap[test.OutputFile.Hash],
			})
		}
		reqGroups = append(reqGroups, &runpb.TestGroup{
			Cases: reqTests,
			Score: group.Score,
		})
	}
	evalReq := &runpb.EvaluateRequest{
		SubmissionId: string(submission.SubmissionID),
		Program:      compileResponse.Program,
		Groups:       reqGroups,
		TimeLimitMs:  problem.CurrentVersion.TimeLimMS,
		MemLimitKb:   problem.CurrentVersion.MemLimKB,
	}
	if validatorProgram != nil {
		evalReq.Validator = &runpb.CustomValidator{Program: validatorProgram}
	}
	stream, err := runner.Evaluate(ctx, evalReq)
	if err != nil {
		return err
	}
	logger.Infof("eval: %v", evalReq)
	for {
		verdict, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		logger.Infof("Verdict: %v", verdict)
		switch res := verdict.Result.(type) {
		case *runpb.EvaluateResponse_TestCase:
			// TODO report this
		case *runpb.EvaluateResponse_TestGroup:
			// TODO report this
		case *runpb.EvaluateResponse_Submission:
			run.Verdict = models.Verdict(res.Submission.Verdict.String())
			run.Status = models.StatusSuccessful
			run.Score = res.Submission.Score
			if err := submissions.UpdateRun(ctx, run, submissions.UpdateRunArgs{Fields: []submissions.RunField{submissions.RunFieldStatus, submissions.RunFieldVerdict, submissions.RunFieldScore}}); err != nil {
				return err
			}
		case nil:
			logger.Warningf("Got empty response")
		default:
			return fmt.Errorf("Got unexpected result %T", res)
		}

	}
	return nil
}

func main() {
	flag.Parse()
	defer logger.Init("masterjudge", false, true, os.Stderr).Close()

	var err error
	runner, err = rclient.NewClient()
	if err != nil {
		logger.Fatalf("Failed creating runner client: %v", err)
	}
	filehandler = fhclient.NewClient()

	judgeQueue := make(chan *models.SubmissionRun, 1000)
	if err := queue.StartQueue(context.Background(), judgeQueue); err != nil {
		logger.Fatalf("Failed queue startup: %v", err)
	}
	for w := 1; w <= 10; w++ {
		go func() {
			for {
				if err := judge(context.Background(), <-judgeQueue); err != nil {
					logger.Errorf("Failed judging submission: %v", err)
				}
			}
		}()
	}

	lis, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	if err != nil {
		logger.Fatalf("failed to create server: %v", err)
	}
	service.Register(grpcServer)
	logger.Infof("serving on %v", *listenAddr)
	grpcServer.Serve(lis)
}
