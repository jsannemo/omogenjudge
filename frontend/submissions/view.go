package submissions

import (
	"fmt"
	"strconv"

	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/problems"
	"github.com/jsannemo/omogenjudge/storage/submissions"

	"github.com/gorilla/mux"
)

type ViewParams struct {
	Problems   models.ProblemMap
	Submission *models.Submission
}

func ViewHandler(r *request.Request) (request.Response, error) {
	vars := mux.Vars(r.Request)
	subId, err := strconv.ParseInt(vars[paths.SubmissionIdArg], 10, 32)
	if err != nil {
		return request.BadRequest("Non-numeric ID"), nil
	}
	subs, err := submissions.ListSubmissions(r.Request.Context(),
		submissions.ListArgs{WithFiles: true, WithRun: true, WithGroupRuns: true},
		submissions.ListFilter{Submissions: &submissions.SubmissionFilter{SubmissionIDs: []int32{int32(subId)}}})
	if err != nil {
		return nil, err
	}
	if len(subs) == 0 {
		return request.NotFound(), nil
	}
	sub := subs[0]
	if sub.AccountID != r.Context.UserID {
		return request.NotFound(), nil
	}
	probs, err := problems.List(r.Request.Context(),
		problems.ListArgs{WithStatements: problems.StmtTitles, WithTests: problems.TestsGroups},
		problems.Problems(sub.ProblemID))
	if err != nil {
		return nil, err
	}
	if len(probs) != 1 {
		return request.Error(fmt.Errorf("could not find problem for submission")), nil
	}
	if !r.Context.CanSeeProblem(probs[0]) {
		return request.NotFound(), nil
	}
	return request.Template("submissions_view", &ViewParams{probs.AsMap(), sub}), nil
}
