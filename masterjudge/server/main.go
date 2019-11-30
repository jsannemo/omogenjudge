package main

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/google/logger"
	"github.com/jsannemo/omogenjudge/masterjudge/queue"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"

	filepb "github.com/jsannemo/omogenjudge/filehandler/api"
	fhclient "github.com/jsannemo/omogenjudge/filehandler/client"
	"github.com/jsannemo/omogenjudge/masterjudge/service"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
	rclient "github.com/jsannemo/omogenjudge/runner/client"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/problems"
	"github.com/jsannemo/omogenjudge/storage/submissions"
	"github.com/jsannemo/omogenjudge/util/go/filestore"
)

var (
	listenAddr = flag.String("master_listen_addr", "127.0.0.1:61813", "The master server address to listen to in the format host:port")
)

var runner runpb.RunServiceClient
var filehandler filepb.FileHandlerServiceClient

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

func getValidatorProgram(val *models.OutputValidator) (*runpb.Program, error) {
	program := &runpb.Program{LanguageId: val.ValidatorLanguageID.String}
	contents, err := filestore.GetFile(val.ValidatorSourceZIP.URL)
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
	logger.Infof("Judging run %d", run.SubmissionRunID)
	subs, err := submissions.ListSubmissions(ctx, submissions.ListArgs{WithFiles: true}, submissions.ListFilter{SubmissionID: []int32{run.SubmissionID}})
	if err != nil {
		return err
	}
	if len(subs) != 1 {
		return fmt.Errorf("could not find submission %d", run.SubmissionID)
	}
	submission := subs[0]
	var compileErr *error = nil
	compileOutput := make(chan *runpb.CompileResponse)
	go compile(ctx, submission, compileOutput, &compileErr)
	run.Status = models.StatusCompiling
	if err := submissions.UpdateRun(ctx, run, submissions.UpdateRunArgs{Fields: []submissions.RunField{submissions.RunFieldStatus}}); err != nil {
		return fmt.Errorf("failed settin run as compiling: %v", err)
	}
	probs, err := problems.List(ctx, problems.ListArgs{WithTests: problems.TestsAll}, problems.ListFilter{ProblemID: []int32{submission.ProblemID}})
	if err != nil {
		return fmt.Errorf("failed retreiving problem info: %v", err)
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
		return fmt.Errorf("failed ensuring file: %v", err)
	}
	testPathMap := make(map[string]string)
	for i, handle := range testHandles {
		testPathMap[handle.Sha256Hash] = resp.Paths[i]
	}

	var validatorProgram *runpb.CompiledProgram
	validator := problem.CurrentVersion.OutputValidator
	if !validator.ValidatorSourceZIP.Nil() {
		outputValidator, err := getValidatorProgram(validator)
		if err != nil {
			return err
		}
		validatorID := validator.ValidatorSourceZIP.Hash
		resp, err := runner.CompileCached(ctx,
			&runpb.CompileCachedRequest{
				Identifier: validatorID.String,
				Request: &runpb.CompileRequest{
					Program:    outputValidator,
					OutputPath: fmt.Sprintf("/var/lib/omogen/tmps/val-%s", validatorID),
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
		run.Status = models.StatusCompilationFailed
		run.Verdict = models.VerdictUnjudged
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
		SubmissionId: strconv.Itoa(int(submission.SubmissionID)),
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
	atGroup := 0
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
			// TODO(jsannemo): report this
		case *runpb.EvaluateResponse_TestGroup:
			groupRun := &models.TestGroupRun{
				SubmissionRunID: run.SubmissionRunID,
				TestGroupID:     problem.CurrentVersion.TestGroups[atGroup].TestGroupID,
				Evaluation: models.Evaluation{
					Score:       res.TestGroup.Score,
					TimeUsageMS: res.TestGroup.TimeUsageMs,
					Verdict:     models.VerdictFromRunVerdict(res.TestGroup.Verdict),
				},
			}
			if err := submissions.CreateGroupRun(ctx, groupRun); err != nil {
				return err
			}
			atGroup++
		case *runpb.EvaluateResponse_Submission:
			run.Verdict = models.VerdictFromRunVerdict(res.Submission.Verdict)
			run.Status = models.StatusSuccessful
			run.Score = res.Submission.Score
			if err := submissions.UpdateRun(ctx, run, submissions.UpdateRunArgs{Fields: []submissions.RunField{submissions.RunFieldStatus, submissions.RunFieldVerdict, submissions.RunFieldScore}}); err != nil {
				return err
			}
		case nil:
			logger.Warningf("Got empty response")
		default:
			return fmt.Errorf("unexpected eval result: %T", res)
		}

	}
	return nil
}

func main() {
	flag.Parse()
	defer logger.Init("omogenjudge-master", false, true, os.Stderr).Close()

	var err error
	runner, err = rclient.NewClient()
	if err != nil {
		logger.Fatalf("Failed creating runner client: %v", err)
	}
	filehandler = fhclient.NewClient()

	judgeQueue := make(chan *models.SubmissionRun, 1000)
	for w := 1; w <= 4; w++ {
		go func() {
			for {
				if err := judge(context.Background(), <-judgeQueue); err != nil {
					logger.Errorf("Failed judging submission: %v", err)
				}
			}
		}()
	}
	if err := queue.StartQueue(context.Background(), judgeQueue); err != nil {
		logger.Fatalf("Failed queue startup: %v", err)
	}

	lis, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	if err != nil {
		logger.Fatalf("failed to create server: %v", err)
	}
	if err := service.Register(grpcServer); err != nil {
		logger.Fatalf("failed to register server: %v", err)
	}
	logger.Infof("serving on %v", *listenAddr)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
}
