package problems

import (
	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

func VerifyProblem(problem *toolspb.Problem) (*toolspb.VerifyProblemResponse, error) {
	var errors []string
	var warnings []string

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

	testgroupReporter := util.NewReporter()
	if err := verifyTestdata(problem, testgroupReporter); err != nil {
		return nil, err
	}
	testgroupReporter.AddFailures(&errors, &warnings)

	return &toolspb.VerifyProblemResponse{
		VerifiedProblem: problem,
		Errors:          errors,
		Warnings:        warnings,
	}, nil
}
