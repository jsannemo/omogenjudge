// Handler for registering
package users

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/users"
)

var registerTemplates = template.Must(template.ParseFiles(
	"frontend/users/register.tpl",
	"frontend/templates/header.tpl",
	"frontend/templates/nav.tpl",
	"frontend/templates/footer.tpl",
))

// LogoutHandler handles new user requests
func RegisterHandler(r *request.Request) (request.Response, error) {
	rootUrl, err := paths.Route(paths.Home).URL()
	if err != nil {
		return nil, err
	}
	root := rootUrl.String()
	if r.Context.UserId != 0 {
		return request.Redirect(root), nil
	}
	if r.Request.Method == http.MethodPost {
		username := r.Request.FormValue("username")
		password := r.Request.FormValue("password")

		userid, err := users.CreateUser(r.Request.Context(), username, password)
		if err == users.ErrUserExists {
      // TODO show an error message for this instead on the registration page
			return request.Error(errors.New("Username in use")), nil
		} else if err != nil {
			return request.Error(err), nil
		}
		r.Context.UserId = userid
		return request.Redirect(root), nil
	}
	return request.Template(registerTemplates, "page", nil), nil
}
