package models

import (
	"golang.org/x/text/language"

	"github.com/jsannemo/omogenjudge/frontend/paths"
)

type Course struct {
	CourseId  int32  `db:"course_id"`
	ShortName string `db:"course_short_name"`

	Locs     []*CourseLoc
	Chapters ChapterList
}

func (c *Course) Loc(preferred []language.Tag) *CourseLoc {
	var has []language.Tag
	userPrefs := append(preferred, language.Make("en"), language.Make("sv"))
	for _, loc := range c.Locs {
		has = append(has, language.Make(loc.Language))
	}
	var matcher = language.NewMatcher(has)
	_, index, _ := matcher.Match(userPrefs...)
	return c.Locs[index]
}

func (c *Course) Link() string {
	return paths.Route(paths.Course, paths.CourseNameArg, c.ShortName)
}

type CourseList []*Course

func (cl CourseList) AsMap() CourseMap {
	cm := make(CourseMap)
	for _, c := range cl {
		cm[c.CourseId] = c
	}
	return cm
}

type CourseMap map[int32]*Course

func (c CourseMap) Ids() []int32 {
	var ids []int32
	for id, _ := range c {
		ids = append(ids, id)
	}
	return ids
}

type CourseLoc struct {
	CourseId int32      `db:"course_id"`
	Language string     `db:"course_language"`
	Name     HTMLString `db:"course_name"`
	Summary  HTMLString `db:"course_summary"`
	Desc     HTMLString `db:"course_description"`
}

type Chapter struct {
	Course    *Course
	CourseId  int32  `db:"course_id"`
	ChapterId int32  `db:"course_chapter_id"`
	ShortName string `db:"chapter_short_name"`

	Locs     []*ChapterLoc
	Sections SectionList
}

func (c *Chapter) Loc(preferred []language.Tag) *ChapterLoc {
	var has []language.Tag
	userPrefs := append(preferred, language.Make("en"), language.Make("sv"))
	for _, loc := range c.Locs {
		has = append(has, language.Make(loc.Language))
	}
	var matcher = language.NewMatcher(has)
	_, index, _ := matcher.Match(userPrefs...)
	return c.Locs[index]
}

func (c *Chapter) Link() string {
	return paths.Route(paths.CourseChapter,
		paths.CourseNameArg, c.Course.ShortName,
		paths.CourseChapterNameArg, c.ShortName)
}

func (c *Chapter) NextChapter() *Chapter {
	for i, ch := range c.Course.Chapters {
		if c.ChapterId == ch.ChapterId && i+1 != len(c.Course.Chapters) {
			return c.Course.Chapters[i+1]
		}
	}
	return nil
}

type ChapterList []*Chapter

func (cl ChapterList) AsMap() ChapterMap {
	cm := make(ChapterMap)
	for _, c := range cl {
		cm[c.ChapterId] = c
	}
	return cm
}

type ChapterMap map[int32]*Chapter

func (cm ChapterMap) First() *Chapter {
	for _, ch := range cm {
		return ch
	}
	return nil
}

type ChapterLoc struct {
	ChapterId int32      `db:"course_chapter_id"`
	Language  string     `db:"chapter_language"`
	Name      HTMLString `db:"chapter_name"`
	Summary   HTMLString `db:"chapter_summary"`
	Desc      HTMLString `db:"chapter_description"`
}

type Section struct {
	Chapter   *Chapter
	ChapterId int32  `db:"course_chapter_id"`
	SectionId int32  `db:"course_section_id"`
	ShortName string `db:"section_short_name"`

	Locs []*SectionLoc
}

func (s *Section) Loc(preferred []language.Tag) *SectionLoc {
	var has []language.Tag
	userPrefs := append(preferred, language.Make("en"), language.Make("sv"))
	for _, loc := range s.Locs {
		has = append(has, language.Make(loc.Language))
	}
	var matcher = language.NewMatcher(has)
	_, index, _ := matcher.Match(userPrefs...)
	return s.Locs[index]
}

func (s *Section) Link() string {
	return paths.Route(paths.CourseSection,
		paths.CourseNameArg, s.Chapter.Course.ShortName,
		paths.CourseChapterNameArg, s.Chapter.ShortName,
		paths.CourseSectionNameArg, s.ShortName)
}

func (s *Section) NextSection() *Section {
	for i, sec := range s.Chapter.Sections {
		if s.SectionId == sec.SectionId && i+1 != len(s.Chapter.Sections) {
			return s.Chapter.Sections[i+1]
		}
	}
	ch := s.Chapter.NextChapter()
	if ch != nil {
		return ch.Sections[0]
	}
	return nil
}

type SectionList []*Section

func (cl SectionList) AsMap() SectionMap {
	cm := make(SectionMap)
	for _, c := range cl {
		cm[c.SectionId] = c
	}
	return cm
}

type SectionMap map[int32]*Section

type SectionLoc struct {
	SectionId int32      `db:"course_section_id"`
	Language  string     `db:"section_language"`
	Name      HTMLString `db:"section_name"`
	Summary   HTMLString `db:"section_summary"`
	Contents  string     `db:"section_contents"`
}
