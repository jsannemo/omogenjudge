package problems

import (
	"github.com/google/logger"
	"io/ioutil"
	"os"
	"path/filepath"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
	"github.com/jsannemo/omogenjudge/runner/language"
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
	testGroupMap, err := parseTestdata(path, testgroupReporter)
	if err != nil {
		return nil, err
	}
	testgroupReporter.AddFailures(&errors, &warnings)
	testGroups := make([]*toolspb.TestGroup, 0, len(testGroupMap))
	for _, v := range testGroupMap {
		testGroups = append(testGroups, v)
	}

	outputValidatorReporter := util.NewReporter()
	outputValidator, err := parseOutputValidator(path, outputValidatorReporter)
	if err != nil {
		return nil, err
	}
	outputValidatorReporter.AddFailures(&errors, &warnings)

	inputValidatorReporter := util.NewReporter()
	inputValidators, err := parseInputValidators(path, inputValidatorReporter)
	if err != nil {
		return nil, err
	}
	inputValidatorReporter.AddFailures(&errors, &warnings)

	submissionReporter := util.NewReporter()
	submissions, err := parseSubmissions(path, submissionReporter)
	if err != nil {
		return nil, err
	}
	submissionReporter.AddFailures(&errors, &warnings)

	logger.Infof("submissions: %v", submissions)

	problem := &toolspb.Problem{
		Statements:      statements,
		Metadata:        metadata,
		TestGroups:      testGroups,
		OutputValidator: outputValidator,
		InputValidators: inputValidators,
		Submissions:     submissions,
	}
	return &toolspb.ParseProblemResponse{
		ParsedProblem: problem,
		Errors:        errors,
		Warnings:      warnings,
	}, nil
}

func parseProgram(mpath string, dir bool, reporter util.Reporter) (*runpb.Program, error) {
	program := &runpb.Program{}
	if dir {
		err := filepath.Walk(mpath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			dat, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			p, err := filepath.Rel(mpath, path)
			if err != nil {
				return err
			}
			program.Sources = append(program.Sources, &runpb.SourceFile{
				Path:     p,
				Contents: string(dat),
			})
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		dat, err := ioutil.ReadFile(mpath)
		if err != nil {
			return nil, err
		}
		program.Sources = append(program.Sources, &runpb.SourceFile{
			Path:     filepath.Base(mpath),
			Contents: string(dat),
		})
	}
	err := language.GuessLanguage(program)
	if err != nil {
		reporter.Err("Failed guessing language of program %v: %v", program, err)
		return nil, nil
	}
	return program, nil
}
