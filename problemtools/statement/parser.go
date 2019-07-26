// Problem statement parsing.
package statement

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/logger"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	util "github.com/jsannemo/omogenjudge/problemtools/util"
)

type parseFilterFunc func(string) (bool, error)
type parseFunc func(string, util.Reporter) (*toolspb.ProblemStatement, error)

type parser struct {
	name   string
	filter parseFilterFunc
	parser parseFunc
}

var parsers = make([]*parser, 0)

func ParseStatements(path string, reporter util.Reporter) ([]*toolspb.ProblemStatement, error) {
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

// parseStatement parses a problem statement from a folder with a matching statement parser
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

func registerParser(parser *parser) {
	logger.Infof("Registering parser: %v", parser.name)
	parsers = append(parsers, parser)
}
