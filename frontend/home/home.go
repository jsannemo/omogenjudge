package home

import (
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/submissions"
)

type problemData struct {
	Scores map[string]int32
	Groups []*models.TestGroup
	Score  int32
}

type HomeParams struct {
	Problems map[int32]*problemData
}

func HomeHandler(r *request.Request) (request.Response, error) {
	if r.Context.User != nil && r.Context.Contest.Started() {
		var probIDs []int32
		points := make(map[int32]*problemData)
		for _, p := range r.Context.Contest.Problems {
			probIDs = append(probIDs, p.ProblemID)
			points[p.ProblemID] = &problemData{
				Groups: p.Problem.CurrentVersion.TestGroups,
				Scores: make(map[string]int32),
			}
		}
		subs, err := submissions.ListSubmissions(r.Request.Context(), submissions.ListArgs{
			WithRun:       true,
			WithGroupRuns: true,
		}, submissions.ListFilter{UserID: []int32{r.Context.User.AccountID}, ProblemID: probIDs})
		if err != nil {
			return nil, err
		}
		for _, s := range subs {
			if s.CurrentRun.Waiting() {
				continue
			}
			if !r.Context.Contest.Within(s.Created) {
				continue
			}
			for _, run := range s.CurrentRun.GroupRuns {
				data := points[s.ProblemID]
				if run.Score > data.Scores[run.TestGroupName] {
					data.Scores[run.TestGroupName] = run.Score
				}
			}
		}
		for _, p := range points {
			for _, s := range p.Scores {
				p.Score += s
			}
		}
		return request.Template("home_home", HomeParams{points}), nil
	}
	return request.Template("home_home", nil), nil
}
