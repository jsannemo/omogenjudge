package courses

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

func insertSectionLoc(ctx context.Context, sl *models.SectionLoc, tx *sqlx.Tx) error {
	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO
      course_section_localization(course_section_id, section_language, section_contents, section_summary, section_name)
      VALUES($1, $2, $3, $4, $5)`,
		sl.SectionId, sl.Language, sl.Contents, sl.Summary, sl.Name); err != nil {
		return err
	}
	return nil
}

func insertSection(ctx context.Context, s *models.Section, tx *sqlx.Tx) error {
	if err := tx.QueryRowContext(
		ctx,
		"INSERT INTO course_section(course_chapter_id, section_short_name) VALUES($1, $2) RETURNING course_section_id",
		s.ChapterId, s.ShortName).Scan(&s.SectionId); err != nil {
		return err
	}

	for _, loc := range s.Locs {
		loc.SectionId = s.SectionId
		if err := insertSectionLoc(ctx, loc, tx); err != nil {
			return err
		}
	}
	return nil
}

func insertChapterLoc(ctx context.Context, cl *models.ChapterLoc, tx *sqlx.Tx) error {
	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO
      course_chapter_localization(course_chapter_id, chapter_language, chapter_summary, chapter_description, chapter_name)
      VALUES($1, $2, $3, $4, $5)`,
		cl.ChapterId, cl.Language, cl.Summary, cl.Desc, cl.Name); err != nil {
		return err
	}
	return nil
}

func insertChapter(ctx context.Context, ch *models.Chapter, tx *sqlx.Tx) error {
	if err := tx.QueryRowContext(
		ctx,
		"INSERT INTO course_chapter(course_id, chapter_short_name) VALUES($1, $2) RETURNING course_chapter_id",
		ch.CourseId, ch.ShortName).Scan(&ch.ChapterId); err != nil {
		return err
	}

	for _, loc := range ch.Locs {
		loc.ChapterId = ch.ChapterId
		if err := insertChapterLoc(ctx, loc, tx); err != nil {
			return err
		}
	}
	for _, sec := range ch.Sections {
		sec.ChapterId = ch.ChapterId
		if err := insertSection(ctx, sec, tx); err != nil {
			return err
		}
	}
	return nil
}

func insertCourseLoc(ctx context.Context, cl *models.CourseLoc, tx *sqlx.Tx) error {
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO
      course_localization(course_id, course_language, course_name, course_summary, course_description)
      VALUES($1, $2, $3, $4, $5)`,
		cl.CourseId, cl.Language, cl.Name, cl.Summary, cl.Desc); err != nil {
		return err
	}
	return nil
}

func insertCourse(ctx context.Context, c *models.Course, tx *sqlx.Tx) error {
	if err := tx.QueryRowContext(ctx, "INSERT INTO course(course_short_name) VALUES($1) RETURNING course_id", c.ShortName).Scan(&c.CourseId); err != nil {
		return err
	}

	for _, loc := range c.Locs {
		loc.CourseId = c.CourseId
		if err := insertCourseLoc(ctx, loc, tx); err != nil {
			return err
		}
	}

	for _, ch := range c.Chapters {
		ch.CourseId = c.CourseId
		if err := insertChapter(ctx, ch, tx); err != nil {
			return err
		}
	}

	return nil
}

func Create(ctx context.Context, c *models.Course) error {
	err := db.InTransaction(ctx, func(tx *sqlx.Tx) error {
		return insertCourse(ctx, c, tx)
	})
	return err
}
