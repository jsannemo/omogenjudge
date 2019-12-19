package home

import (
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/problems"
	"github.com/jsannemo/omogenjudge/storage/submissions"
)

func HomeHandler(r *request.Request) (request.Response, error) {
	if r.Context.Contest != nil {
		return contestHome(r)
	}
	return mainHome(r)
}

type homeParams struct {
	Submissions submissions.SubmissionList
	Problems    models.ProblemMap
}

func mainHome(r *request.Request) (request.Response, error) {
	subs, err := submissions.ListSubmissions(r.Request.Context(), submissions.ListArgs{
		WithRun:      true,
		WithAccounts: true,
	}, submissions.ListFilter{OnlyAvailable: true})
	if err != nil {
		return nil, err
	}
	probs, err := problems.List(r.Request.Context(), problems.ListArgs{}, problems.Problems(subs.ProblemIDs()...))
	if err != nil {
		return nil, err
	}
	return request.Template("home_home", &homeParams{Submissions: subs, Problems: probs.AsMap()}), nil
}

type problemData struct {
	Scores map[string]int32
	Groups []*models.TestGroup
	Score  int32
}

type ContestHomeParams struct {
	Problems map[int32]*problemData
}

func contestHome(r *request.Request) (request.Response, error) {
	team := r.Context.Team
	contest := r.Context.Contest
	if team != nil && contest.Started(team) {
		var probIDs []int32
		points := make(map[int32]*problemData)
		for _, p := range contest.Problems {
			probIDs = append(probIDs, p.ProblemID)
			points[p.ProblemID] = &problemData{
				Groups: p.Problem.CurrentVersion.TestGroups,
				Scores: make(map[string]int32),
			}
		}
		subs, err := submissions.ListSubmissions(r.Request.Context(), submissions.ListArgs{
			WithRun:       true,
			WithGroupRuns: true,
		}, submissions.ListFilter{
			Users:    &submissions.UserFilter{team.MemberIDs()},
			Problems: &submissions.ProblemFilter{probIDs}})
		if err != nil {
			return nil, err
		}
		for _, s := range subs {
			if s.CurrentRun.Waiting() {
				continue
			}
			if !r.Context.Contest.Within(s.Created, team) {
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
		return request.Template("home_contest", ContestHomeParams{points}), nil
	}
	return request.Template("home_contest", nil), nil
}
