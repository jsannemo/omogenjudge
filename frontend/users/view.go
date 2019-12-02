package users

import (
	"github.com/gorilla/mux"

	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/problems"
	"github.com/jsannemo/omogenjudge/storage/submissions"
	"github.com/jsannemo/omogenjudge/storage/users"
)

type ListParams struct {
	Submissions submissions.SubmissionList
	Problems    models.ProblemMap
	Filtered    bool
	Username    string
}

func ViewHandler(r *request.Request) (request.Response, error) {
	if r.Context.Contest != nil && !r.Context.Contest.FullStart() {
		return request.NotFound(), nil
	}
	vars := mux.Vars(r.Request)
	userName := vars[paths.UserNameArg]
	user, err := users.ListUsers(r.Request.Context(), users.ListArgs{}, users.ListFilter{Username: userName})
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return request.NotFound(), err
	}
	userID := user.Single().AccountID
	subFilter := submissions.ListFilter{Users: &submissions.UserFilter{[]int32{userID}}}

	if r.Context.Contest != nil {
		var cProbs []int32
		for _, p := range r.Context.Contest.Problems {
			cProbs = append(cProbs, p.ProblemID)
		}
		subFilter.Problems = &submissions.ProblemFilter{cProbs}
	}
	subs, err := submissions.ListSubmissions(
		r.Request.Context(),
		submissions.ListArgs{WithRun: true},
		subFilter)
	if err != nil {
		return nil, err
	}
	probIDs := subs.ProblemIDs()

	probs, err := problems.List(r.Request.Context(),
		problems.ListArgs{WithStatements: problems.StmtTitles, WithTests: problems.TestsGroups},
		problems.Problems(probIDs...))
	if err != nil {
		return nil, err
	}
	return request.Template("users_view",
		&ListParams{Submissions: subs, Problems: probs.AsMap(), Filtered: userID != r.Context.UserID, Username: userName}), nil
}
