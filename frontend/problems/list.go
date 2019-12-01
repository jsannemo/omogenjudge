package problems

import (
	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/models"
)

type Params struct {
	Problems []*models.Problem
}

// ListHandler is the request handler for the problem list.
func ListHandler(r *request.Request) (request.Response, error) {
	return request.Redirect(paths.Route(paths.Home)), nil
}
