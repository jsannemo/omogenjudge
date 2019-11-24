package problems

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/logger"
	"gopkg.in/yaml.v2"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

type problemLimits struct {
	Multiplier int32
	Time       float64
	Memory     int32
}

type problemJudging struct {
	Limits problemLimits
}

type problemMetadata struct {
	Author  string
	License string
	Judging problemJudging
}

func toLicense(l string, reporter util.Reporter) toolspb.License {
	switch l {
	case "permission":
		return toolspb.License_BY_PERMISSION
	case "cc by-sa 3":
		return toolspb.License_CC_BY_SA_3
	case "public domain":
		return toolspb.License_PUBLIC_DOMAIN
	case "private":
		return toolspb.License_PRIVATE
	}
	reporter.Err("Invalid license: %v", l)
	return toolspb.License_LICENSE_UNSPECIFIED
}

func parseMetadata(path string, reporter util.Reporter) (*toolspb.Metadata, error) {
	metadataPath := filepath.Join(path, "problem.yaml")
	logger.Infof("Looking for metadata path %s", metadataPath)
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		reporter.Err("There was no problem.yaml file")
		return nil, nil
	}
	dat, err := ioutil.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}

	var md problemMetadata
	err = yaml.Unmarshal([]byte(dat), &md)
	if err != nil {
		reporter.Err("Invalid problem yaml: %v", err)
		return nil, nil
	}
	limits := md.Judging.Limits
	timeLimit := limits.Time
	timeMultiplier := limits.Multiplier
	memLimit := limits.Memory
	if memLimit == 0 {
		memLimit = 1000
		reporter.Warn("No explicit memory limit set: using default 1000 MB")
	}
	if timeMultiplier == 0 && timeLimit == 0 {
		timeMultiplier = 4
	}
	lic := toLicense(md.License, reporter)
	return &toolspb.Metadata{
		ProblemId: filepath.Base(path),
		Limits: &toolspb.Limits{
			TimeLimitMs:         int32(1000 * timeLimit),
			MemoryLimitKb:       int32(1000 * memLimit),
			TimeLimitMultiplier: timeMultiplier,
		},
		Author:  md.Author,
		License: lic,
	}, nil
}
