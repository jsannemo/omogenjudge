// Handler for logging in a user
package users

import (
	"errors"
	"net/http"

	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/users"
)

// LoginHandler handles login requests
func LoginHandler(r *request.Request) (request.Response, error) {
	root := paths.Route(paths.Home)
	if r.Context.UserId != 0 {
		return request.Redirect(root), nil
	}
	if r.Request.Method == http.MethodPost {
		username := r.Request.FormValue("username")
		password := r.Request.FormValue("password")

		user, err := users.Authenticate(r.Request.Context(), username, password)
		if err == users.ErrInvalidLogin {
      // TODO show an error message for this instead on the registration page
			return request.Error(errors.New("Incorrect login details")), nil
		} else if err != nil {
			return nil, err
		}
		r.Context.UserId = user.AccountId
		return request.Redirect(root), nil
	}
	return request.Template("users_login", nil), nil
}
