package contests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jsannemo/omogenjudge/storage/problems"
	"time"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
	"github.com/jsannemo/omogenjudge/storage/contests"
	"github.com/jsannemo/omogenjudge/storage/models"
)

type config struct {
	Title            string            `yaml:"title"`
	ShortName        string            `yaml:"shortname"`
	HostName         string            `yaml:"hostname"`
	ProblemSet       map[string]string `yaml:"problemset"`
	Duration         string            `yaml:"duration"`
	Start            *time.Time        `yaml:"start"`
	HiddenScoreboard bool              `yaml:"hidden_scoreboard"`
}

func InstallContest(ctx context.Context, req *toolspb.InstallContestRequest) (*toolspb.InstallContestResponse, error) {
	config := &config{}
	reporter := util.NewReporter()
	if err := util.ParseYamlString(req.ContestYaml, config); err != nil {
		reporter.Err("Failed parsing contest yaml: %v", err)
	}
	if config.Title == "" {
		reporter.Err("Title is empty")
	}
	if config.ShortName == "" {
		reporter.Err("Short name is empty")
	}
	duration, err := time.ParseDuration(config.Duration)
	if err != nil {
		reporter.Err("Failed parsing duration: %v", err)
	}
	if duration < 0 {
		reporter.Err("Negative duration is not valid: %v", duration)
	}
	contest := &models.Contest{
		Title:            config.Title,
		ShortName:        config.ShortName,
		Duration:         duration,
		HostName:         sql.NullString{String: config.HostName, Valid: config.HostName != ""},
		HiddenScoreboard: config.HiddenScoreboard,
	}
	for label, shortname := range config.ProblemSet {
		problem, err := problems.List(ctx, problems.ListArgs{}, problems.ListFilter{ShortName: shortname})
		if err != nil {
			return nil, fmt.Errorf("failed querying for problem: %v", err)
		}
		if len(problem) != 1 {
			reporter.Err("could not find problem %s", shortname)
		} else {
			contest.Problems = append(contest.Problems, &models.ContestProblem{ProblemID: problem[0].ProblemID, Label: label})
		}
	}
	if config.Start != nil {
		contest.StartTime = sql.NullTime{Time: *config.Start, Valid: true}
	}
	if !reporter.HasError() {
		err := contests.CreateContest(ctx, contest)
		if err != nil {
			if err == contests.ErrShortNameExists {
				if err := contests.UpdateContest(ctx, contest); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}
	return &toolspb.InstallContestResponse{
		Errors:   reporter.Errors(),
		Warnings: reporter.Warnings(),
	}, nil
}
