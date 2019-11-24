package submissions

import (
	"context"
	"fmt"
	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
	"strings"
)

// A ListArgs controls what data to include in submissions.
type ListArgs struct {
	// Whether to include submission files in the query.
	WithFiles bool
	// Whether to include the current run in the query.
	WithRun bool
}

// A ListFilter controls what submissions to search for.  Only one of SubmissionID and UserID may be set.
type ListFilter struct {
	SubmissionID []int32
	UserID       []int32
	ProblemID    []int32
}

// A SubmissionList is a slice of Submissions.
type SubmissionList []*models.Submission

// ListSubmissions searches for a list of submissions.
func ListSubmissions(ctx context.Context, args ListArgs, filterArgs ListFilter) (SubmissionList, error) {
	conn := db.Conn()
	var params []interface{}
	var filters []string
	joins := ""
	fields := ""
	if len(filterArgs.SubmissionID) != 0 {
		filters = append(filters, db.SetInParamInt(`submission.submission_id IN (%s)`, &params, filterArgs.SubmissionID))
	}
	if len(filterArgs.UserID) != 0 {
		filters = append(filters, db.SetInParamInt(`account_id IN (%s)`, &params, filterArgs.UserID))
	}
	if len(filterArgs.ProblemID) != 0 {
		filters = append(filters, db.SetInParamInt(`problem_id IN (%s)`, &params, filterArgs.ProblemID))
	}
	filter := ""
	if len(filters) > 0 {
		filter = fmt.Sprintf(" WHERE %s ", strings.Join(filters, " AND "))
	}
	if args.WithRun {
		joins = "LEFT JOIN submission_run ON current_run = submission_run_id"
		fields = `
			,
			submission_run.submission_run_id "submission_run.submission_run_id",
			submission_run.submission_id "submission_run.submission_id",
			submission_run.problem_version_id "submission_run.problem_version_id",
			submission_run.date_created "submission_run.date_created",
			submission_run.time_usage_ms "submission_run.time_usage_ms",
			submission_run.status "submission_run.status",
			submission_run.score "submission_run.score",
			submission_run.verdict "submission_run.verdict",
			submission_run.compile_error "submission_run.compile_error"
`
	}
	query := fmt.Sprintf(`
		SELECT
			submission.submission_id, account_id, problem_id, language, submission.date_created
			%s
		FROM submission
		%s
		%s
		ORDER BY submission.submission_id DESC`,
		fields,
		joins,
		filter)
	var subs SubmissionList
	if err := conn.SelectContext(ctx, &subs, query, params...); err != nil {
		return nil, fmt.Errorf("could not read submissions: %v", err)
	}
	if args.WithFiles {
		for _, sub := range subs {
			if err := conn.SelectContext(ctx, &sub.Files, `SELECT * FROM submission_file WHERE submission_id = $1`, sub.SubmissionID); err != nil {
				return nil, fmt.Errorf("could not read submission files: %v", err)
			}
		}
	}
	return subs, nil
}
