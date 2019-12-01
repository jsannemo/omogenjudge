package contests

import (
	"fmt"
	"html/template"
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
	if r.Context.Contest == nil {
		return request.NotFound(), nil
	}
	contest := r.Context.Contest
	team := r.Context.Team
	if !contest.CanSeeScoreboard(team) {
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

	scoreboard := makeScoreboard(teams, subs, r.Context.Contest, team)

	return request.Template("contest_scoreboard", scoreboard), nil
}

type scoreboardTeam struct {
	Team        *models.Team
	Rank        int
	Scores      map[int32]int32
	ScoreCols   map[int32]template.CSS
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

func makeScoreboard(teams models.TeamList, subs submissions.SubmissionList, contest *models.Contest, viewer *models.Team) interface{} {
	maxScore := contest.MaxScore()
	var scp []*scoreboardProblem
	probs := make(map[int32]*models.Problem)
	for _, p := range contest.Problems {
		probs[p.ProblemID] = p.Problem
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
			ScoreCols:   make(map[int32]template.CSS),
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
		submitter := accTeam[sub.AccountID]
		if sub.CurrentRun.Waiting() {
			continue
		}
		if contest.FlexibleEndTime.Valid {
			elapsed := contest.ElapsedFor(sub.Created, submitter.Team)
			if elapsed >= contest.Duration || elapsed < 0 {
				continue
			}
			if viewer != nil {
				now := contest.ElapsedFor(time.Now(), viewer)
				if elapsed > now {
					continue
				}
			}
		} else if sub.Created.After(contest.FullEndTime()) || sub.Created.Before(contest.StartTime.Time) {
			continue
		}
		submitter.Submissions[sub.ProblemID]++
		inc := false
		for _, tg := range sub.CurrentRun.GroupRuns {
			if tg.Score > submitter.TgScores[sub.ProblemID][tg.TestGroupName] {
				inc = true
				submitter.TgScores[sub.ProblemID][tg.TestGroupName] = tg.Score
			}
		}
		if !inc && submitter.Submissions[sub.ProblemID] != 1 {
			continue
		}
		submitter.Times[sub.ProblemID] = sub.Created.Sub(contest.StartFor(submitter.Team))
	}
	var rankedTeams []*scoreboardTeam
	for _, team := range sc {
		for pid, p := range team.TgScores {
			ppoints := int32(0)
			for _, g := range p {
				ppoints += g
			}
			team.Scores[pid] = ppoints
			team.ScoreCols[pid] = scoreCol(ppoints, probs[pid].CurrentVersion.MaxScore())
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

func scoreCol(score int32, maxScore int32) template.CSS {
	// background-color: hsl(111, 67%, 85%);
	// background-color: hsl(2, 100%, 95%);
	frac := float32(score) / float32(maxScore)
	h := int32(frac*(111-2) + 2)
	s := int32(frac*(67-100) + 100)
	l := int32(frac*(85-95) + 95)
	return template.CSS(fmt.Sprintf("hsl(%d, %d%%, %d%%)", h, s, l))
}
