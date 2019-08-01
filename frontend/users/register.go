package users

import (
	"errors"
	"net/http"

	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/users"
)

func RegisterHandler(r *request.Request) (request.Response, error) {
  root := paths.Route(paths.Home)
	if r.Context.UserId != 0 {
		return request.Redirect(root), nil
	}
	if r.Request.Method == http.MethodPost {
		username := r.Request.FormValue("username")
		password := r.Request.FormValue("password")
    user := &models.Account{
      Username: username,
    }
    user.SetPassword(password)

		err := users.Create(r.Request.Context(), user)
		if err == users.ErrUserExists {
      // TODO show an error message for this instead on the registration page
			return request.Error(errors.New("Username in use")), nil
		} else if err != nil {
			return nil, err
		}
		r.Context.UserId = user.AccountId
		return request.Redirect(root), nil
	}
	return request.Template("users_register", nil), nil
}
