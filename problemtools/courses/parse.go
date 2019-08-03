package courses

import (
	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

// ParseProblem parses the problem at the given path into the API format.
func ParseCourse(path string) (*toolspb.ParseCourseResponse, error) {
	var errors []string
	var warnings []string

	courseReporter := util.NewReporter()
	course, err := parseCourse(path, courseReporter)
	if err != nil {
		return nil, err
	}
	courseReporter.AddFailures(&errors, &warnings)

	return &toolspb.ParseCourseResponse{
		ParsedCourse: course,
		Errors:       errors,
		Warnings:     warnings,
	}, nil
}
