// Parsing of problem testdata.
package testdata

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	util "github.com/jsannemo/omogenjudge/problemtools/util"
)

// ParseTestdata parses the testdata of the problem at the given path.
func ParseTestdata(path string, reporter util.Reporter) ([]*toolspb.TestGroup, error) {
	testgroups := make([]*toolspb.TestGroup, 0)
	testdataPath := filepath.Join(path, "testdata")

	if _, err := os.Stat(testdataPath); os.IsNotExist(err) {
		reporter.Err("Problem had no testdata folder")
		return testgroups, nil
	}

	files, err := ioutil.ReadDir(testdataPath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			testgroup, err := parseGroup(filepath.Join(testdataPath, f.Name()), reporter)
			if err != nil {
				return nil, err
			}
			// The testgroup can be null due to a reporter error that is not a run-time error
			if testgroup != nil {
				testgroups = append(testgroups, testgroup)
			}
		}
	}
	return testgroups, nil
}

type groupConfig struct {
	Visibility string
}

func defaultConfig() groupConfig {
	return groupConfig{
		Visibility: "hidden",
	}
}

func parseConfig(path string, reporter util.Reporter) (groupConfig, error) {
	config := defaultConfig()
	configPath := filepath.Join(path, "testgroup.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil
	}

	dat, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal([]byte(dat), &config)
	if err != nil {
		reporter.Err("Invalid config yaml: %v", err)
		return config, nil
	}

	if config.Visibility != "public" && config.Visibility != "hidden" {
		reporter.Err("Visibility value %s is invalid (expected public or hidden)", config.Visibility)
	}

	return config, nil
}

func parseTests(path string, reporter util.Reporter) ([]*toolspb.TestCase, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	hasInput := make(map[string]bool)
	hasOutput := make(map[string]bool)

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		filePath := filepath.Join(path, f.Name())
		if strings.HasSuffix(filePath, ".in") {
			basePath := strings.TrimSuffix(filePath, ".in")
			hasInput[basePath] = true
		} else if strings.HasSuffix(filePath, ".ans") {
			basePath := strings.TrimSuffix(filePath, ".ans")
			hasOutput[basePath] = true
		}
	}
	for baseName, _ := range hasInput {
		if !hasOutput[baseName] {
			reporter.Err("Test case %s has no matching output", baseName)
		}
	}

	for baseName, _ := range hasOutput {
		if !hasInput[baseName] {
			reporter.Err("Test case %s has no matching input", baseName)
		}
	}

	var cases []*toolspb.TestCase
	for baseName, _ := range hasInput {
		if hasOutput[baseName] {
			tcName := filepath.Base(baseName)
			cases = append(cases, &toolspb.TestCase{
				Name:       tcName,
				InputPath:  baseName + ".in",
				OutputPath: baseName + ".ans",
			})
		}
	}
	return cases, nil
}

func parseGroup(path string, reporter util.Reporter) (*toolspb.TestGroup, error) {
	config, err := parseConfig(path, reporter)
	if err != nil {
		return nil, err
	}
	tests, err := parseTests(path, reporter)
	if err != nil {
		return nil, err
	}
	return &toolspb.TestGroup{
		Name:       filepath.Base(path),
		PublicSamples: config.Visibility == "public",
		Tests:         tests,
	}, nil
}
