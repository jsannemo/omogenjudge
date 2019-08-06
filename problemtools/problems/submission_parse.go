package problems

import (
	"fmt"
	"github.com/google/logger"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

type submissionMetadata struct {
	RequiredFailures  []string `yaml:"required_failures"`
	AllowedFailures   []string `yaml:"allowed_failures"`
	ExcludeFromTiming bool     `yaml:"exclude_from_timing"`
}

func parseSubmissions(path string, reporter util.Reporter) ([]*toolspb.Submission, error) {
	subPath := filepath.Join(path, "submissions")
	metadata := make(map[string]submissionMetadata)
	metadataPath := filepath.Join(subPath, "submissions.yaml")
	if _, err := os.Stat(metadataPath); !os.IsNotExist(err) {
		dat, err := ioutil.ReadFile(metadataPath)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(dat, &metadata)
		logger.Infof("metadata %s %s", metadata, dat)
	}

	if _, err := os.Stat(subPath); os.IsNotExist(err) {
		return nil, nil
	}
	files, err := ioutil.ReadDir(subPath)
	if err != nil {
		return nil, fmt.Errorf("failed opening submission path: %v", err)
	}
	var submissions []*toolspb.Submission
	for _, f := range files {
		if f.Name() == "submissions.yaml" {
			continue
		}
		program, err := parseProgram(filepath.Join(subPath, f.Name()), f.IsDir(), reporter)
		if err != nil {
			return nil, fmt.Errorf("failed parsing submission: %v", err)
		}
		if program != nil {
			logger.Infof("found submission %v", program)
			subMetadata, _ := metadata[f.Name()]
			subMetadata.AllowedFailures = append(subMetadata.AllowedFailures, subMetadata.RequiredFailures...)
			subMetadata.AllowedFailures = append(subMetadata.AllowedFailures, "ac")
			for _, failure := range subMetadata.AllowedFailures {
				if failure == "tle" {
					subMetadata.ExcludeFromTiming = true
				}
			}
			submissions = append(submissions, &toolspb.Submission{
				Submission:   program,
				UseForTiming: !subMetadata.ExcludeFromTiming,
				Constraint:   toConstraint(subMetadata, reporter),
				Name:         f.Name(),
			})
		}
	}
	return submissions, nil
}

func toConstraint(metadata submissionMetadata, rep util.Reporter) *toolspb.SubmissionConstraint {
	return &toolspb.SubmissionConstraint{
		RequiredFailures: toVerdicts(metadata.RequiredFailures, rep),
		AllowedFailures:  toVerdicts(metadata.AllowedFailures, rep),
	}
}

func toVerdicts(strings []string, rep util.Reporter) []runpb.Verdict {
	var res []runpb.Verdict
	for _, verdict := range strings {
		res = append(res, toVerdict(verdict, rep))
	}
	return res
}

func toVerdict(s string, rep util.Reporter) runpb.Verdict {
	switch s {
	case "tle":
		return runpb.Verdict_TIME_LIMIT_EXCEEDED
	case "wa":
		return runpb.Verdict_WRONG_ANSWER
	case "rte":
		return runpb.Verdict_RUN_TIME_ERROR
	case "ac":
		return runpb.Verdict_ACCEPTED
	}
	rep.Err("Invalid verdict in submissions.yaml: %s", s)
	return runpb.Verdict_VERDICT_UNSPECIFIED
}
