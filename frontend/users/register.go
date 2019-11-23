package users

import (
	"net/http"

	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/users"
)

type RegisterParams struct {
	Error string
}

func RegisterHandler(r *request.Request) (request.Response, error) {
	root := paths.Route(paths.Home)
	if r.Context.UserID != 0 {
		return request.Redirect(root), nil
	}
	if r.Request.Method == http.MethodPost {
		username := r.Request.FormValue("username")
		password := r.Request.FormValue("password")
		user := &models.Account{
			Username: username,
		}
		user.SetPassword(password)

		err := users.CreateUser(r.Request.Context(), user)
		if err == users.ErrUserExists {
			return request.Template("users_register", &RegisterParams{Error: "Användarnamnet är upptaget"}), nil
		} else if err != nil {
			return nil, err
		}
		r.Context.UserID = user.AccountID
		return request.Redirect(root), nil
	}
	return request.Template("users_register", &RegisterParams{}), nil
}
