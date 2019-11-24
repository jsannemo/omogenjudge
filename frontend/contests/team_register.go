package contests

import (
	"net/http"

	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/contests"
	"github.com/jsannemo/omogenjudge/storage/models"
)

// RegisterHandler is the request handler for the team register action.
func RegisterHandler(r *request.Request) (request.Response, error) {
	if r.Context.User == nil {
		return request.Redirect(paths.Route(paths.Login)), nil
	}
	// User was already registered.
	if r.Context.Team != nil {
		return request.Redirect(paths.Route(paths.Home)), nil
	}
	// Don't allow registration after the contest ends.
	if r.Context.Contest.Over() {
		return request.Redirect(paths.Route(paths.Home)), nil
	}
	if r.Request.Method == http.MethodPost {
		err := contests.CreateTeam(r.Request.Context(), &models.Team{
			ContestID: r.Context.Contest.ContestID,
			Members:   []*models.TeamMember{&models.TeamMember{AccountID: r.Context.UserID}},
		})
		if err != nil {
			return nil, err
		}
		return request.Redirect(paths.Route(paths.Home)), nil
	}
	return request.Template("contest_team_register", nil), nil
}
