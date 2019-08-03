package courses

import (
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/courses"
	"github.com/jsannemo/omogenjudge/storage/models"
)

type Params struct {
	Courses []*models.Course
}

func ListHandler(r *request.Request) (request.Response, error) {
	courses := courses.List(r.Request.Context(), courses.ListArgs{Content: courses.ContentNone}, courses.ListFilter{})
	return request.Template("courses_list", &Params{courses}), nil
}
