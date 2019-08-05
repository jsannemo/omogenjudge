package problems

import (
	"regexp"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

var isProblemId = regexp.MustCompile(`^[a-z0-9]+$`).MatchString

func verifyMetadata(problem *toolspb.Problem, reporter util.Reporter) error {
	if problem.Metadata.ProblemId == "" {
		reporter.Err("Empty problem id")
	}

	if !isProblemId(problem.Metadata.ProblemId) {
		reporter.Err("Problem ID must consist of [a-z0-9]+")
	}

	if problem.Metadata.License == toolspb.License_LICENSE_UNSPECIFIED {
		reporter.Err("No license set")
	}

	if problem.Metadata.Author == "" {
		reporter.Err("No author set")
	}

	timeLimit := problem.Metadata.Limits.TimeLimitMs
	if 0 > timeLimit || timeLimit > 60*1000 {
		reporter.Err("Time limit out of bounds: %v", timeLimit)
	}

	memLimit := problem.Metadata.Limits.MemoryLimitKb
	if 0 > memLimit || memLimit > 5*1024*1025 {
		reporter.Err("Memory limit out of bounds: %v", memLimit)
	}
	return nil
}
