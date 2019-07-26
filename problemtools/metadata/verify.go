package metadata

import (
	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

// VerifyMetadata verifies the metadata of a problem, reporting potential verification errors.
func VerifyMetadata(problem *toolspb.Problem, reporter util.Reporter) error {
	if problem.Metadata.ProblemId == "" {
		reporter.Err("Empty problem id")
	}
	return nil
}
