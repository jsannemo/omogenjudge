package problems

import (
	"context"
	"fmt"
	"path/filepath"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

func verifyInputValidators(ctx context.Context, path string, problem *toolspb.Problem, runner runpb.RunServiceClient, rep util.Reporter) ([]*runpb.CompiledProgram, error) {
	var validators []*runpb.CompiledProgram
	for i, val := range problem.InputValidators {
		resp, err := runner.Compile(ctx, &runpb.CompileRequest{
			Program:    val,
			OutputPath: filepath.Join(path, fmt.Sprintf("input_validator_compiled_%d", i)),
		})
		if err != nil {
			return nil, err
		}
		if resp.Program == nil {
			rep.Err("compilation of input validator failed: %v", resp.CompilationError)
			return nil, nil
		}
		validators = append(validators, resp.Program)
	}
	return validators, nil
}
