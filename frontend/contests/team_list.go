package contests

import (
	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/contests"
	"github.com/jsannemo/omogenjudge/storage/models"
)

type Params struct {
	Teams []*models.Team
}

// ListHandler is the request handler for the team list.
func ListHandler(r *request.Request) (request.Response, error) {
	if r.Context.Contest == nil {
		return request.Redirect(paths.Route(paths.Home)), nil
	}

	teams, err := contests.ListTeams(r.Request.Context(), contests.TeamFilter{ContestID: r.Context.Contest.ContestID})
	if err != nil {
		return nil, err
	}
	return request.Template("contests_team_list", &Params{teams}), nil
}
