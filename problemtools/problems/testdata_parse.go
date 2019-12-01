package problems

import (
	"github.com/google/logger"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	util "github.com/jsannemo/omogenjudge/problemtools/util"
)

func parseTestdata(path string, reporter util.Reporter) (map[string]*toolspb.TestGroup, error) {
	testgroups := make(map[string]*toolspb.TestGroup)
	testdataPath := filepath.Join(path, "data")

	if _, err := os.Stat(testdataPath); os.IsNotExist(err) {
		reporter.Err("Problem had no data folder")
		return testgroups, nil
	}

	config, err := parseConfig(testdataPath, reporter)
	if err != nil {
		return nil, err
	}

	logger.Infof("testdata metadata: %v", config)

	files, err := ioutil.ReadDir(testdataPath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			testgroup, err := parseGroup(filepath.Join(testdataPath, f.Name()), configFor(f.Name(), config), testgroups, reporter)
			if err != nil {
				return nil, err
			}
			// The testgroup can be null due to a reporter error that is not a run-time error
			if testgroup != nil {
				testgroups[f.Name()] = testgroup
			}
		}
	}
	return testgroups, nil
}

type testGroupConfig struct {
	Score       int32
	Visibility  string
	InputFlags  map[string]string `yaml:"input_flags"`
	OutputFlags map[string]string `yaml:"output_flags"`
	Include     string
}

func defaultConfig() testGroupConfig {
	return testGroupConfig{
		Score:       0,
		Visibility:  "hidden",
		InputFlags:  map[string]string{},
		OutputFlags: map[string]string{},
		Include:     "",
	}
}

func configFor(group string, configs map[string]testGroupConfig) testGroupConfig {
	config := defaultConfig()
	if def, ok := configs["default"]; ok {
		for k, v := range def.InputFlags {
			config.InputFlags[k] = v
		}
		for k, v := range def.OutputFlags {
			config.OutputFlags[k] = v
		}
	}
	if group == "sample" || group == "samples" {
		config.Visibility = "public"
	}
	if v, ok := configs[group]; ok {
		config.Score = v.Score
		config.Include = v.Include
		if v.Visibility != "" {
			config.Visibility = v.Visibility
		}
		for k, v := range v.InputFlags {
			config.InputFlags[k] = v
		}
		for k, v := range v.OutputFlags {
			config.OutputFlags[k] = v
		}
	}
	return config
}

func parseConfig(path string, reporter util.Reporter) (map[string]testGroupConfig, error) {
	configPath := filepath.Join(path, "testdata.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil
	}

	dat, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config map[string]testGroupConfig
	err = yaml.Unmarshal([]byte(dat), &config)
	if err != nil {
		reporter.Err("Invalid config yaml: %v", err)
		return nil, nil
	}

	for _, v := range config {
		if v.Visibility != "" && v.Visibility != "public" && v.Visibility != "hidden" {
			reporter.Err("Visibility value %s is invalid (expected public or hidden)", v.Visibility)
		}
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
				FullName:   filepath.Base(path) + "/" + tcName,
				InputPath:  baseName + ".in",
				OutputPath: baseName + ".ans",
			})
		}
	}
	return cases, nil
}

func parseGroup(path string, config testGroupConfig, groups map[string]*toolspb.TestGroup, reporter util.Reporter) (*toolspb.TestGroup, error) {
	tests, err := parseTests(path, reporter)
	if err != nil {
		return nil, err
	}
	// Ignore groups without test cases
	if len(tests) == 0 {
		return nil, nil
	}
	inc := strings.Split(strings.TrimSpace(config.Include), " ")
	for _, v := range inc {
		if v == "" {
			continue
		}
		tg, ok := groups[v]
		if !ok {
			reporter.Err("Could not include group %s, was not defined yet", v)
		} else {
			tests = append(tests, tg.Tests...)
		}
	}
	return &toolspb.TestGroup{
		Name:          filepath.Base(path),
		PublicSamples: config.Visibility == "public",
		Score:         config.Score,
		InputFlags:    config.InputFlags,
		OutputFlags:   config.OutputFlags,
		Tests:         tests,
	}, nil
}
