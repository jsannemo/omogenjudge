package courses

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

type ContentOpt byte

const (
	// Only load the title of the course
	ContentNone ContentOpt = iota
	// Only load the content of the course and titles of the chapters and sections
	ContentCourse
	// Only load the content of the chapter, and titles of the course and sections
	ContentChapter
	// Only load the content of the section, and titles of the course and chapter
	ContentSection
)

type ListArgs struct {
	Content ContentOpt
}

type ListFilter struct {
	ShortName        string
	ChapterShortName string
	SectionShortName string
}

func listQuery(args ListArgs, filter ListFilter) (string, []interface{}) {
	var params []interface{}
	filterStr := ""
	if filter.ShortName != "" {
		params = append(params, filter.ShortName)
		filterStr = fmt.Sprintf("WHERE course_short_name = $%d", len(params))
	}
	return fmt.Sprintf("SELECT * FROM course %s ORDER BY course_short_name ASC", filterStr), params
}

func List(ctx context.Context, args ListArgs, filter ListFilter) models.CourseList {
	conn := db.Conn()
	query, params := listQuery(args, filter)
	var courses models.CourseList
	if err := conn.SelectContext(ctx, &courses, query, params...); err != nil {
		panic(err)
	}
	includeContent(ctx, courses.AsMap(), args.Content, filter)
	return courses
}

func includeCourseContent(ctx context.Context, courses models.CourseMap, opt ContentOpt) {
	if len(courses) == 0 {
		return
	}
	fields := "course_id, course_language, course_name, course_summary"
	if opt == ContentCourse {
		fields = fields + ", course_desc"
	}

	ids := courses.Ids()
	query, args, err := sqlx.In("SELECT * FROM course_localization WHERE course_id IN (?);", ids)
	if err != nil {
		panic(err)
	}
	conn := db.Conn()
	query = conn.Rebind(query)
	var cs []*models.CourseLoc
	if err := conn.SelectContext(ctx, &cs, query, args...); err != nil {
		panic(err)
	}
	for _, c := range cs {
		course := courses[c.CourseId]
		course.Locs = append(course.Locs, c)
	}
}

func chapterMap(ctx context.Context, c *models.Course, opt ContentOpt, filter ListFilter) models.ChapterMap {
	filterStr := "WHERE course_id = $1"
	params := []interface{}{c.CourseId}
	// TODO implement
	/*
		if filter.ChapterShortName != "" {
			filterStr = filterStr + " AND chapter_short_name = $2"
			params = append(params, filter.ChapterShortName)
		}
	*/

	query := "SELECT * FROM course_chapter " + filterStr + " ORDER BY course_chapter_id ASC"
	if err := db.Conn().SelectContext(ctx, &c.Chapters, query, params...); err != nil {
		panic(err)
	}
	chs := c.Chapters.AsMap()
	for _, ch := range chs {
		ch.Course = c
	}

	fields := "course_chapter_id, chapter_language, chapter_name, chapter_summary"
	if opt == ContentChapter {
		fields = fields + ", chapter_description"
	}

	query = "SELECT " + fields + " FROM course_chapter_localization NATURAL JOIN course_chapter " + filterStr
	var locs []*models.ChapterLoc
	if err := db.Conn().SelectContext(ctx, &locs, query, params...); err != nil {
		panic(err)
	}
	for _, loc := range locs {
		ch := chs[loc.ChapterId]
		ch.Locs = append(ch.Locs, loc)
	}
	return chs
}

func includeSubContent(ctx context.Context, c *models.Course, opt ContentOpt, filter ListFilter) {
	chs := chapterMap(ctx, c, opt, filter)
	if len(chs) == 0 {
		return
	}

	filterStr := "WHERE course_id = $1"
	params := []interface{}{c.CourseId}
	// implement in a good way
	/*
			if filter.ChapterShortName != "" {
				filterStr = filterStr + " AND course_chapter_id = $2"
				params = append(params, chs.First().ChapterId)
		    if filter.SectionShortName != "" {
		      filterStr = filterStr + " AND section_short_name = $3"
		      params = append(params, filter.SectionShortName)
		    }
		  }*/

	var secs models.SectionList
	query := "SELECT course_chapter_id, course_section_id, section_short_name FROM course_section INNER JOIN course_chapter USING (course_chapter_id) " + filterStr + " ORDER BY course_section_id"
	if err := db.Conn().SelectContext(ctx, &secs, query, params...); err != nil {
		panic(err)
	}
	for _, sec := range secs {
		ch := chs[sec.ChapterId]
		ch.Sections = append(ch.Sections, sec)
		sec.Chapter = ch
	}
	smap := secs.AsMap()
	var slocs []*models.SectionLoc
	fields := "course_section_id, section_language, section_name, section_summary"
	if opt == ContentSection {
		fields = fields + ", section_contents"
	}
	query =
		`SELECT ` + fields + ` 
  FROM course_section_localization
  INNER JOIN course_section USING (course_section_id)
  INNER JOIN course_chapter USING (course_chapter_id) ` + filterStr
	if err := db.Conn().SelectContext(ctx, &slocs, query, params...); err != nil {
		panic(err)
	}
	for _, sloc := range slocs {
		sec := smap[sloc.SectionId]
		sec.Locs = append(sec.Locs, sloc)
	}
}

func includeContent(ctx context.Context, courses models.CourseMap, opt ContentOpt, filter ListFilter) {
	includeCourseContent(ctx, courses, opt)
	if opt != ContentNone {
		for _, course := range courses {
			includeSubContent(ctx, course, opt, filter)
		}
	}
}
