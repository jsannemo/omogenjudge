package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/google/logger"

	filepb "github.com/jsannemo/omogenjudge/filehandler/api"
	fhclient "github.com/jsannemo/omogenjudge/filehandler/client"
	"github.com/jsannemo/omogenjudge/masterjudge/queue"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
	rclient "github.com/jsannemo/omogenjudge/runner/client"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/problems"
	"github.com/jsannemo/omogenjudge/storage/submissions"
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

func compile(ctx context.Context, s *models.Submission, output chan<- *runpb.CompiledProgram, outerr **error) {
	response, err := runner.Compile(ctx,
		&runpb.CompileRequest{
			Program:    s.ToRunnerProgram(),
			OutputPath: fmt.Sprintf("/var/lib/omogen/submissions/%d/program", s.SubmissionId),
		})
	if err != nil {
		*outerr = &err
	} else {
		output <- response.Program
	}
	close(output)
}

func judge(ctx context.Context, submission *models.Submission) error {
	slot := slotPool.Get()
	defer slotPool.Put(slot)
	logger.Infof("Judging submission %d", submission.SubmissionId)
	var compileErr *error = nil
	compileOutput := make(chan *runpb.CompiledProgram)
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
	compiledProgram := <-compileOutput
	if compiledProgram == nil {
		submission.Status = models.StatusCompilationFailed
		submissions.Update(ctx, submission, submissions.UpdateArgs{Fields: []submissions.Field{submissions.FieldStatus}})
		logger.Infof("Compilation failure: %v", *compileErr)
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
			Program:      compiledProgram,
			Cases:        reqTests,
			TimeLimitMs:  Problem.TimeLimMs,
			MemLimitKb:   Problem.MemLimKb,
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
	select {}
}
