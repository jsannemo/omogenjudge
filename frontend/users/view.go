package users

import (
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/submissions"
	"github.com/jsannemo/omogenjudge/storage/problems"
	"github.com/jsannemo/omogenjudge/storage/models"
)

type ListParams struct {
	Submissions models.SubmissionList
  Problems models.ProblemMap
}

func ViewHandler(r *request.Request) (request.Response, error) {
  userId := r.Context.UserId
  subs := submissions.List(r.Request.Context(), submissions.ListArgs{}, submissions.ListFilter{UserId: userId})
  probs := problems.List(r.Request.Context(), problems.ListArgs{WithStatements: problems.StmtTitles}, problems.ListFilter{Submissions: subs}).AsMap()
	return request.Template("users_view", &ListParams{subs, probs}), nil
}
