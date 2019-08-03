package eval

import (
	"fmt"
	"os"
	"path/filepath"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
	"github.com/jsannemo/omogenjudge/runner/diff"
	"github.com/jsannemo/omogenjudge/runner/runners"
)

type Evaluator struct {
	root    string
	env     *runners.Env
	program runners.Program
}

func NewEvaluator(root string, program runners.Program) (*Evaluator, error) {
	env, err := runners.NewEnv(filepath.Join(root, "env"))
	if err != nil {
		return nil, fmt.Errorf("failed creating Env: %v", err)
	}
	return &Evaluator{root, env, program}, nil
}

type TestCase struct {
	Name       string
	InputPath  string
	OutputPath string
}

type Result struct {
	TestCaseVerdict   runpb.Verdict
	SubmissionVerdict runpb.Verdict
}

func (e *Evaluator) Evaluate(testCases []*TestCase, timeLimMs, memLimitKb int64, results chan<- *Result) error {
	defer close(results)
	defer e.env.Clear()
	e.program.SetArgs(&runners.ProgramArgs{
		InputPath:     e.env.PathFor("input", false),
		OutputPath:    e.env.PathFor("output", true),
		ErrorPath:     e.env.PathFor("error", true),
		TimeLimitMs:   timeLimMs,
		MemoryLimitKb: memLimitKb,
	})

	verdict := runpb.Verdict_ACCEPTED
	for _, tc := range testCases {
		tcPath := filepath.Join(e.root, tc.Name)
		if err := os.MkdirAll(tcPath, 0755); err != nil {
			return err
		}
		outPath := filepath.Join(tcPath, "output")
		if _, err := os.Create(outPath); err != nil {
			return err
		}
		errPath := filepath.Join(tcPath, "error")
		if _, err := os.Create(errPath); err != nil {
			return err
		}

		if err := e.env.LinkFile(tc.InputPath, "input", false); err != nil {
			return err
		}
		if err := e.env.LinkFile(outPath, "output", true); err != nil {
			return err
		}
		if err := e.env.LinkFile(errPath, "error", true); err != nil {
			return err
		}

		exit, err := e.program.Execute()
		if err != nil {
			return err
		}

		if crashed(exit) {
			results <- &Result{TestCaseVerdict: runpb.Verdict_RUN_TIME_ERROR}
			verdict = runpb.Verdict_RUN_TIME_ERROR
			break
		} else if timedOut(exit) {
			results <- &Result{TestCaseVerdict: runpb.Verdict_TIME_LIMIT_EXCEEDED}
			verdict = runpb.Verdict_TIME_LIMIT_EXCEEDED
			break
		} else {
			hasDiff, err := diffOutput(tc.OutputPath, outPath)
			if err != nil {
				return err
			}
			if hasDiff {
				results <- &Result{TestCaseVerdict: runpb.Verdict_WRONG_ANSWER}
				verdict = runpb.Verdict_WRONG_ANSWER
				break
			} else {
				results <- &Result{TestCaseVerdict: runpb.Verdict_ACCEPTED}
			}
		}

		e.env.Clear()
	}
	results <- &Result{SubmissionVerdict: verdict}
	return nil
}

func diffOutput(refPath, outPath string) (bool, error) {
	refFile, err := os.Open(refPath)
	if err != nil {
		return false, err
	}
	outFile, err := os.Open(outPath)
	if err != nil {
		return false, err
	}
	res, err := diff.Diff(refFile, outFile)
	return !res.Match, err
}

func crashed(res *runners.ExecResult) bool {
	return (res.ExitReason == runners.Exited && res.ExitCode != 0) || res.ExitReason == runners.Signaled
}

func timedOut(res *runners.ExecResult) bool {
	return res.ExitReason == runners.TimedOut
}
