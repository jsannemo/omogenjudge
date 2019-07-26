package problems

import (
	"html/template"
	"net/http"

	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/storage/problems"
	"github.com/jsannemo/omogenjudge/storage/submissions"

	"github.com/gorilla/mux"
)

var submitTemplates = template.Must(template.ParseFiles(
	"frontend/problems/submit.tpl",
	"frontend/templates/header.tpl",
	"frontend/templates/nav.tpl",
	"frontend/templates/footer.tpl",
))
type SubmitParams struct {
	Problem *problems.Problem
}

func SubmitHandler(r *request.Request) (request.Response, error) {
	loginUrl, err := paths.Route(paths.Login).URL()
	if err != nil {
		return nil, err
	}
  // TODO save current page location
	if r.Context.UserId == 0 {
		return request.Redirect(loginUrl.String()), nil
	}

	vars := mux.Vars(r.Request)
	problem, err := problems.GetProblem(r.Request.Context(), vars[paths.ProblemNameArg], true)
	if err != nil {
		return request.Error(err), nil
	}

	if r.Request.Method == http.MethodPost {
		submit := r.Request.FormValue("submission")
    s := &submissions.Submission{
      AccountId: r.Context.UserId,
      ProblemId: problem.ProblemId,
      Files: []*submissions.SubmissionFile{
        &submissions.SubmissionFile{
          Path: "main.cpp",
          Contents: []byte(submit),
        },
      },
    }
    err := submissions.CreateSubmission(r.Request.Context(), s)
    if err != nil {
      return request.Error(err), nil
    }
  }

	return request.Template(submitTemplates, "page", &SubmitParams{problem}), nil
}
