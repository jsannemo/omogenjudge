package eval

import (
  "os/exec"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

  _ "github.com/google/logger"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
	"github.com/jsannemo/omogenjudge/runner/diff"
	"github.com/jsannemo/omogenjudge/runner/runners"
)

type Evaluator struct {
	root        string
	env         *runners.Env
	program     runners.Program
	valenv      *runners.Env
	validator   runners.Program
	EvaluateAll bool
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

func (e *Evaluator) resetPermissions() error {
	cmd := exec.Command("/usr/bin/omogenjudge-permissionfixer", filepath.Base(e.root))
  return cmd.Run()
}

func (e *Evaluator) Evaluate(testCases []*TestCase, timeLimMs, memLimitKb int64, results chan<- *Result) error {
  if err := e.resetPermissions(); err != nil {
    return err
  }
	defer close(results)
	defer e.env.Clear()
	if e.valenv != nil {
		defer e.valenv.Clear()
	}
  outPath := e.env.PathFor("output", true)
	e.program.SetArgs(&runners.ProgramArgs{
		InputPath:     e.env.PathFor("input", false),
		OutputPath:    outPath,
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

		if err := e.env.LinkFile(tc.InputPath, "input", false); err != nil {
			return err
		}

		exit, err := e.program.Execute()
		if err != nil {
			return err
		}
    if err := e.resetPermissions(); err != nil {
      return err
    }

		if exit.Crashed() {
			results <- &Result{TestCaseVerdict: runpb.Verdict_RUN_TIME_ERROR}
			verdict = runpb.Verdict_RUN_TIME_ERROR
		} else if exit.TimedOut() {
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
        if err := e.resetPermissions(); err != nil {
          return err
        }

				if exit.TimedOut() {
					return fmt.Errorf("output validator timed out")
				}
				if exit.CrashedWith(42) {
					results <- &Result{TestCaseVerdict: runpb.Verdict_ACCEPTED}
				} else if exit.CrashedWith(43) {
					results <- &Result{TestCaseVerdict: runpb.Verdict_WRONG_ANSWER}
					verdict = runpb.Verdict_WRONG_ANSWER
				} else {
					// TODO handle error
					dat, err := ioutil.ReadFile(e.valenv.PathFor("error", true))
          if err != nil {
            return fmt.Errorf("could not read output validator errors: %v", err)
          }
					dat2, err := ioutil.ReadFile(e.valenv.PathFor("output", true))
          if err != nil {
            return fmt.Errorf("could not read output validator output: %v", err)
          }
					return fmt.Errorf("output validator crashed: %v", string(dat) + " " + string(dat2))
				}
			} else {
				hasDiff, err := diffOutput(tc.OutputPath, outPath)
				if err != nil {
					return err
				}
				if hasDiff {
					results <- &Result{TestCaseVerdict: runpb.Verdict_WRONG_ANSWER}
					verdict = runpb.Verdict_WRONG_ANSWER
				} else {
					results <- &Result{TestCaseVerdict: runpb.Verdict_ACCEPTED}
				}
			}
		}

		e.env.Clear()
		if e.valenv != nil {
			e.valenv.Clear()
		}
		if verdict != runpb.Verdict_ACCEPTED && !e.EvaluateAll {
			break
		}
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
