package problems

import (
	"golang.org/x/text/language"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

func verifyStatements(problem *toolspb.Problem, reporter util.Reporter) error {
	if len(problem.Statements) == 0 {
		reporter.Err("Problem had no statements")
	}
	for _, statement := range problem.Statements {
		if statement.Title == "" {
			reporter.Err("Statement for language %s had no title", statement.LanguageCode)
		}
		if statement.StatementHtml == "" {
			reporter.Err("Statement for language %s was empty", statement.LanguageCode)
		}
		if _, err := language.Parse(statement.LanguageCode); err != nil {
			reporter.Err("Invalid language statement: %v", err)
		}
	}
	return nil
}
