package problems

import (
	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	problemutils "github.com/jsannemo/omogenjudge/problemtools/util"
)

func verifyIncludedFiles(includedFiles []*toolspb.IncludedFiles, reporter problemutils.Reporter) {
	for _, f := range includedFiles {
		if len(f.FileContents) == 0 {
			reporter.Warn("Language %v had included code directory but no files", f.LanguageId)
		}
	}
}
