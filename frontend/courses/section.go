package courses

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/gorilla/mux"

	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/storage/courses"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/problems"
)

type SectionParams struct {
	Section *models.Section
	Output  template.HTML
}

func SectionHandler(r *request.Request) (request.Response, error) {
	vars := mux.Vars(r.Request)
	courses := courses.List(r.Request.Context(), courses.ListArgs{Content: courses.ContentSection}, courses.ListFilter{
		ShortName:        vars[paths.CourseNameArg],
		ChapterShortName: vars[paths.CourseChapterNameArg],
		SectionShortName: vars[paths.CourseSectionNameArg],
	})
	if len(courses) == 0 {
		return request.NotFound(), nil
	}
	course := courses[0]
	if len(course.Chapters) == 0 {
		return request.NotFound(), nil
	}
	chapter := course.Chapters[0]
	if len(chapter.Sections) == 0 {
		return request.NotFound(), nil
	}
	section := chapter.Sections[0]

	tpl := template.New("").Funcs(map[string]interface{}{
		"loadProblem": func(shortName string) *models.Problem {
			return problems.List(r.Request.Context(), problems.ListArgs{WithStatements: problems.StmtAll, WithTests: problems.TestsSamples}, problems.ListFilter{ShortName: shortName})[0]
		},
		"ctx": func() *request.RequestContext {
			return r.Context
		},
	})
	tpl, err := tpl.ParseFiles("frontend/templates/courses/content-helpers/helpers.tpl")
	if err != nil {
		return nil, err
	}
	tpl, err = tpl.Parse(section.Loc(r.Context.Locales).Contents)
	if err != nil {
		return nil, err
	}
	var rendered bytes.Buffer
	if err := tpl.Execute(&rendered, nil); err != nil {
		return nil, err
	}
	fmt.Println(rendered.String())
	return request.Template("courses_section", &SectionParams{section, template.HTML(rendered.String())}), nil
}
