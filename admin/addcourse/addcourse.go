package main

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/logger"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	ptclient "github.com/jsannemo/omogenjudge/problemtools/client"
	"github.com/jsannemo/omogenjudge/storage/courses"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/util/go/cli"
	futil "github.com/jsannemo/omogenjudge/util/go/files"
)

func toStorageSectionLocs(sec *toolspb.Section) []*models.SectionLoc {
	var locs []*models.SectionLoc
	for _, loc := range sec.SectionLocs {
		locs = append(locs, &models.SectionLoc{Language: loc.LanguageCode, Name: models.HTMLString(loc.Title), Summary: models.HTMLString(loc.Summary), Contents: loc.Contents})
	}
	return locs
}

func toStorageSection(sec *toolspb.Section) *models.Section {
	return &models.Section{
		ShortName: sec.SectionId,
		Locs:      toStorageSectionLocs(sec),
	}
}

func toStorageSections(toolSecs []*toolspb.Section) []*models.Section {
	var secs []*models.Section
	for _, ch := range toolSecs {
		secs = append(secs, toStorageSection(ch))
	}
	return secs
}

func toStorageChapterLocs(ch *toolspb.Chapter) []*models.ChapterLoc {
	var locs []*models.ChapterLoc
	for _, loc := range ch.ChapterLocs {
		locs = append(locs, &models.ChapterLoc{Language: loc.LanguageCode, Name: models.HTMLString(loc.Name), Summary: models.HTMLString(loc.Summary), Desc: models.HTMLString(loc.Description)})
	}
	return locs
}

func toStorageChapter(ch *toolspb.Chapter) *models.Chapter {
	return &models.Chapter{
		ShortName: ch.ChapterId,
		Locs:      toStorageChapterLocs(ch),
		Sections:  toStorageSections(ch.Sections),
	}
}

func toStorageChapters(toolChs []*toolspb.Chapter) []*models.Chapter {
	var chs []*models.Chapter
	for _, ch := range toolChs {
		chs = append(chs, toStorageChapter(ch))
	}
	return chs
}

func toStorageCourseLocs(course *toolspb.Course) []*models.CourseLoc {
	var locs []*models.CourseLoc
	for _, loc := range course.CourseLocs {
		locs = append(locs, &models.CourseLoc{Language: loc.LanguageCode, Name: models.HTMLString(loc.Name), Summary: models.HTMLString(loc.Summary), Desc: models.HTMLString(loc.Description)})
	}
	return locs
}

func toStorageCourse(course *toolspb.Course) *models.Course {
	return &models.Course{
		ShortName: course.CourseId,
		Locs:      toStorageCourseLocs(course),
		Chapters:  toStorageChapters(course.Chapters),
	}
}

func installCourse(path string) error {
	logger.Infof("Installing course %s", path)

	tmp, err := ioutil.TempDir("/tmp", "omogeninstall")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	if err := os.Chmod(tmp, 0755); err != nil {
		return err
	}
	npath := filepath.Join(tmp, filepath.Base(path))
	if err := futil.CopyDirectory(path, npath); err != nil {
		return err
	}

	ctx := context.Background()

	client := ptclient.NewClient()
	parsed, err := client.ParseCourse(ctx, &toolspb.ParseCourseRequest{
		CoursePath: npath,
	})
	if err != nil {
		logger.Fatalf("Failed parsing problem: %v", err)
	}
	for _, warnMsg := range parsed.Warnings {
		logger.Warningln(warnMsg)
	}
	for _, errMsg := range parsed.Errors {
		logger.Errorln(errMsg)
	}
	if len(parsed.Errors) != 0 {
		logger.Errorf("Problem had errors; will not install")
		return nil
	}
	if len(parsed.Warnings) != 0 {
		if !cli.RequestConfirmation("The problem had warnings; do you still want to install it?") {
			return nil
		}
	}
	logger.Infof("parsed %v", parsed)

	course := toStorageCourse(parsed.ParsedCourse)
	return courses.Create(ctx, course)
}

func main() {
	flag.Parse()
	path := flag.Arg(0)
	path, err := filepath.Abs(path)
	if err != nil {
		logger.Fatal(err)
	}

	if err := installCourse(path); err != nil {
		logger.Fatalf("failed installing course: %v", err)
	}
}
