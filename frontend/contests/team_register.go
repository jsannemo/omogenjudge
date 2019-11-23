package contests

import (
	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/contests"
	"github.com/jsannemo/omogenjudge/storage/models"
)

type Params struct {
}

// ListHandler is the request handler for the team list.
func ListHandler(r *request.Request) (request.Response, error) {
	if r.Context.Contest == nil {
		return request.Redirect(paths.Route(paths.Home)), nil
	}
	if r.Context.User == nil {
		return request.Redirect(paths.Route(paths.Login)), nil
	}
	teams, err := contests.CreateTeam(r.Request.Context(), contests.TeamFilter{ContestID: r.Context.Contest.ContestID})
	if err != nil {
		return nil, err
	}
	return request.Redirect(paths.Route(paths.Home)), nil
}
