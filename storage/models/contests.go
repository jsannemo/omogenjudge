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
	HostName  sql.NullString `db:"host_name"`
	StartTime sql.NullTime   `db:"start_time"`
	Duration  time.Duration  `db:"duration"`
	// Whether the scoreboard should be hidden for contestants during the contest.
	HiddenScoreboard bool `db:"hidden_scoreboard"`
	Problems         []*ContestProblem
}

func (c *Contest) Started() bool {
	if c.StartTime.Valid {
		return c.StartTime.Time.Before(time.Now())
	} else {
		return false
	}
}

func (c *Contest) Over() bool {
	return c.EndTime().Before(time.Now())
}

func (c *Contest) EndTime() time.Time {
	return c.StartTime.Time.Add(c.Duration)
}

func (c *Contest) UntilStart() time.Duration {
	return c.StartTime.Time.Sub(time.Now())
}

func (c *Contest) UntilEnd() time.Duration {
	return c.EndTime().Sub(time.Now())
}

func (c *Contest) Within(time time.Time) bool {
	return !c.StartTime.Time.After(time) && !c.EndTime().Before(time)
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

type TeamList []*Team

type TeamMember struct {
	TeamID    int32   `db:"team_id"`
	AccountID int32   `db:"account_id"`
	Account   Account `db:"account"`
}

type TeamMemberList []*TeamMember
