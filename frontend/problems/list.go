package problems

import (
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/problems"
)

type Params struct {
	Problems []*models.Problem
}

// Request handler for listing problem
func ListHandler(r *request.Request) (request.Response, error) {
	problems := problems.List(r.Request.Context(), problems.ListArgs{WithStatements: problems.StmtTitles}, problems.ListFilter{})
	return request.Template("problems_list", &Params{problems}), nil
}
