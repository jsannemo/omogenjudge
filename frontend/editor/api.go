package editor

import (
  "net/http"

	"github.com/jsannemo/omogenjudge/frontend/request"
)

func ApiFiles(r *request.Request) (request.Response, error) {
  if r.Request.Method == http.MethodPost {
    return request.Raw("test"), nil
  }
  return request.Raw("test"), nil
}
