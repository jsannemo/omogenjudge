package problems

import (
	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/problems"
)

type Params struct {
	Problems []*models.Problem
}

// ListHandler is the request handler for the problem list.
func ListHandler(r *request.Request) (request.Response, error) {
	if !r.Context.Contest.FullStart() {
		return request.Redirect(paths.Route(paths.Home)), nil
	}

	probs, err := problems.List(r.Request.Context(), problems.ListArgs{WithStatements: problems.StmtTitles}, problems.ListFilter{})
	if err != nil {
		return nil, err
	}
	return request.Template("problems_list", &Params{probs}), nil
}
