package models

// A TestGroup is a group of test cases that belong together for evaluation purposes.
type TestGroup struct {
	TestGroupID      int32 `db:"problem_testgroup_id"`
	ProblemVersionID int32 `db:"problem_version_id"`
	// A user-visible name for the test group.
	Name string `db:"testgroup_name"`
	// Whether the tests in the group should be displayed to users.
	PublicVisibility bool `db:"public_visibility"`
	// The score that should be awarded if the test group is passed.
	Score int32 `db:"score"`
	// The flags that should be provided to output validators when evaluating the group.
	OutputValidatorFlags string `db:"output_validator_flags"`
	Tests                []*TestCase
}

// A TestCase represents a single test case that submissions can be evaluated on.
type TestCase struct {
	TestGroupID int32 `db:"problem_testgroup_id"`
	TestCaseID  int32 `db:"problem_testcase_id"`
	// A user-visible name for the test case.
	Name       string      `db:"testcase_name"`
	InputFile  *StoredFile `db:"input_file"`
	OutputFile *StoredFile `db:"output_file"`
}
