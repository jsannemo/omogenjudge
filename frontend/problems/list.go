// Handler for listing problems
package problems

import (
	"html/template"

	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/problems"
)

var listTemplates = template.Must(template.ParseFiles(
	"frontend/problems/list.tpl",
	"frontend/templates/header.tpl",
	"frontend/templates/nav.tpl",
	"frontend/templates/footer.tpl",
))

type Params struct {
	Problems problems.ProblemMap
}

// Request handler for listing problem
func ListHandler(r *request.Request) (request.Response, error) {
	problems, err := problems.ListProblems(r.Request.Context())
	if err != nil {
		return request.Error(err), nil
	}
	return request.Template(listTemplates, "page", &Params{problems}), nil
}
