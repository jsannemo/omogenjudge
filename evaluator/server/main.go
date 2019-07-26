package main

import (
  "context"
  "flag"
  "fmt"
  "path/filepath"
  "errors"

  "github.com/google/logger"

	"github.com/jsannemo/omogenjudge/evaluator/queue"
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
	stream, err := runner.Run(context.Background())
	validatorStream, err := runner.Run(context.Background())
  if err != nil {
    return err
  }
  verdict := "AC"
  for _, test := range tests {
    err := stream.Send(&runpb.RunRequest{
      Program: compiledProgram,
      InputPath: testPathMap[test.InputFile.Hash],
      OutputPath: fmt.Sprintf("/var/lib/omogen/submissions/%d/%d/output", id, test.TestCaseId),
      ErrorPath: fmt.Sprintf("/var/lib/omogen/submissions/%d/%d/error", id, test.TestCaseId),
      TimeLimitMs: 2000,
      MemoryLimitKb: 512 * 1024,
    })
    if err != nil {
      return err
    }
    in, err := stream.Recv()
    if err != nil {
      return err
    }
    logger.Info(in)
    switch x := in.Exit.(type) {
		case *runpb.RunResponse_TimeLimitExceeded:
      verdict = "TLE"
      break
		case *runpb.RunResponse_MemoryLimitExceeded:
      verdict = "RTE"
      break
		case *runpb.RunResponse_Signaled:
      verdict = "RTE"
      break
		case *runpb.RunResponse_Exited:
      if x.Exited.ExitCode != 0 {
        verdict = "RTE"
      }
		default:
      return errors.New("No exit set")
		}

    err = validatorStream.Send(&runpb.RunRequest{
      Program: &runpb.CompiledProgram {
        ProgramRoot: "/usr/bin/diff",
        LanguageId: "cmd",
      },
      InputPath: fmt.Sprintf("/var/lib/omogen/submissions/%d/%d/output", id, test.TestCaseId),
      OutputPath: fmt.Sprintf("/var/lib/omogen/submissions/%d/%d/valoutput", id, test.TestCaseId),
      ErrorPath: fmt.Sprintf("/var/lib/omogen/submissions/%d/%d/valerror", id, test.TestCaseId),
      TimeLimitMs: 2000,
      MemoryLimitKb: 512 * 1024,
      ExtraReadPaths: []string{
        fmt.Sprintf("/var/lib/omogen/submissions/%d/%d", id, test.TestCaseId),
        filepath.Dir(testPathMap[test.OutputFile.Hash]),
      },
      Args: []string{
        testPathMap[test.OutputFile.Hash],
        "-",
        "-b",
      },
    })
    if err != nil {
      return err
    }
    val, err := validatorStream.Recv()
    if err != nil {
      return err
    }
    logger.Info(val)
    switch x := val.Exit.(type) {
		case *runpb.RunResponse_Exited:
      if x.Exited.ExitCode == 0 {
      } else if x.Exited.ExitCode == 1 {
        verdict = "WA"
      } else {
        return errors.New("Validator crashed :(")
      }
		default:
      return errors.New("Invalid exit for validator :(")
		}
  }
  logger.Infof("Verdict: %s", verdict)
  stream.CloseSend()
  validatorStream.CloseSend()
  return nil
}

func main() {
	flag.Parse()
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
