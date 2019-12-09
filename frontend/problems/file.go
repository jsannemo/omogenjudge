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
	path := vars[paths.ProblemFileArg]
	problem, err := getProblemIfVisible(r, vars[paths.ProblemNameArg], problems.ListArgs{})
	if err != nil {
		return nil, err
	}
	if problem == nil {
		return request.NotFound(), nil
	}

	file, err := problems.GetStatementFile(r.Request.Context(), shortName, path)
	if err == sql.ErrNoRows {
		return request.NotFound(), nil
	}
	if err != nil {
		return nil, err
	}
	content, err := file.Content.FileData()
	if file.Attachment {
		r.Writer.Header().Set("Content-Disposition", "attachment")
	}
	return request.RawBytes(content), err
}
