package home

import (
	"github.com/jsannemo/omogenjudge/frontend/request"
)

func HomeHandler(r *request.Request) (request.Response, error) {
	return request.Template("home_home", nil), nil
}
