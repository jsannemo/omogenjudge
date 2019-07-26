// Problem metadata parsing.
package metadata

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

// The following structs mirror the Yaml configuration structure of the metadata.yaml file
type limits struct {
	Time   float64
	Memory int32
}

type judging struct {
	Limits limits
}

type metadata struct {
	Author  string
	License string
	Judging judging
}

// ParseMetadata parses the metadata of a problem, reporting potential parsing errors.
func ParseMetadata(path string, reporter util.Reporter) (*toolspb.Metadata, error) {
	metadataPath := filepath.Join(path, "metadata.yaml")
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		reporter.Err("There was no metadata.yaml file")
		return nil, nil
	}

	dat, err := ioutil.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}

	var md metadata
	err = yaml.Unmarshal([]byte(dat), &md)
	if err != nil {
		reporter.Err("Invalid metadata yaml: %v", err)
		return nil, nil
	}
	timeLimit := md.Judging.Limits.Time
	if 0 > timeLimit || timeLimit > 120 {
		reporter.Err("Time limit out of bounds: %v", timeLimit)
	}
	memLimit := md.Judging.Limits.Memory
	if memLimit == 0 {
		memLimit = 1024
		reporter.Warn("No explicit memory limit set: using default 1024 MB")
	}
	if 0 > memLimit || memLimit > 5 * 1024 {
		reporter.Err("Memory limit out of bounds: %v", memLimit)
	}
	return &toolspb.Metadata{
		ProblemId: filepath.Base(path),
		Limits: &toolspb.Limits{
			TimeLimitMilliseconds: int32(1000 * timeLimit),
			MemoryLimitMegabytes:  memLimit,
		},
	}, nil
}
