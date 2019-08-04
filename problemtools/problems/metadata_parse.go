package problems

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

type problemLimits struct {
	Time   float64
	Memory int32
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
	}
	reporter.Err("Invalid license: %v", l)
	return toolspb.License_LICENSE_UNSPECIFIED
}

func parseMetadata(path string, reporter util.Reporter) (*toolspb.Metadata, error) {
	metadataPath := filepath.Join(path, "metadata.yaml")
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		reporter.Err("There was no metadata.yaml file")
		return nil, nil
	}

	dat, err := ioutil.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}

	var md problemMetadata
	err = yaml.Unmarshal([]byte(dat), &md)
	if err != nil {
		reporter.Err("Invalid metadata yaml: %v", err)
		return nil, nil
	}
	timeLimit := md.Judging.Limits.Time
	memLimit := md.Judging.Limits.Memory
	if memLimit == 0 {
		memLimit = 1000
		reporter.Warn("No explicit memory limit set: using default 1000 MB")
	}
	lic := toLicense(md.License, reporter)
	return &toolspb.Metadata{
		ProblemId: filepath.Base(path),
		Limits: &toolspb.Limits{
			TimeLimitMs:   int32(1000 * timeLimit),
			MemoryLimitKb: int32(1000 * memLimit),
		},
		Author:  md.Author,
		License: lic,
	}, nil
}
