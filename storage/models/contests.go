package models

import (
	"database/sql"
	"strings"
	"time"
)

// A Contest is a contest with a number of problems team can register for and compete in.
type Contest struct {
	ContestID int32 `db:"contest_id"`
	// Full name of the contest.
	Title string `db:"title"`
	// Short name of the contest, for use in e.g. URLs.
	ShortName string `db:"short_name"`
	// A host name for the contest. Web requests to this hostname will resolve to this contest.
	HostName        sql.NullString `db:"host_name"`
	StartTime       sql.NullTime   `db:"start_time"`
	FlexibleEndTime sql.NullTime   `db:"selection_window_end_time"`
	Duration        time.Duration  `db:"duration"`
	// Whether the scoreboard should be hidden for contestants during the contest.
	HiddenScoreboard bool `db:"hidden_scoreboard"`
	Problems         []*ContestProblem
}

func (c *Contest) Flexible() bool {
	return c.FlexibleEndTime.Valid
}

func (c *Contest) Started(team *Team) bool {
	if c.FlexibleEndTime.Valid {
		if team == nil {
			return false
		}
		return team.StartTime.Valid
	}
	return c.FullStart()
}

func (c *Contest) FullStart() bool {
	if c.StartTime.Valid {
		return c.StartTime.Time.Before(time.Now())
	} else {
		return false
	}
}

func (c *Contest) Over(team *Team) bool {
	if !c.FlexibleEndTime.Valid {
		return c.FullEndTime().Before(time.Now())
	}
	if !team.StartTime.Valid {
		// Not started yet.
		return false
	}
	return c.EndTime(team).Before(time.Now())
}

func (c *Contest) FullOver() bool {
	return c.FullEndTime().Before(time.Now())
}

func (c *Contest) EndTime(team *Team) time.Time {
	if c.FlexibleEndTime.Valid {
		return team.StartTime.Time.Add(c.Duration)
	}
	return c.FullEndTime()
}

func (c *Contest) FullEndTime() time.Time {
	if c.FlexibleEndTime.Valid {
		return c.FlexibleEndTime.Time
	}
	return c.StartTime.Time.Add(c.Duration)
}

func (c *Contest) ElapsedFor(t time.Time, team *Team) time.Duration {
	if team.StartTime.Valid {
		return t.Sub(team.StartTime.Time)
	}
	return t.Sub(c.StartTime.Time)
}

func (c *Contest) UntilStart() time.Duration {
	return c.StartTime.Time.Sub(time.Now())
}

func (c *Contest) UntilEnd(team *Team) time.Duration {
	return c.EndTime(team).Sub(time.Now())
}

func (c *Contest) UntilFullEnd() time.Duration {
	return c.FullEndTime().Sub(time.Now())
}

func (c *Contest) Within(time time.Time, team *Team) bool {
	if c.FlexibleEndTime.Valid {
		if !team.StartTime.Valid {
			return false
		}
		return !team.StartTime.Time.After(time) && !c.EndTime(team).Before(time)
	}
	return !c.StartTime.Time.After(time) && !c.FullEndTime().Before(time)
}

func (c *Contest) MaxScore() int32 {
	res := int32(0)
	for _, p := range c.Problems {
		res += p.Problem.CurrentVersion.MaxScore()
	}
	return res
}

func (c *Contest) CanSeeScoreboard(team *Team) bool {
	if c.HiddenScoreboard {
		return c.FullOver() || (team != nil && c.Over(team))
	}
	if !c.Started(team) {
		return false
	}
	return true
}

func (c *Contest) HasProblem(problemID int32) bool {
	for _, p := range c.Problems {
		if p.ProblemID == problemID {
			return true
		}
	}
	return false
}

// A ContestProblem is a problem with associated metadata that appears in a contest.
type ContestProblem struct {
	ContestID int32 `db:"contest_id"`
	ProblemID int32 `db:"problem_id"`
	Problem   *Problem
	Label     string
}

type Team struct {
	TeamID    int32          `db:"team_id"`
	ContestID int32          `db:"contest_id"`
	TeamName  sql.NullString `db:"team_name"`
	StartTime sql.NullTime   `db:"start_time"`
	Members   TeamMemberList
}

func (t *Team) DisplayName() string {
	if t.TeamName.Valid {
		return t.TeamName.String
	}
	var names []string
	for _, m := range t.Members {
		names = append(names, m.Account.Username)
	}
	return strings.Join(names, ", ")
}

func (t *Team) Link() string {
	return t.Members[0].Account.Link()
}

func (t *Team) MemberIDs() []int32 {
	var ids []int32
	for _, k := range t.Members {
		ids = append(ids, k.AccountID)
	}
	return ids
}

type TeamList []*Team

type TeamMember struct {
	TeamID    int32   `db:"team_id"`
	AccountID int32   `db:"account_id"`
	Account   Account `db:"account"`
}

type TeamMemberList []*TeamMember
