package models

import (
	"database/sql"
	"fmt"
)

type ProblemVersion struct {
	ProblemVersionId int32            `db:"problem_version_id"`
	ProblemId        int32            `db:"problem_id"`
	TimeLimMs        int32            `db:"time_limit_ms"`
	MemLimKb         int32            `db:"memory_limit_kb"`
	OutputValidator  *OutputValidator `db:"problem_output_validator"`
	TestGroups       []*TestGroup
}

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

func (p *ProblemVersion) Tests() []*TestCase {
	var tests []*TestCase
	for _, group := range p.TestGroups {
		tests = append(tests, group.Tests...)
	}
	return tests
}

func (p *ProblemVersion) TestDataFiles() FileList {
	var files FileList
	for _, tc := range p.Tests() {
		files = append(files, tc.InputFile, tc.OutputFile)
	}
	return files
}

func (p *ProblemVersion) TimeLimString() string {
	return fmt.Sprintf("%.1g s", float64(p.TimeLimMs)/1000)
}

func (p *ProblemVersion) MemLimString() string {
	return fmt.Sprintf("%.1g GB", float64(p.MemLimKb)/1000/1000)
}

type OutputValidator struct {
	ValidatorLanguageId sql.NullString     `db:"language_id"`
	ValidatorSourceZip  *NilableStoredFile `db:"validator_source_zip"`
}

func (ov *OutputValidator) Nil() bool {
	return !ov.ValidatorSourceZip.NotNil()
}
