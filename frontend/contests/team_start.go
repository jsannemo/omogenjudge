package contests

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/contests"
)

// StartHandler is the request handler for the team register action.
func StartHandler(r *request.Request) (request.Response, error) {
	team := r.Context.Team
	if team == nil {
		return request.Redirect(paths.Route(paths.Home)), nil
	}
	if !r.Context.Contest.Flexible() || team.StartTime.Valid || r.Context.Contest.FullOver() {
		return request.Redirect(paths.Route(paths.Home)), nil
	}
	if r.Request.Method == http.MethodPost {
		team.StartTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		err := contests.UpdateTeam(r.Request.Context(), team)
		if err != nil {
			return nil, err
		}
		return request.Redirect(paths.Route(paths.Home)), nil
	}
	return request.Redirect(paths.Route(paths.Home)), nil
}
