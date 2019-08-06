package problems

import (
	"os"
	"path/filepath"

	util "github.com/jsannemo/omogenjudge/problemtools/util"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

func parseOutputValidator(path string, reporter util.Reporter) (*runpb.Program, error) {
	validatorPath := filepath.Join(path, "output_validator")
	if _, err := os.Stat(validatorPath); os.IsNotExist(err) {
		return nil, nil
	}
	return parseProgram(validatorPath, true, reporter)
}
