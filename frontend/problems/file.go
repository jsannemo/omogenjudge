package problems

import (
	"database/sql"

	"github.com/gorilla/mux"

	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/problems"
)

func FileHandler(r *request.Request) (request.Response, error) {
	vars := mux.Vars(r.Request)
	shortName := vars[paths.ProblemNameArg]
	lang := vars[paths.ProblemLangArg]
	path := vars[paths.ProblemFileArg]
	problem, err := getProblemIfVisible(r, vars[paths.ProblemNameArg], problems.ListArgs{})
	if err != nil {
		return nil, err
	}
	if problem == nil {
		return request.NotFound(), nil
	}

	file, err := problems.GetStatementFile(r.Request.Context(), shortName, lang, path)
	if err == sql.ErrNoRows {
		return request.NotFound(), nil
	}
	if err != nil {
		return nil, err
	}
	content, err := file.Content.FileData()
	return request.RawBytes(content), err
}
