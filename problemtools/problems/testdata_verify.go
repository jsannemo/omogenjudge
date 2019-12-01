package problems

import (
	"context"
	"github.com/jsannemo/omogenjudge/util/go/cli"
	"regexp"
	"strings"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

var isTestGroupName = regexp.MustCompile(`^[a-z0-9]+$`).MatchString
var isTestCaseName = regexp.MustCompile(`^[a-z0-9\-_]+$`).MatchString

func verifyTestdata(ctx context.Context, problem *toolspb.Problem, validators []*runpb.CompiledProgram, runner runpb.RunServiceClient, reporter util.Reporter) error {
	for _, g := range problem.TestGroups {
		if len(g.Tests) == 0 {
			reporter.Err("Empty test group %v", g.Name)
		}
		if !isTestGroupName(g.Name) {
			reporter.Err("Invalid test group name: %v [a-z0-9]+", g.Name)
		}
		for _, tc := range g.Tests {
			if !isTestCaseName(tc.Name) {
				reporter.Err("Invalid test case name: %v [a-z0-9\\-_]", tc.Name)
			}
		}
		if err := verifyTestCaseFormats(ctx, g, validators, runner, reporter); err != nil {
			return err
		}
	}
	return nil
}

func verifyTestCaseFormats(ctx context.Context, group *toolspb.TestGroup, validators []*runpb.CompiledProgram, runner runpb.RunServiceClient, reporter util.Reporter) error {
	var inputFiles []string
	var names []string
	for _, tc := range group.Tests {
		inputFiles = append(inputFiles, tc.InputPath)
		names = append(names, tc.FullName)
	}
	for _, validator := range validators {
		resp, err := runner.SimpleRun(ctx, &runpb.SimpleRunRequest{
			Program:    validator,
			InputFiles: inputFiles,
			Arguments:  cli.FormatFlagMap(group.InputFlags),
		})
		if err != nil {
			return err
		}
		for i, res := range resp.Results {
			if res.Timeout {
				reporter.Err("test case %s caused validator to time out")
			} else if res.ExitCode != 42 {
				msg := ""
				if res.Stdout != "" {
					msg = msg + strings.TrimSpace(res.Stdout)
				}
				if res.Stderr != "" {
					msg = msg + strings.TrimSpace(res.Stderr)
				}
				reporter.Err("test case %s failed validation: '%s'", names[i], msg)
			}
		}
	}

	return nil
}
