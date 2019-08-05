package problems

import (
	"regexp"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

var isTestGroupName = regexp.MustCompile(`^[a-z0-9]+$`).MatchString
var isTestCaseName = regexp.MustCompile(`^[a-z0-9\-_]+$`).MatchString

func verifyTestdata(problem *toolspb.Problem, reporter util.Reporter) error {
	for _, g := range problem.TestGroups {
		if len(g.Tests) == 0 {
			reporter.Err("Empty test group %v", g.Name)
		}
		if !isTestGroupName(g.Name) {
			reporter.Err("Invalid test group name: %v [a-z0-9]+", g.Name)
		}
		for _, tc := range g.Tests {
			if !isTestCaseName(tc.Name) {
				reporter.Err("Invalid test cas ename: %v [a-z0-9\\-_]", tc.Name)
			}
		}
	}
	return nil
}
