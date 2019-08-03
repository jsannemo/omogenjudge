package middleware

import (
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/users"
)

func readUser(r *request.Request) (request.Response, error) {
	if r.Context.UserId != 0 {
		user, err := users.Get(r.Request.Context(), r.Context.UserId)
		if err == users.ErrNoSuchUser {
			r.Context.UserId = 0
		} else if err != nil {
			return nil, err
		} else {
			r.Context.User = user
		}
	}
	return nil, nil
}
