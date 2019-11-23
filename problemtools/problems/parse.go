package problems

import (
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
	parseReporter := util.NewReporter()
	statements, err := parseStatements(path, parseReporter)
	if err != nil {
		return nil, err
	}

	metadata, err := parseMetadata(path, parseReporter)
	if err != nil {
		return nil, err
	}

	testGroupMap, err := parseTestdata(path, parseReporter)
	if err != nil {
		return nil, err
	}
	testGroups := make([]*toolspb.TestGroup, 0, len(testGroupMap))
	for _, v := range testGroupMap {
		testGroups = append(testGroups, v)
	}

	outputValidator, err := parseOutputValidator(path, parseReporter)
	if err != nil {
		return nil, err
	}

	inputValidators, err := parseInputValidators(path, parseReporter)
	if err != nil {
		return nil, err
	}

	submissions, err := parseSubmissions(path, parseReporter)
	if err != nil {
		return nil, err
	}

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
		Infos:         parseReporter.Infos(),
		Warnings:      parseReporter.Warnings(),
		Errors:        parseReporter.Errors(),
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
