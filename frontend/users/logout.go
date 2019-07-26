// Handler for logging out a user
package users

import (
	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
)

// LogoutHandler handles logout requests
func LogoutHandler(r *request.Request) (request.Response, error) {
	r.Context.UserId = 0
	root, err := paths.Route(paths.Home).URL()
	if err != nil {
		return nil, err
	}
	return request.Redirect(root.String()), nil
}
