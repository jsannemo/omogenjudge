package problems

import (
	"io/ioutil"
	"os"
	"path/filepath"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	util "github.com/jsannemo/omogenjudge/problemtools/util"
)

type statementFilterFunc func(path string) (bool, error)
type statementParseFunc func(path string, problemName string, statementFiles map[string]string, reporter util.Reporter) (*toolspb.ProblemStatement, error)

type statementParser struct {
	name   string
	filter statementFilterFunc
	parser statementParseFunc
}

var parsers = []*statementParser{
	markdownParser(),
}

func parseStatements(path string, reporter util.Reporter) (*toolspb.ProblemStatements, error) {
	statements := &toolspb.ProblemStatements{
		StatementFiles: make(map[string]string),
		Attachments:    make(map[string]string),
	}
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
			err := parseStatement(statements, filepath.Join(statementPath, f.Name()), filepath.Base(path), reporter)
			if err != nil {
				return nil, err
			}
		}
	}
	attachmentPath := filepath.Join(path, "attachments")
	attachements, err := ioutil.ReadDir(attachmentPath)
	for _, f := range attachements {
		if !f.IsDir() {
			statements.Attachments[f.Name()] = filepath.Join(attachmentPath, f.Name())
		} else {
			// TODO(jsannemo): support zipping of directories
		}
	}

	return statements, nil
}

func parseStatement(statements *toolspb.ProblemStatements, path string, problemName string, reporter util.Reporter) error {
	var found = false
	var parsedStatement *toolspb.ProblemStatement
	for _, p := range parsers {
		match, err := p.filter(path)
		if err != nil {
			return err
		}
		if match {
			if found {
				reporter.Err("Statement matched multiple parsers")
				break
			}
			found = true
			parsedStatement, err = p.parser(path, problemName, statements.StatementFiles, reporter)
		}
	}
	statements.Statements = append(statements.Statements, parsedStatement)
	return nil
}
