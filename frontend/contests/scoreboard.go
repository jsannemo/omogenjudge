package contests

import (
	"sort"
	"time"

	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/contests"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/submissions"
)

// ScoreboardHandler is the request handler for the contest scoreboard.
func ScoreboardHandler(r *request.Request) (request.Response, error) {
	if !r.Context.Contest.Started() {
		return request.Redirect(paths.Route(paths.ContestTeams)), nil
	}

	teams, err := contests.ListTeams(r.Request.Context(), contests.TeamFilter{ContestID: r.Context.Contest.ContestID})
	if err != nil {
		return nil, err
	}
	var accountIDs []int32
	for _, t := range teams {
		for _, a := range t.Members {
			accountIDs = append(accountIDs, a.AccountID)
		}
	}

	var probIDs []int32
	for _, p := range r.Context.Contest.Problems {
		probIDs = append(probIDs, p.ProblemID)
	}
	subs, err := submissions.ListSubmissions(
		r.Request.Context(),
		submissions.ListArgs{WithRun: true},
		submissions.ListFilter{ProblemID: probIDs, UserID: accountIDs})
	if err != nil {
		return nil, err
	}

	scoreboard := makeScoreboard(teams, subs, r.Context.Contest)

	return request.Template("contest_scoreboard", scoreboard), nil
}

type scoreboardTeam struct {
	Team        *models.Team
	Rank        int
	Scores      map[int32]int32
	Submissions map[int32]int
	TotalScore  int32
	Times       map[int32]time.Duration
}

type scoreboardProblem struct {
	Label   string
	Problem *models.Problem
}

type scoreboard struct {
	Teams    []*scoreboardTeam
	Problems []*scoreboardProblem
	MaxScore int32
}

func makeScoreboard(teams models.TeamList, subs submissions.SubmissionList, contest *models.Contest) interface{} {
	maxScore := int32(0)
	var scp []*scoreboardProblem
	for _, p := range contest.Problems {
		maxScore += p.Problem.CurrentVersion.MaxScore()
		scp = append(scp, &scoreboardProblem{
			Label:   p.Label,
			Problem: p.Problem,
		})
	}
	sort.Slice(scp, func(i, j int) bool {
		return scp[i].Label < scp[j].Label
	})

	accTeam := make(map[int32]*scoreboardTeam)
	sc := make(map[int32]*scoreboardTeam)
	for _, t := range teams {
		sc[t.TeamID] = &scoreboardTeam{
			Team:        t,
			Scores:      make(map[int32]int32),
			Submissions: make(map[int32]int),
			Times:       make(map[int32]time.Duration),
		}
		for _, a := range t.Members {
			accTeam[a.AccountID] = sc[t.TeamID]
		}
	}

	for i := len(subs) - 1; i >= 0; i-- {
		sub := subs[i]
		if sub.Created.Before(contest.StartTime.Time) || sub.Created.After(contest.EndTime()) {
			continue
		}
		team := accTeam[sub.AccountID]
		if sub.CurrentRun.Score <= team.Scores[sub.ProblemID] {
			continue
		}
		team.Submissions[sub.ProblemID]++
		team.Scores[sub.ProblemID] = sub.CurrentRun.Score
		team.Times[sub.ProblemID] = sub.Created.Sub(contest.StartTime.Time)
	}
	var rankedTeams []*scoreboardTeam
	for _, team := range sc {
		for _, v := range team.Scores {
			team.TotalScore += v
		}
		rankedTeams = append(rankedTeams, team)
	}
	sort.Slice(rankedTeams, func(i, j int) bool {
		return rankedTeams[i].TotalScore > rankedTeams[j].TotalScore
	})
	cur := int32(-1)
	rank := -1
	for i := 0; i < len(rankedTeams); i++ {
		if rankedTeams[i].TotalScore != cur {
			rank = i + 1
		}
		rankedTeams[i].Rank = rank
		cur = rankedTeams[i].TotalScore
	}
	return &scoreboard{
		Teams:    rankedTeams,
		Problems: scp,
		MaxScore: maxScore,
	}
}
