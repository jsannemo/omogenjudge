package problems

import (
	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

// ParseProblem parses the problem at the given path into the API format.
func ParseProblem(path string) (*toolspb.ParseProblemResponse, error) {
	var errors []string
	var warnings []string

	statementReporter := util.NewReporter()
	statements, err := parseStatements(path, statementReporter)
	if err != nil {
		return nil, err
	}
	statementReporter.AddFailures(&errors, &warnings)

	metadataReporter := util.NewReporter()
	metadata, err := parseMetadata(path, metadataReporter)
	if err != nil {
		return nil, err
	}
	metadataReporter.AddFailures(&errors, &warnings)

	testgroupReporter := util.NewReporter()
	testgroups, err := parseTestdata(path, testgroupReporter)
	if err != nil {
		return nil, err
	}
	testgroupReporter.AddFailures(&errors, &warnings)

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
