package editor

import (
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/frontend/util"
)

type ViewParams struct {
	Languages []*util.Language
}

func ViewHandler(r *request.Request) (request.Response, error) {
	return request.Template("editor_view", &ViewParams{util.Languages()}), nil
}
