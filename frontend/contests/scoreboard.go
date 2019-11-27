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
		submissions.ListArgs{WithRun: true, WithGroupRuns: true},
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
	TgScores    map[int32]map[string]int32
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
	probs := make(map[int32]*models.Problem)
	for _, p := range contest.Problems {
		probs[p.ProblemID] = p.Problem
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
			TgScores:    make(map[int32]map[string]int32),
			Submissions: make(map[int32]int),
			Times:       make(map[int32]time.Duration),
		}
		for _, a := range t.Members {
			accTeam[a.AccountID] = sc[t.TeamID]
		}
		for _, p := range contest.Problems {
			sc[t.TeamID].TgScores[p.ProblemID] = make(map[string]int32)
		}
	}

	for i := len(subs) - 1; i >= 0; i-- {
		sub := subs[i]
		if !contest.Within(sub.Created) || sub.CurrentRun.Waiting() {
			continue
		}
		team := accTeam[sub.AccountID]
		team.Submissions[sub.ProblemID]++
		inc := false
		for _, tg := range sub.CurrentRun.GroupRuns {
			if tg.Score > team.TgScores[sub.ProblemID][tg.TestGroupName] {
				inc = true
				team.TgScores[sub.ProblemID][tg.TestGroupName] = tg.Score
			}
		}
		if !inc && team.Submissions[sub.ProblemID] != 1 {
			continue
		}
		team.Times[sub.ProblemID] = sub.Created.Sub(contest.StartTime.Time)
	}
	var rankedTeams []*scoreboardTeam
	for _, team := range sc {
		for pid, p := range team.TgScores {
			ppoints := int32(0)
			for _, g := range p {
				ppoints += g
			}
			team.Scores[pid] = ppoints
		}
		for _, v := range team.Scores {
			team.TotalScore += v
		}
		rankedTeams = append(rankedTeams, team)
	}
	sort.Slice(rankedTeams, func(i, j int) bool {
		if rankedTeams[i].TotalScore != rankedTeams[j].TotalScore {
			return rankedTeams[i].TotalScore > rankedTeams[j].TotalScore
		}
		return rankedTeams[i].Team.DisplayName() < rankedTeams[j].Team.DisplayName()
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
