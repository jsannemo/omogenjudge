package courses

import (
	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/courses"
	"github.com/jsannemo/omogenjudge/storage/models"

	"github.com/gorilla/mux"
)

type ChapterParams struct {
	Chapter *models.Chapter
}

func ChapterHandler(r *request.Request) (request.Response, error) {
	vars := mux.Vars(r.Request)
	shortName := vars[paths.CourseChapterNameArg]
	courses := courses.List(r.Request.Context(), courses.ListArgs{Content: courses.ContentChapter}, courses.ListFilter{
		ShortName:        vars[paths.CourseNameArg],
		ChapterShortName: vars[paths.CourseChapterNameArg],
	})
	if len(courses) == 0 {
		return request.NotFound(), nil
	}
	course := courses[0]
	if len(course.Chapters) == 0 {
		return request.NotFound(), nil
	}
	chapter, err := course.Chapters.ShortName(shortName)
	if err != nil {
		return request.NotFound(), nil
	}
	return request.Template("courses_chapter", &ChapterParams{chapter}), nil
}
