package problems

import (
	"context"
	"path/filepath"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

func verifyOutputValidator(ctx context.Context, path string, problem *toolspb.Problem, runner runpb.RunServiceClient, rep util.Reporter) error {
	if problem.OutputValidator == nil {
		return nil
	}
	resp, err := runner.Compile(ctx, &runpb.CompileRequest{
		Program:    problem.OutputValidator,
		OutputPath: filepath.Join(path, "output_validator_compiled"),
	})
	if err != nil {
		return err
	}
	if resp.Program == nil {
		rep.Err("compilation of output validator failed: %v", resp.CompilationError)
		return nil
	}
	return nil
}
