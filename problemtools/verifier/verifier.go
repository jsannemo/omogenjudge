// Overall problem verification
package verifier

import (
	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/metadata"
	"github.com/jsannemo/omogenjudge/problemtools/statement"
	"github.com/jsannemo/omogenjudge/problemtools/testdata"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

// VerifyProblem verifies a given problem that has already been parsed.
func VerifyProblem(problem *toolspb.Problem) (*toolspb.VerifyProblemResponse, error) {
	var errors []string
	var warnings []string

	statementReporter := util.NewReporter()
	if err := statement.VerifyStatements(problem, statementReporter); err != nil {
		return nil, err
	}
	errors, warnings = statementReporter.AddFailures(errors, warnings)

	metadataReporter := util.NewReporter()
	if err := metadata.VerifyMetadata(problem, metadataReporter); err != nil {
		return nil, err
	}
	errors, warnings = metadataReporter.AddFailures(errors, warnings)

	testgroupReporter := util.NewReporter()
	if err := testdata.VerifyTestdata(problem, testgroupReporter); err != nil {
		return nil, err
	}
	errors, warnings = testgroupReporter.AddFailures(errors, warnings)

	return &toolspb.VerifyProblemResponse{
		VerifiedProblem: problem,
		Errors:          errors,
		Warnings:        warnings,
	}, nil
}
