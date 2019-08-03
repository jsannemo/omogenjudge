package courses

import (
	"io/ioutil"
	"os"
	"path/filepath"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

type chapterMetadata struct {
	Name     map[string]string
	Sections []string
}

func parseChapter(path string, rep util.Reporter) (*toolspb.Chapter, error) {
	metadata := chapterMetadata{}
	if err := util.ParseYaml(filepath.Join(path, "metadata.yaml"), &metadata); err != nil {
		rep.Err("Failed parsing metadata: %v", err)
		return nil, nil
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var locs []*toolspb.ChapterLoc
	for _, f := range files {
		if f.IsDir() {
			descPath := filepath.Join(path, f.Name(), "desc.md")
			res, err := parseMarkdown(descPath, rep)
			if err != nil {
				if !os.IsNotExist(err) {
					rep.Err("Failed parsing chapter description: %v", err)
					return nil, nil
				}
			} else {
				locs = append(locs, &toolspb.ChapterLoc{
					LanguageCode: f.Name(),
					Name:         res.title,
					Summary:      res.summary,
					Description:  res.output,
				})
			}
		}
	}

	var sections []*toolspb.Section
	for _, secName := range metadata.Sections {
		sec, err := parseSection(path, secName, rep)
		if err != nil {
			return nil, err
		}
		if sec != nil {
			sections = append(sections, sec)
		}
	}
	return &toolspb.Chapter{
		ChapterId:   filepath.Base(path),
		Sections:    sections,
		ChapterLocs: locs,
	}, nil
}
