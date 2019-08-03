package courses

import (
	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/courses"
	"github.com/jsannemo/omogenjudge/storage/models"

	"github.com/gorilla/mux"
)

type CourseParams struct {
	Course *models.Course
}

func CourseHandler(r *request.Request) (request.Response, error) {
	vars := mux.Vars(r.Request)
	courses := courses.List(r.Request.Context(), courses.ListArgs{Content: courses.ContentCourse}, courses.ListFilter{ShortName: vars[paths.CourseNameArg]})
	if len(courses) == 0 {
		return request.NotFound(), nil
	}
	course := courses[0]
	return request.Template("courses_course", &CourseParams{course}), nil
}
