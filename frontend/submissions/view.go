package submissions

import (
	"strconv"

	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/problems"
	"github.com/jsannemo/omogenjudge/storage/submissions"

	"github.com/gorilla/mux"
)

type ViewParams struct {
	Problem    *models.Problem
	Submission *models.Submission
}

func ViewHandler(r *request.Request) (request.Response, error) {
	vars := mux.Vars(r.Request)
	subId, err := strconv.Atoi(vars[paths.SubmissionIdArg])
	if err != nil {
		return request.BadRequest("Non-numeric ID"), nil
	}
	subs := submissions.List(r.Request.Context(), submissions.ListArgs{WithFiles: true}, submissions.ListFilter{SubmissionId: int32(subId)})
	if len(subs) == 0 {
		return request.NotFound(), nil
	}
	sub := subs[0]
	prob := problems.List(r.Request.Context(), problems.ListArgs{WithStatements: problems.StmtTitles}, problems.ListFilter{ProblemId: sub.ProblemId})[0]
	return request.Template("submissions_view", &ViewParams{prob, sub}), nil
}
