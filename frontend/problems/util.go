package problems

import (
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/problems"
)

func getProblemIfVisible(req *request.Request, shortname string, args problems.ListArgs) (*models.Problem, error) {
	probs, err := problems.List(req.Request.Context(), args, problems.ListFilter{ShortName: shortname})
	if err != nil {
		return nil, err
	}
	if len(probs) == 0 {
		return nil, nil
	}
	problem := probs[0]
	if !req.Context.CanSeeProblem(problem) {
		return nil, nil
	}
	return problem, nil

}
