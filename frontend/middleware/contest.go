package middleware

import (
	"fmt"

	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/contests"
)

// readContest is a processor that stores the current contest data in the request context.
func readContest(r *request.Request) (request.Response, error) {
	hostname := r.Request.Host
	cs, err := contests.ListContests(r.Request.Context(), contests.ListArgs{WithProblems: true}, contests.ListFilter{HostName: hostname})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve current contest: %v", err)
	} else if len(cs) > 0 {
		r.Context.Contest = cs.Latest()
	} else {
		return request.NotFound(), nil
	}
	return nil, nil
}

// readTeam is a processor that stores the logged-in team in the request context.
func readTeam(r *request.Request) (request.Response, error) {
	if r.Context.Contest == nil || r.Context.User == nil {
		return nil, nil
	}
	teams, err := contests.ListTeams(r.Request.Context(), contests.TeamFilter{
		ContestID: r.Context.Contest.ContestID,
		AccountID: r.Context.User.AccountID,
	})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve current team: %v", err)
	} else if len(teams) > 0 {
		r.Context.Team = teams[0]
	}
	return nil, nil
}
