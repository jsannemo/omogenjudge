package eval

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
	"github.com/jsannemo/omogenjudge/runner/diff"
	"github.com/jsannemo/omogenjudge/runner/runners"
)

type Evaluator struct {
	root      string
	env       *runners.Env
	program   runners.Program
	valenv    *runners.Env
	validator runners.Program
}

func NewEvaluator(root string, program runners.Program, validator runners.Program) (*Evaluator, error) {
	env, err := runners.NewEnv(filepath.Join(root, "env"))
	if err != nil {
		return nil, fmt.Errorf("failed creating Env: %v", err)
	}
	eval := &Evaluator{
		root:      root,
		env:       env,
		program:   program,
		validator: validator}
	if validator != nil {
		valenv, err := runners.NewEnv(filepath.Join(root, "valenv"))
		if err != nil {
			return nil, fmt.Errorf("failed creating Env: %v", err)
		}
		eval.valenv = valenv
	}
	return eval, nil
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
	if e.valenv != nil {
		defer e.valenv.Clear()
	}
	e.program.SetArgs(&runners.ProgramArgs{
		InputPath:     e.env.PathFor("input", false),
		OutputPath:    e.env.PathFor("output", true),
		ErrorPath:     e.env.PathFor("error", true),
		TimeLimitMs:   timeLimMs,
		MemoryLimitKb: memLimitKb,
	})
	if e.validator != nil {
		e.validator.SetArgs(&runners.ProgramArgs{
			InputPath:  e.valenv.PathFor("team_output", false),
			OutputPath: e.valenv.PathFor("output", true),
			ErrorPath:  e.valenv.PathFor("error", true),
			// TODO make this configurable
			TimeLimitMs:   2000,
			MemoryLimitKb: 1000 * 1000,
			ExtraArgs: []string{
				e.valenv.PathFor("input", false),
				e.valenv.PathFor("team_output", false),
				e.valenv.PathFor("judge_answer", false),
				filepath.Join(e.valenv.WriteRoot, "feedback"),
			},
		})
	}

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
		} else {
			if e.validator != nil {
				if err := e.valenv.LinkFile(tc.InputPath, "input", false); err != nil {
					return err
				}
				if err := e.valenv.LinkFile(outPath, "team_output", false); err != nil {
					return err
				}
				if err := e.valenv.LinkFile(tc.OutputPath, "judge_answer", false); err != nil {
					return err
				}
				exit, err := e.validator.Execute()
				if err != nil {
					return err
				}
				if timedOut(exit) {
					return fmt.Errorf("output validator timed out")
				}
				if crashedWith(exit, 42) {
					results <- &Result{TestCaseVerdict: runpb.Verdict_ACCEPTED}
				} else if crashedWith(exit, 43) {
					results <- &Result{TestCaseVerdict: runpb.Verdict_WRONG_ANSWER}
					verdict = runpb.Verdict_WRONG_ANSWER
					break
				} else {
					dat, _ := ioutil.ReadFile(e.valenv.PathFor("error", true))
					return fmt.Errorf("output validator crashed: %v", string(dat))
				}
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

func crashedWith(res *runners.ExecResult, code int) bool {
	return res.ExitReason == runners.Exited && res.ExitCode == code
}

func crashed(res *runners.ExecResult) bool {
	return (res.ExitReason == runners.Exited && res.ExitCode != 0) || res.ExitReason == runners.Signaled
}

func timedOut(res *runners.ExecResult) bool {
	return res.ExitReason == runners.TimedOut
}
