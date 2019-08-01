package problems

import (
	"net/http"

	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/storage/problems"
	"github.com/jsannemo/omogenjudge/storage/submissions"
	"github.com/jsannemo/omogenjudge/storage/models"

	"github.com/gorilla/mux"
)

type SubmitParams struct {
	Problem *models.Problem
}

func SubmitHandler(r *request.Request) (request.Response, error) {
	loginUrl := paths.Route(paths.Login)
  // TODO save current page location
	if r.Context.UserId == 0 {
		return request.Redirect(loginUrl), nil
	}

	vars := mux.Vars(r.Request)
  problems := problems.List(r.Request.Context(), problems.ListArgs{WithStatements: problems.StmtTitles}, problems.ListFilter{ShortName: vars[paths.ProblemNameArg]})
  if len(problems) == 0 {
    return request.NotFound(), nil
  }
  problem := problems[0]

	if r.Request.Method == http.MethodPost {
		submit := r.Request.FormValue("submission")
    s := &models.Submission{
      AccountId: r.Context.UserId,
      ProblemId: problem.ProblemId,
      Files: []*models.SubmissionFile{
        &models.SubmissionFile{
          Path: "main.cpp",
          Contents: submit,
        },
      },
    }
    err := submissions.Create(r.Request.Context(), s)
    if err != nil {
      return request.Error(err), nil
    }
    subUrl := paths.Route(paths.Submission, paths.SubmissionIdArg, s.SubmissionId)
    return request.Redirect(subUrl), nil
  }

	return request.Template("problems_submit", &SubmitParams{problem}), nil
}
