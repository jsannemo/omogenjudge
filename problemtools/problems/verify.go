package problems

import (
	"context"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

func VerifyProblem(ctx context.Context, req *toolspb.VerifyProblemRequest, runner runpb.RunServiceClient) (*toolspb.VerifyProblemResponse, error) {
	problem := req.ProblemToVerify
	path := req.ProblemPath
	verifyReporter := util.NewReporter()

	if err := verifyStatements(problem, verifyReporter); err != nil {
		return nil, err
	}

	inputValidators, err := verifyInputValidators(ctx, path, problem, runner, verifyReporter)
	if err != nil {
		return nil, err
	}

	outputValidator, err := verifyOutputValidator(ctx, path, problem, runner, verifyReporter)
	if err != nil {
		return nil, err
	}

	verifyIncludedFiles(problem.IncludedFiles, verifyReporter)

	if err := verifyTestdata(ctx, problem, inputValidators, runner, verifyReporter); err != nil {
		return nil, err
	}

	if err := verifySubmissions(ctx, problem, outputValidator, runner, verifyReporter); err != nil {
		return nil, err
	}

	if err := verifyMetadata(problem, verifyReporter); err != nil {
		return nil, err
	}

	return &toolspb.VerifyProblemResponse{
		VerifiedProblem: problem,
		Infos:           verifyReporter.Infos(),
		Warnings:        verifyReporter.Warnings(),
		Errors:          verifyReporter.Errors(),
	}, nil
}
