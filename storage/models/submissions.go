package models

import (
	"database/sql"
	"time"

	"github.com/jsannemo/omogenjudge/frontend/paths"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

// A Submission is a user submission to a problem.
type Submission struct {
	// SubmissionID The numeric ID of the submission. This may exposed externally.
	SubmissionID int32 `db:"submission_id"`
	// The ID of the author of the submission.
	AccountID int32 `db:"account_id"`
	// The ID of the problem the submission was for.
	ProblemID int32 `db:"problem_id"`
	// The tag of the langauge the submission was made in.
	Language string
	// The creation date of the submission.
	Created time.Time `db:"date_created"`
	// The files the submission consists of.
	Files []*SubmissionFile

	Runs []*SubmissionRun
}

// ToRunnerProgram serializes submission to a program that can be compiled by the RunService.
func (s *Submission) ToRunnerProgram() *runpb.Program {
	var files []*runpb.SourceFile
	for _, file := range s.Files {
		files = append(files, &runpb.SourceFile{
			Path:     file.Path,
			Contents: file.Contents,
		})
	}
	return &runpb.Program{
		Sources:    files,
		LanguageId: s.Language,
	}
}

// Link returns the link to the details of the submission.
func (s *Submission) Link() string {
	return paths.Route(paths.Submission, paths.SubmissionIdArg, s.SubmissionID)
}

// A SubmissionFile is a file that is part of a user submission.
type SubmissionFile struct {
	SubmissionID int32 `db:"submission_id"`
	// The path that the file should be placed at, relative to the compilation working directory.
	Path     string `db:"file_path"`
	Contents string `db:"file_contents"`
}

// An Evaluation is a collection of the various judgements of a given submission.
// The evaluation may represent a single test case, a group or the evaluation of an entire submission.
type Evaluation struct {
	Score       int32   `db:"score"`
	TimeUsageMS int32   `db:"time_usage_ms"`
	Verdict     Verdict `db:"verdict"`
}

// A SubmissionRun is a particular judge execution of a submission. A submission may have multiple runs when rejudged, for
// example if a new version of the problem is installed.
type SubmissionRun struct {
	SubmissionID     int32 `db:"submission_id"`
	SubmissionRunID  int32 `db:"submission_run_id"`
	ProblemVersionID int32 `db:"problem_version_id"`
	Evaluation
	// The current judge workflow status of this run.
	Status Status
	// Whether the run is public or not. Runs can be private during a rejudge of a set of subimssions, for example.
	Public       bool
	CompileError sql.NullString `db:"compile_error"`
	Created      time.Time      `db:"date_created"`
}

// A TestGroupRun is a particular judge execution of a test group.
type TestGroupRun struct {
	SubmissionRunID int32 `db:"submission_run_id"`
	TestGroupID     int32 `db:"problem_testgroup_id"`
	Evaluation
	Created time.Time `db:"date_created"`
}

// A TestCaseRun is a particular judge execution of a test case.
type TestCaseRun struct {
	SubmissionRunID int32 `db:"submission_run_id"`
	TestCaseID      int32 `db:"problem_testcase_id"`
	Evaluation
	Created time.Time `db:"date_created"`
}
