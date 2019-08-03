package models

type TestGroupMap map[int32]*TestGroup
type TestGroupList []*TestGroup

func (tl TestGroupList) AsMap() TestGroupMap {
	tm := make(TestGroupMap)
	for _, g := range tl {
		tm[g.TestGroupId] = g
	}
	return tm
}

type TestGroup struct {
	ProblemId        int32  `db:"problem_id"`
	TestGroupId      int32  `db:"problem_testgroup_id"`
	Name             string `db:"testgroup_name"`
	PublicVisibility bool   `db:"public_visibility"`
	Tests            []*TestCase
}

type TestCase struct {
	TestGroupId int32       `db:"problem_testgroup_id"`
	TestCaseId  int32       `db:"problem_testcase_id"`
	Name        string      `db:"testcase_name"`
	InputFile   *StoredFile `db:"input_file"`
	OutputFile  *StoredFile `db:"output_file"`
}
