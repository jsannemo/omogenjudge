// Overall handler for problem parsing.
package parser

import (
	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/metadata"
	"github.com/jsannemo/omogenjudge/problemtools/statement"
	"github.com/jsannemo/omogenjudge/problemtools/testdata"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

// ParseProblem parses the problem at the given path into the API format.
func ParseProblem(path string) (*toolspb.ParseProblemResponse, error) {
	var errors []string
	var warnings []string

	statementReporter := util.NewReporter()
	statements, err := statement.ParseStatements(path, statementReporter)
	if err != nil {
		return nil, err
	}
	errors, warnings = statementReporter.AddFailures(errors, warnings)

	metadataReporter := util.NewReporter()
	metadata, err := metadata.ParseMetadata(path, metadataReporter)
	if err != nil {
		return nil, err
	}
	errors, warnings = metadataReporter.AddFailures(errors, warnings)

	testgroupReporter := util.NewReporter()
	testgroups, err := testdata.ParseTestdata(path, testgroupReporter)
	if err != nil {
		return nil, err
	}
	errors, warnings = testgroupReporter.AddFailures(errors, warnings)

	problem := &toolspb.Problem{
		Statements: statements,
		Metadata:   metadata,
		TestGroups: testgroups,
	}
	return &toolspb.ParseProblemResponse{
		ParsedProblem: problem,
		Errors:        errors,
		Warnings:      warnings,
	}, nil
}
