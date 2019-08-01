// Handler for logging out a user
package users

import (
	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
)

// LogoutHandler handles logout requests
func LogoutHandler(r *request.Request) (request.Response, error) {
	r.Context.UserId = 0
	return request.Redirect(paths.Route(paths.Home)), nil
}
