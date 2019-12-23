package submissions

import (
	"context"
	"fmt"
	"strings"

	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

// A ListArgs controls what data to include in submissions.
type ListArgs struct {
	// Whether to include submission files in the query.
	WithFiles bool
	// Whether to include the current run in the query.
	WithRun       bool
	WithGroupRuns bool
	WithAccounts  bool
}

// A ListFilter controls what submissions to search for.
type ListFilter struct {
	Submissions   *SubmissionFilter
	Users         *UserFilter
	Problems      *ProblemFilter
	OnlyAvailable bool
}

type SubmissionFilter struct {
	SubmissionIDs []int32
}

type UserFilter struct {
	UserIDs []int32
}

type ProblemFilter struct {
	ProblemIDs []int32
}

// A SubmissionList is a slice of Submissions.
type SubmissionList []*models.Submission

func (lists SubmissionList) ProblemIDs() []int32 {
	ids := make(map[int32]bool)
	for _, sub := range lists {
		ids[sub.ProblemID] = true
	}
	var res []int32
	for id, _ := range ids {
		res = append(res, id)
	}
	return res
}

// ListSubmissions searches for a list of submissions.
func ListSubmissions(ctx context.Context, args ListArgs, filterArgs ListFilter) (SubmissionList, error) {
	if len(filterArgs.Users.UserIDs) == 0 {
		return SubmissionList{}, nil
	}
	conn := db.Conn()
	var params []interface{}
	var filters []string
	joins := ""
	fields := ""
	if filterArgs.Submissions != nil {
		filters = append(filters, db.SetInParamInt(`submission.submission_id IN (%s)`, &params, filterArgs.Submissions.SubmissionIDs))
	}
	if filterArgs.Users != nil {
		filters = append(filters, db.SetInParamInt(`account_id IN (%s)`, &params, filterArgs.Users.UserIDs))
	}
	if filterArgs.Problems != nil {
		filters = append(filters, db.SetInParamInt(`problem_id IN (%s)`, &params, filterArgs.Problems.ProblemIDs))
	}
	if filterArgs.OnlyAvailable {
		filters = append(filters, "public_from <= current_timestamp")
		joins += " LEFT JOIN problem USING(problem_id) "
	}
	filter := ""
	if len(filters) > 0 {
		filter = fmt.Sprintf(" WHERE %s ", strings.Join(filters, " AND "))
	}
	if args.WithRun {
		joins += " LEFT JOIN submission_run ON current_run = submission_run_id "
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
	if args.WithAccounts {
		joins += " LEFT JOIN account USING (account_id) "
		fields += `
			,
			account.account_id "account.account_id",
			account.username "account.username"
`
	}
	if args.WithGroupRuns {
		fields += `,
			to_json(ARRAY(
				SELECT json_build_object(
					'testgroup_name', problem_testgroup.testgroup_name,
					'problem_testgroup_id', g.problem_testgroup_id,
					'time_usage_ms', g.time_usage_ms, 
					'score', g.score, 
					'verdict', g.verdict)
				FROM submission_group_run g
				LEFT JOIN problem_testgroup ON g.problem_testgroup_id = problem_testgroup.problem_testgroup_id
				WHERE g.submission_run_id = submission_run.submission_run_id)) "submission_run.group_runs"`
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
