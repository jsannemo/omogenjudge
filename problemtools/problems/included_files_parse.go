package problems

import (
	"io/ioutil"
	"os"
	"path/filepath"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	problemutils "github.com/jsannemo/omogenjudge/problemtools/util"
)

func parseIncludedFiles(path string, _ problemutils.Reporter) ([]*toolspb.IncludedFiles, error) {
	inclusionPath := filepath.Join(path, "include")
	if _, err := os.Stat(inclusionPath); os.IsNotExist(err) {
		return nil, nil
	}
	inclusionFiles, err := ioutil.ReadDir(inclusionPath)
	if err != nil {
		return nil, err
	}
	var allIncluded []*toolspb.IncludedFiles
	for _, f := range inclusionFiles {
		if f.IsDir() {
			includedFiles := &toolspb.IncludedFiles{
				LanguageId:   f.Name(),
				FileContents: make(map[string]string),
			}
			langPath := filepath.Join(inclusionPath, f.Name())
			langFiles, err := ioutil.ReadDir(langPath)
			if err != nil {
				return nil, err
			}
			for _, f := range langFiles {
				if !f.IsDir() {
					content, err := ioutil.ReadFile(filepath.Join(langPath, f.Name()))
					if err != nil {
						return nil, err
					}
					includedFiles.FileContents[f.Name()] = string(content)
				}
			}
			allIncluded = append(allIncluded, includedFiles)
		}
	}
	return allIncluded, nil
}
