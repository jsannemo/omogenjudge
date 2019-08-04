package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
  "database/sql"
	"sync"

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
)

var (
	listenAddr = flag.String("master_listen_addr", "127.0.0.1:61813", "The master server address to listen to in the format host:port")
)

var runner runpb.RunServiceClient
var filehandler filepb.FileHandlerServiceClient
var slotPool = sync.Pool{}

func init() {
	slotPool.Put(0)
	slotPool.Put(1)
	slotPool.Put(2)
	slotPool.Put(3)
}

func compile(ctx context.Context, s *models.Submission, output chan<- *runpb.CompileResponse, outerr **error) {
	response, err := runner.Compile(ctx,
		&runpb.CompileRequest{
			Program:    s.ToRunnerProgram(),
			OutputPath: fmt.Sprintf("/var/lib/omogen/submissions/%d/program", s.SubmissionId),
		})
	if err != nil {
		*outerr = &err
	} else {
		output <- response
	}
	close(output)
}

func judge(ctx context.Context, submission *models.Submission) error {
	slot := slotPool.Get()
	defer slotPool.Put(slot)
	logger.Infof("Judging submission %d", submission.SubmissionId)
	var compileErr *error = nil
	compileOutput := make(chan *runpb.CompileResponse)
	go compile(ctx, submission, compileOutput, &compileErr)
	submission.Status = models.StatusCompiling
	submissions.Update(ctx, submission, submissions.UpdateArgs{Fields: []submissions.Field{submissions.FieldStatus}})
	problems := problems.List(ctx, problems.ListArgs{WithTests: problems.TestsAll}, problems.ListFilter{ProblemId: submission.ProblemId})
	if len(problems) == 0 {
		return fmt.Errorf("requested problem %d, got %d problems", submission.ProblemId, len(problems))
	}
	problem := problems[0]
	testHandles := problem.TestDataFiles().ToHandlerFiles()
	resp, err := filehandler.EnsureFile(ctx, &filepb.EnsureFileRequest{Handles: testHandles})
	if err != nil {
		return err
	}
	testPathMap := make(map[string]string)
	for i, handle := range testHandles {
		testPathMap[handle.Sha256Hash] = resp.Paths[i]
	}
	compileResponse := <-compileOutput
	if compileResponse == nil {
    return fmt.Errorf("Compilation crashed: %v", compileErr)
  }
  if compileResponse.Program == nil {
		submission.Status = models.StatusCompilationFailed
    submission.CompileError = sql.NullString{compileResponse.CompilationError, true}
		submissions.Update(ctx, submission, submissions.UpdateArgs{Fields: []submissions.Field{submissions.FieldStatus, submissions.FieldCompileError}})
		return nil
	}
	submission.Status = models.StatusRunning
	submissions.Update(ctx, submission, submissions.UpdateArgs{Fields: []submissions.Field{submissions.FieldStatus}})

	tests := problem.Tests()
	var reqTests []*runpb.TestCase
	for _, test := range tests {
		reqTests = append(reqTests, &runpb.TestCase{
			Name:       fmt.Sprintf("test-%d", test.TestCaseId),
			InputPath:  testPathMap[test.InputFile.Hash],
			OutputPath: testPathMap[test.OutputFile.Hash],
		})
	}
	stream, err := runner.Evaluate(ctx,
		&runpb.EvaluateRequest{
			SubmissionId: submission.SubmissionId,
			Program:      compileResponse.Program,
			Cases:        reqTests,
			TimeLimitMs:  int64(problem.TimeLimMs),
			MemLimitKb:   int64(problem.MemLimKb),
		})
	if err != nil {
		return err
	}
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
		case *runpb.EvaluateResponse_Submission:
			submission.Verdict = models.Verdict(res.Submission.Verdict.String())
			submission.Status = models.StatusSuccessful
			submissions.Update(ctx, submission, submissions.UpdateArgs{Fields: []submissions.Field{submissions.FieldVerdict, submissions.FieldStatus}})
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

	runner = rclient.NewClient()
	filehandler = fhclient.NewClient()

	judgeQueue := make(chan *models.Submission, 1000)
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
