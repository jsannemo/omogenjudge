package courses

import (
	"io/ioutil"
	"os"
	"path/filepath"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

type courseMetadata struct {
	Chapters []string
}

func parseCourse(path string, rep util.Reporter) (*toolspb.Course, error) {
	metadata := courseMetadata{}
	if err := util.ParseYaml(filepath.Join(path, "metadata.yaml"), &metadata); err != nil {
		rep.Err("Failed parsing metadata: %v", err)
		return nil, nil
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var locs []*toolspb.CourseLoc
	for _, f := range files {
		if f.IsDir() {
			descPath := filepath.Join(path, f.Name(), "desc.md")
			res, err := parseMarkdown(descPath, rep)
			if err != nil {
				if !os.IsNotExist(err) {
					rep.Err("Failed parsing description: %v", err)
					return nil, nil
				}
			} else {
				locs = append(locs, &toolspb.CourseLoc{
					LanguageCode: f.Name(),
					Name:         res.title,
					Summary:      res.summary,
					Description:  res.output,
				})
			}
		}
	}

	var chapters []*toolspb.Chapter
	for _, chName := range metadata.Chapters {
		ch, err := parseChapter(filepath.Join(path, chName), rep)
		if err != nil {
			return nil, err
		}
		if ch != nil {
			chapters = append(chapters, ch)
		}
	}
	return &toolspb.Course{
		CourseId:   filepath.Base(path),
		Chapters:   chapters,
		CourseLocs: locs,
	}, nil
}
