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
	contest := r.Context.Contest
	if contest != nil && !contest.Started(r.Context.Team) {
		return request.NotFound(), nil
	}
	vars := mux.Vars(r.Request)
	probs, err := problems.List(r.Request.Context(),
		problems.ListArgs{WithStatements: problems.StmtAll, WithTests: problems.TestsSamplesAndGroups},
		problems.ListFilter{ShortName: vars[paths.ProblemNameArg]})


	if err != nil {
		return nil, err
	}
	if len(probs) == 0 {
		return request.NotFound(), nil
	}
	problem := probs[0]

	if contest != nil && !contest.HasProblem(problem.ProblemID) {
		return request.NotFound(), nil
	}
	return request.Template("problems_view", &ViewParams{problem}), nil
}
