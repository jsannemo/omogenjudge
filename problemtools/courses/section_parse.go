package courses

import (
	"io/ioutil"
	"os"
	"path/filepath"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

func parseSection(path string, section string, rep util.Reporter) (*toolspb.Section, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var locs []*toolspb.SectionLoc
	for _, f := range files {
		if f.IsDir() {
			contPath := filepath.Join(path, f.Name(), section+".md")
			res, err := parseMarkdown(contPath, rep)
			if err != nil {
				if !os.IsNotExist(err) {
					rep.Err("Failed parsing section content: %v", err)
					return nil, err
				}
			} else {
				locs = append(locs, &toolspb.SectionLoc{
					LanguageCode: f.Name(),
					Title:        res.title,
					Summary:      res.summary,
					Contents:     res.output,
				})
			}
		}
	}

	return &toolspb.Section{
		SectionId:   filepath.Base(section),
		SectionLocs: locs,
	}, nil
}
