package problems

import (
	"context"
	"fmt"
	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
	"io"
	"path/filepath"
)

func verifySubmissions(ctx context.Context, path string, problem *toolspb.Problem, outputValidator *runpb.CompiledProgram, runner runpb.RunServiceClient, reporter util.Reporter) error {

	for i, submission := range problem.Submissions {
		resp, err := runner.Compile(ctx, &runpb.CompileRequest{
			Program:    submission.Submission,
			OutputPath: filepath.Join(path, fmt.Sprintf("submission_%d", i)),
		})
		if err != nil {
			return err
		}
		if resp.Program == nil {
			reporter.Err("compilation of submission failed: %v", resp.CompilationError)
			return nil
		}
		var cases []*runpb.TestCase
		for _, g := range problem.TestGroups {
			for _, tc := range g.Tests {
				cases = append(cases,
					&runpb.TestCase{
						Name:       tc.Name,
						InputPath:  tc.InputPath,
						OutputPath: tc.OutputPath,
					})
			}
		}
		// TODO: this should not use 0
		req := &runpb.EvaluateRequest{
			SubmissionId: 0,
			Program:      resp.Program,
			MemLimitKb:   int64(problem.Metadata.Limits.MemoryLimitKb),
			TimeLimitMs:  int64(problem.Metadata.Limits.TimeLimitMs),
			EvaluateAll:  true,
			Cases:        cases,
		}
		if outputValidator != nil {
			req.Validator = &runpb.CustomValidator{
				Program: outputValidator,
			}
		}

		stream, err := runner.Evaluate(ctx, req)
		if err != nil {
			return err
		}
		observedVerdict := make(map[runpb.Verdict]bool)
		allowedVerdicts := make(map[runpb.Verdict]bool)
		for _, allowed := range submission.Constraint.AllowedFailures {
			allowedVerdicts[allowed] = true
		}
		i := 0
		for {
			verdict, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			switch res := verdict.Result.(type) {
			case *runpb.EvaluateResponse_TestCase:
				v := res.TestCase.Verdict
				observedVerdict[v] = true
				if !allowedVerdicts[v] {
					reporter.Err("Submission %s got forbidden verdict %s on case %s", submission.Name, v, cases[i].Name)
				}
			default:
			}
			i++
		}
		for _, required := range submission.Constraint.RequiredFailures {
			if !observedVerdict[required] {
				reporter.Err("Submission %s required verdict %s, but did not cause it", submission.Name, required)
			}
		}
	}
	return nil
}
