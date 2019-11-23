package middleware

import (
	"fmt"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/users"
)

// readUser is a processor that stores the logged-in user account data in the request context.
func readUser(r *request.Request) (request.Response, error) {
	if r.Context.UserID != 0 {
		user, err := users.ListUsers(r.Request.Context(), users.ListArgs{}, users.ListFilter{AccountID: []int32{r.Context.UserID}})
		if err != nil {
			return nil, fmt.Errorf("could not retrieve current user: %v", err)
		} else if len(user) == 1 {
			r.Context.User = user[0]
		} else {
			return nil, fmt.Errorf("did not find current user %v in database", r.Context.UserID)
		}
	}
	return nil, nil
}

