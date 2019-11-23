package problems

import (
	"context"
	"fmt"
	"io"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
	"github.com/jsannemo/omogenjudge/util/go/files"
	"github.com/jsannemo/omogenjudge/util/go/strings"
	"github.com/jsannemo/omogenjudge/util/go/users"
)

func verifySubmissions(ctx context.Context, path string, problem *toolspb.Problem, outputValidator *runpb.CompiledProgram, runner runpb.RunServiceClient, reporter util.Reporter) error {
	// If a time limit is specified, use it. Otherwise, we set a high value so we can determine time limit from the
	// submissions.
	tl := problem.Metadata.Limits.TimeLimitMs
	if tl == 0 {
		tl = 60 * 1000
	}
	maxTime := int32(0)
	for _, submission := range problem.Submissions {
		if !submission.UseForTiming {
			continue
		}
		time, err := verifySubmission(ctx, submission, tl, problem, outputValidator, runner, reporter)
		if err != nil {
			return err
		}
		if time > maxTime {
			maxTime = time
		}
	}
	// If we had no time limit specified, set it to multiplier * max time rounded upwards.
	if problem.Metadata.Limits.TimeLimitMs == 0 {
		problem.Metadata.Limits.TimeLimitMs = (problem.Metadata.Limits.TimeLimitMultiplier * maxTime + 999) / 1000 * 1000
	}
	for _, submission := range problem.Submissions {
		if submission.UseForTiming {
			continue
		}
		if _, err := verifySubmission(ctx, submission, problem.Metadata.Limits.TimeLimitMs, problem, outputValidator, runner, reporter); err != nil {
			return err
		}
	}
	return nil
}

func verifySubmission(ctx context.Context, submission *toolspb.Submission, timelim int32, problem *toolspb.Problem, outputValidator *runpb.CompiledProgram, runner runpb.RunServiceClient, reporter util.Reporter) (int32, error) {
	// TODO(jsannemo): remove the submission directory afterwards
	id := strings.RandStr(8) // 8*6 = 48 bits of entropy
	root := files.NewFileBase(fmt.Sprintf("/var/lib/omogen/submissions/%s", id))
	root.Gid = users.OmogenClientsID()
	root.GroupWritable = true
	if err := root.Mkdir("."); err != nil {
		return 0, fmt.Errorf("Could not create submission directory: %v", err)
	}
	path, err := root.FullPath("compiled")
	if err != nil {
		panic(err)
	}
	resp, err := runner.Compile(ctx, &runpb.CompileRequest{
		Program:    submission.Submission,
		OutputPath: path,
	})
	if err != nil {
		return 0, err
	}
	if resp.Program == nil {
		reporter.Err("compilation of submission failed: %v", resp.CompilationError)
		return 0, nil
	}
	var groups []*runpb.TestGroup
	for _, g := range problem.TestGroups {
		var cases []*runpb.TestCase
		for _, tc := range g.Tests {
			cases = append(cases,
				&runpb.TestCase{
					Name:       tc.Name,
					InputPath:  tc.InputPath,
					OutputPath: tc.OutputPath,
				})
		}
		groups = append(groups, &runpb.TestGroup{
			Cases: cases,
			Score: g.Score,
		})
	}
	req := &runpb.EvaluateRequest{
		SubmissionId: id,
		Program:      resp.Program,
		MemLimitKb:   problem.Metadata.Limits.MemoryLimitKb,
		TimeLimitMs:  timelim,
		EvaluateAll:  true,
		Groups:       groups,
	}
	if outputValidator != nil {
		req.Validator = &runpb.CustomValidator{
			Program: outputValidator,
		}
	}

	stream, err := runner.Evaluate(ctx, req)
	if err != nil {
		return 0, err
	}
	observedVerdict := make(map[runpb.Verdict]bool)
	allowedVerdicts := make(map[runpb.Verdict]bool)
	for _, allowed := range submission.Constraint.AllowedFailures {
		allowedVerdicts[allowed] = true
	}
	groupIdx := 0
	tcIdx := 0
	time := int32(0)
	for {
		verdict, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
		switch res := verdict.Result.(type) {
		case *runpb.EvaluateResponse_TestCase:
			v := res.TestCase.Verdict
			observedVerdict[v] = true
			if !allowedVerdicts[v] {
				reporter.Err("Submission %s got forbidden verdict %s on %s", submission.Name, v, groups[groupIdx].Cases[tcIdx].Name)
			}
			tcIdx++
		case *runpb.EvaluateResponse_TestGroup:
			groupIdx++
			tcIdx = 0
		case *runpb.EvaluateResponse_Submission:
			if res.Submission.Score != submission.Constraint.ExpectedScore {
				reporter.Err("Submission %s expected score %d got %d", submission.Name, submission.Constraint.ExpectedScore, res.Submission.Score)
			}
			reporter.Info("Submission %s took time %f s, got score %d", submission.Name, float32(res.Submission.TimeUsageMs) / 1000.0, res.Submission.Score)
			time = res.Submission.TimeUsageMs
		default:
		}
	}
	for _, required := range submission.Constraint.RequiredFailures {
		if !observedVerdict[required] {
			reporter.Err("Submission %s required verdict %s, but did not cause it", submission.Name, required)
		}
	}
	return time, nil
}
