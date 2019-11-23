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
	Problems   problems.ProblemMap
	Submission *models.Submission
}

func ViewHandler(r *request.Request) (request.Response, error) {
	vars := mux.Vars(r.Request)
	subId, err := strconv.Atoi(vars[paths.SubmissionIdArg])
	if err != nil {
		return request.BadRequest("Non-numeric ID"), nil
	}
	subs, err := submissions.ListSubmissions(r.Request.Context(), submissions.ListArgs{WithFiles: true}, submissions.ListFilter{SubmissionID: []int32{int32(subId)}})
	if err != nil {
		return nil, err
	}
	if len(subs) == 0 {
		return request.NotFound(), nil
	}
	sub := subs[0]
	probs, err := problems.List(r.Request.Context(), problems.ListArgs{WithStatements: problems.StmtTitles}, problems.ListFilter{ProblemId: []int32{sub.ProblemID}})
	if err != nil {
		return nil, err
	}
	return request.Template("submissions_view", &ViewParams{probs.AsMap(), sub}), nil
}
