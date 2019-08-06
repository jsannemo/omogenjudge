package problems

import (
	"context"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

func VerifyProblem(ctx context.Context, req *toolspb.VerifyProblemRequest, runner runpb.RunServiceClient) (*toolspb.VerifyProblemResponse, error) {
	var errors []string
	var warnings []string

	problem := req.ProblemToVerify
	path := req.ProblemPath

	statementReporter := util.NewReporter()
	if err := verifyStatements(problem, statementReporter); err != nil {
		return nil, err
	}
	statementReporter.AddFailures(&errors, &warnings)

	metadataReporter := util.NewReporter()
	if err := verifyMetadata(problem, metadataReporter); err != nil {
		return nil, err
	}
	metadataReporter.AddFailures(&errors, &warnings)

	inputValidatorReporter := util.NewReporter()
	inputValidators, err := verifyInputValidators(ctx, path, problem, runner, inputValidatorReporter)
	if err != nil {
		return nil, err
	}
	inputValidatorReporter.AddFailures(&errors, &warnings)

	outputReporter := util.NewReporter()
	if err := verifyOutputValidator(ctx, path, problem, runner, outputReporter); err != nil {
		return nil, err
	}
	outputReporter.AddFailures(&errors, &warnings)

	testgroupReporter := util.NewReporter()
	if err := verifyTestdata(ctx, problem, inputValidators, runner, testgroupReporter); err != nil {
		return nil, err
	}
	testgroupReporter.AddFailures(&errors, &warnings)

	return &toolspb.VerifyProblemResponse{
		VerifiedProblem: problem,
		Errors:          errors,
		Warnings:        warnings,
	}, nil
}
