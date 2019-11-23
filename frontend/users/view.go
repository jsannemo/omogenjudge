package users

import (
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/problems"
	"github.com/jsannemo/omogenjudge/storage/submissions"
)

type ListParams struct {
	Submissions submissions.SubmissionList
	Problems    problems.ProblemMap
}

func ViewHandler(r *request.Request) (request.Response, error) {
	userId := r.Context.UserID
	subs, err := submissions.ListSubmissions(r.Request.Context(), submissions.ListArgs{}, submissions.ListFilter{UserID: userId})
	if err != nil {
		return nil, err
	}
	// TODO(jsannemo) add filter
	probs, err := problems.List(r.Request.Context(), problems.ListArgs{WithStatements: problems.StmtTitles}, problems.ListFilter{})
	if err != nil {
		return nil, err
	}
	return request.Template("users_view", &ListParams{subs, probs.AsMap()}), nil
}
