package problems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jsannemo/omogenjudge/problemtools/util"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

func parseInputValidators(path string, reporter util.Reporter) ([]*runpb.Program, error) {
	validatorPath := filepath.Join(path, "input_validators")
	if _, err := os.Stat(validatorPath); os.IsNotExist(err) {
		return nil, nil
	}
	files, err := ioutil.ReadDir(validatorPath)
	if err != nil {
		return nil, fmt.Errorf("failed opening input validator path: %v", err)
	}
	var programs []*runpb.Program
	for _, f := range files {
		program, err := parseProgram(filepath.Join(validatorPath, f.Name()), f.IsDir(), reporter)
		if err != nil {
			return nil, fmt.Errorf("failed parsing input validator: %v", err)
		}
		if program != nil {
			programs = append(programs, program)
		}
	}
	return programs, nil
}
