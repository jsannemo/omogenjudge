package problems

import (
	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

func verifyMetadata(problem *toolspb.Problem, reporter util.Reporter) error {
	if problem.Metadata.ProblemId == "" {
		reporter.Err("Empty problem id")
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
