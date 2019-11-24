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
	if !r.Context.Contest.Started() {
		return request.NotFound(), nil
	}
	vars := mux.Vars(r.Request)
	problems, err := problems.List(r.Request.Context(), problems.ListArgs{WithStatements: problems.StmtAll, WithTests: problems.TestsSamples}, problems.ListFilter{ShortName: vars[paths.ProblemNameArg]})
	if err != nil {
		return nil, err
	}
	if len(problems) == 0 {
		return request.NotFound(), nil
	}
	problem := problems[0]
	return request.Template("problems_view", &ViewParams{problem}), nil
}
