package problems

import (
	"html/template"

	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/storage/problems"

	"github.com/gorilla/mux"
)

var viewTemplates = template.Must(template.ParseFiles(
	"frontend/problems/view.tpl",
	"frontend/templates/header.tpl",
	"frontend/templates/nav.tpl",
	"frontend/templates/footer.tpl",
))

type ViewParams struct {
	Problem *problems.Problem
}

func ViewHandler(r *request.Request) (request.Response, error) {
	vars := mux.Vars(r.Request)
  problems, err := problems.ListProblems(r.Request.Context(), problems.ListArgs{WithStatements: true, WithSamples: true}, problems.ListFilter{ShortName: vars[paths.ProblemNameArg]})
	if err != nil {
		return request.Error(err), nil
	}
  if len(problems) == 0 {
    return request.NotFound(), nil
  }
  problem := problems[0]
	return request.Template(viewTemplates, "page", &ViewParams{problem}), nil
}
