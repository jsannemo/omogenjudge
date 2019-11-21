package models

type TestGroup struct {
	TestGroupId          int32  `db:"problem_testgroup_id"`
	ProblemVersionId     int32  `db:"problem_version_id"`
	Name                 string `db:"testgroup_name"`
	PublicVisibility     bool   `db:"public_visibility"`
	Score                int32  `db:"score"`
	OutputValidatorFlags string `db:"output_validator_flags"`
	Tests                []*TestCase
}

type TestCase struct {
	TestGroupId int32       `db:"problem_testgroup_id"`
	TestCaseId  int32       `db:"problem_testcase_id"`
	Name        string      `db:"testcase_name"`
	InputFile   *StoredFile `db:"input_file"`
	OutputFile  *StoredFile `db:"output_file"`
}
