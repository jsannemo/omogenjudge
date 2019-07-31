package main

import (
  "context"
  "flag"
  "fmt"
	"os"
  "io"

  "github.com/google/logger"

	"github.com/jsannemo/omogenjudge/judgemaster/queue"
	"github.com/jsannemo/omogenjudge/storage/submissions"
	"github.com/jsannemo/omogenjudge/storage/problems"
  runpb "github.com/jsannemo/omogenjudge/runner/api"
  rclient "github.com/jsannemo/omogenjudge/runner/client"
  filepb "github.com/jsannemo/omogenjudge/filehandler/api"
  fhclient "github.com/jsannemo/omogenjudge/filehandler/client"
)

var runner runpb.RunServiceClient
var filehandler filepb.FileHandlerServiceClient

func compile(s *submissions.Submission, output chan<- *runpb.CompiledProgram, outerr **error) {
  response, err := runner.Compile(context.TODO(),
    &runpb.CompileRequest{
      Program: s.ToRunnerProgram(),
      OutputPath: fmt.Sprintf("/var/lib/omogen/submissions/%d/program", s.SubmissionId),
    })
  if err != nil {
    *outerr = &err
  } else {
    output <- response.Program
  }
  close(output)
}


func judge(id int32) error {
  submission, err := submissions.GetSubmission(context.TODO(), id)
  logger.Infof("Judging submission %d", id)
  if err != nil {
    return err
  }
  var compileErr *error = nil
  compileOutput := make(chan *runpb.CompiledProgram)
  go compile(submission, compileOutput, &compileErr)
  problem, err := problems.GetProblemForJudging(context.TODO(), submission.ProblemId)
  if err != nil {
    return err
  }
  testHandles := problem.TestDataFiles().ToHandlerFiles()
  resp, err := filehandler.EnsureFile(context.TODO(), &filepb.EnsureFileRequest{Handles: testHandles})
  if err != nil {
    return err
  }
  testPathMap := make(map[string]string)
  for i, handle := range testHandles {
    testPathMap[handle.Sha256Hash] = resp.Paths[i]
  }
  compiledProgram := <-compileOutput
  if compiledProgram == nil {
    return *compileErr
  }

  tests := problem.Tests()
  var reqTests []*runpb.TestCase
  for _, test := range tests{
    reqTests = append(reqTests, &runpb.TestCase{
      Name: fmt.Sprintf("test-%d", test.TestCaseId),
      InputPath: testPathMap[test.InputFile.Hash],
      OutputPath: testPathMap[test.OutputFile.Hash],
    })
  }
	stream, err := runner.Evaluate(context.Background(),
    &runpb.EvaluateRequest{
      SubmissionId: id,
      Program: compiledProgram,
      Cases: reqTests,
      TimeLimitMs: 1000,
      MemLimitKb: 1024 * 1000,
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
  }
  return nil
}

func main() {
	flag.Parse()
  defer logger.Init("evaluator", false, true, os.Stderr).Close()

  runner = rclient.NewClient()
  filehandler = fhclient.NewClient()

  judgeQueue := make(chan int32, 1000)
  if err := queue.StartQueue(judgeQueue); err != nil {
    logger.Fatalf("Failed queue startup: %v", err)
  }
  for w := 1; w <= 10; w++ {
    go func() {
      for {
        if err := judge(<-judgeQueue); err != nil {
          logger.Errorf("Failed judging submission: %v", err)
        }
      }
    }()
  }
  select{ }
}
