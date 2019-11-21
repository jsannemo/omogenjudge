package problems

import (
	"github.com/jsannemo/omogenjudge/storage/models"
)

type ProblemMap map[int32]*models.Problem

func (p ProblemMap) AsList() ProblemList {
	var probs ProblemList
	for _, prob := range p {
		probs = append(probs, prob)
	}
	return probs
}

func (p ProblemMap) Ids() []int32 {
	var ids []int32
	for id, _ := range p {
		ids = append(ids, id)
	}
	return ids
}

type ProblemList []*models.Problem

func (pl ProblemList) AsMap() ProblemMap {
	pm := make(ProblemMap)
	for _, p := range pl {
		pm[p.ProblemId] = p
	}
	return pm
}

type StatementList []*models.ProblemStatement

func (sl StatementList) AddTo(pm ProblemMap) {
	for _, s := range sl {
		p := pm[s.ProblemId]
		p.Statements = append(p.Statements, s)
	}
}

type TestGroupMap map[int32]*models.TestGroup
type TestGroupList []*models.TestGroup

func (tl TestGroupList) AsMap() TestGroupMap {
	tm := make(TestGroupMap)
	for _, g := range tl {
		tm[g.TestGroupId] = g
	}
	return tm
}
