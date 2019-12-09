package problems

import (
	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/problems"

	"github.com/gorilla/mux"
)

type ViewParams struct {
	Problem *models.Problem
}

func ViewHandler(r *request.Request) (request.Response, error) {
	vars := mux.Vars(r.Request)
	problem, err := getProblemIfVisible(r, vars[paths.ProblemNameArg],
		problems.ListArgs{WithStatements: problems.StmtAll, WithTests: problems.TestsSamplesAndGroups})
	if err != nil {
		return nil, err
	}
	if problem == nil {
		return request.NotFound(), nil
	}
	return request.Template("problems_view", &ViewParams{problem}), nil
}
