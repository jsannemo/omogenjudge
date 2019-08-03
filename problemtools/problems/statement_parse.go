package problems

import (
	"io/ioutil"
	"os"
	"path/filepath"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	util "github.com/jsannemo/omogenjudge/problemtools/util"
)

type statementFilterFunc func(string) (bool, error)
type statementParseFunc func(string, util.Reporter) (*toolspb.ProblemStatement, error)

type statementParser struct {
	name   string
	filter statementFilterFunc
	parser statementParseFunc
}

var parsers = []*statementParser{
	markdownParser(),
}

func parseStatements(path string, reporter util.Reporter) ([]*toolspb.ProblemStatement, error) {
	statements := make([]*toolspb.ProblemStatement, 0)
	statementPath := filepath.Join(path, "statements")

	if _, err := os.Stat(statementPath); os.IsNotExist(err) {
		reporter.Err("Problem had no statement folder")
		return statements, nil
	}

	files, err := ioutil.ReadDir(statementPath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			statement, err := parseStatement(filepath.Join(statementPath, f.Name()), reporter)
			if err != nil {
				return nil, err
			}
			// Statement can be null due to a reporter error that is not a run-time error
			if statement != nil {
				statements = append(statements, statement)
			}
		}
	}
	return statements, nil
}

func parseStatement(path string, reporter util.Reporter) (*toolspb.ProblemStatement, error) {
	var found = false
	var parsedStatement *toolspb.ProblemStatement
	for _, p := range parsers {
		match, err := p.filter(path)
		if err != nil {
			return nil, err
		}
		if match {
			if found {
				reporter.Err("Statement matched multiple parsers")
				break
			}
			found = true
			parsedStatement, err = p.parser(path, reporter)
		}
	}
	return parsedStatement, nil
}
