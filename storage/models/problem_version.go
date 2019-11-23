package models

import (
	"database/sql"
	"fmt"
)

// A ProblemVersion contains the judging-specific information of a problem that affect submission evaluation.
type ProblemVersion struct {
	ProblemVersionID int32            `db:"problem_version_id"`
	ProblemID        int32            `db:"problem_id"`
	TimeLimMS        int32            `db:"time_limit_ms"`
	MemLimKB         int32            `db:"memory_limit_kb"`
	OutputValidator  *OutputValidator `db:"problem_output_validator"`
	TestGroups       []*TestGroup
}

// Samples returns the test cases that should be publicly visible for the problem.
func (p *ProblemVersion) Samples() []*TestCase {
	var samples []*TestCase
	for _, group := range p.TestGroups {
		if !group.PublicVisibility {
			continue
		}
		samples = append(samples, group.Tests...)
	}
	return samples
}

// TimeLimString returns a human formatted string of the time limit.
func (p *ProblemVersion) TimeLimString() string {
	return fmt.Sprintf("%.1g s", float64(p.TimeLimMS)/1000)
}

// TimeLimString returns a human formatted string of the memory limit.
func (p *ProblemVersion) MemLimString() string {
	return fmt.Sprintf("%.1g GB", float64(p.MemLimKB)/1000/1000)
}

// OutputValidator is a custom validator of the output from a submission problem.
type OutputValidator struct {
	ValidatorLanguageID sql.NullString     `db:"language_id"`
	ValidatorSourceZIP  *NilableStoredFile `db:"validator_source_zip"`
}

// Nil checks whether the given validator is absent.
func (ov *OutputValidator) Nil() bool {
	return ov.ValidatorSourceZIP.Nil()
}
