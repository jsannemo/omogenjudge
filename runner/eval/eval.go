package eval

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	_ "github.com/google/logger"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
	"github.com/jsannemo/omogenjudge/runner/diff"
	"github.com/jsannemo/omogenjudge/runner/runners"
)

type Evaluator struct {
	root        string
	linker      *runners.FileLinker
	program     runners.Program
	valLinker   *runners.FileLinker
	validator   runners.Program
	EvaluateAll bool
}

func NewEvaluator(root string, program runners.Program, validator runners.Program) (*Evaluator, error) {
	fl, err := runners.NewFileLinker(filepath.Join(root, "env"))
	if err != nil {
		return nil, fmt.Errorf("failed creating FileLinker: %v", err)
	}
	eval := &Evaluator{
		root:      root,
		linker:    fl,
		program:   program,
		validator: validator}
	if validator != nil {
		valfl, err := runners.NewFileLinker(filepath.Join(root, "valenv"))
		if err != nil {
			return nil, fmt.Errorf("failed creating validator FileLinker: %v", err)
		}
		eval.valLinker = valfl
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
	defer e.linker.Clear()
	if e.valLinker != nil {
		defer e.valLinker.Clear()
	}
	outPath := e.linker.PathFor("output", true)
	e.program.SetArgs(&runners.ProgramArgs{
		InputPath:     e.linker.PathFor("input", false),
		OutputPath:    outPath,
		ErrorPath:     e.linker.PathFor("error", true),
		TimeLimitMs:   timeLimMs,
		MemoryLimitKb: memLimitKb,
	})
	if e.validator != nil {
		e.validator.SetArgs(&runners.ProgramArgs{
			InputPath:  e.valLinker.PathFor("team_output", false),
			OutputPath: e.valLinker.PathFor("output", true),
			ErrorPath:  e.valLinker.PathFor("error", true),
			// TODO make this configurable
			TimeLimitMs:   2000,
			MemoryLimitKb: 1000 * 1000,
			ExtraArgs: []string{
				e.valLinker.PathFor("input", false),
				e.valLinker.PathFor("judge_answer", false),
				filepath.Join(e.valLinker.PathFor("fedeback", true)),
			},
		})
	}

	verdict := runpb.Verdict_ACCEPTED
	for _, tc := range testCases {
		tcPath := filepath.Join(e.root, tc.Name)
		if err := os.MkdirAll(tcPath, 0755); err != nil {
			return err
		}

		if err := e.linker.LinkFile(tc.InputPath, "input", false); err != nil {
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
				if err := e.valLinker.LinkFile(tc.InputPath, "input", false); err != nil {
					return err
				}
				if err := e.valLinker.LinkFile(outPath, "team_output", false); err != nil {
					return err
				}
				if err := e.valLinker.LinkFile(tc.OutputPath, "judge_answer", false); err != nil {
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
					dat, err := ioutil.ReadFile(e.valLinker.PathFor("error", true))
					if err != nil {
						return fmt.Errorf("could not read output validator errors: %v", err)
					}
					dat2, err := ioutil.ReadFile(e.valLinker.PathFor("output", true))
					if err != nil {
						return fmt.Errorf("could not read output validator output: %v", err)
					}
					return fmt.Errorf("output validator crashed: %v", string(dat)+" "+string(dat2))
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

		e.linker.Clear()
		if e.valLinker != nil {
			e.valLinker.Clear()
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
