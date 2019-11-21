package problems

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

type ListFilter struct {
	ShortName string
	ProblemId []int32
}

type TestOpt byte
type StmtOpt byte

const (
	TestsNone TestOpt = iota
	TestsSamples
	TestsAll

	StmtNone StmtOpt = iota
	StmtTitles
	StmtAll
)

type ListArgs struct {
	WithStatements StmtOpt
	WithTests      TestOpt
}

func List(ctx context.Context, args ListArgs, filter ListFilter) (ProblemList, error) {
	if filter.ShortName != "" && len(filter.ProblemId) != 0 {
		return nil, fmt.Errorf("Only one filter is allowed when listing problems")
	}
	conn := db.Conn()
	var probs ProblemList
	query, params := problemQuery(args, filter)
	if err := conn.SelectContext(ctx, &probs, query, params...); err != nil {
		return nil, err
	}
	if err := includeStatements(ctx, probs.AsMap(), args.WithStatements); err != nil {
		return nil, err
	}
	for _, p := range probs {
		if err := includeTests(ctx, p.CurrentVersion, args.WithTests); err != nil {
			return nil, err
		}
	}
	return probs, nil
}

func problemQuery(args ListArgs, filterArgs ListFilter) (string, []interface{}) {
	var params []interface{}
	filter := ""
	joins := ""
	fields := ""

	query := `
    SELECT
        problem_id, short_name, author, license,
        problem_version.problem_version_id, problem_version.time_limit_ms, problem_version.memory_limit_kb, problem_version.problem_id
        %s
    FROM problem
    LEFT JOIN problem_version ON current_version = problem_version.problem_version_id
    %s
    %s
    `
	if filterArgs.ShortName != "" {
		filter, params = db.SetParam("short_name = $%d", params, filterArgs.ShortName)
	} else if len(filterArgs.ProblemId) != 0 {
		filter, params = db.SetInParamInt("WHERE problem_id IN (%s)", params, filterArgs.ProblemId)
	}

	if args.WithTests == TestsAll {
		joins = "LEFT JOIN problem_output_validator USING(problem_version_id)"
		fields = `, file_hash(validator_source_zip) "problem_version.problem_output_validator.validator_source_zip.hash", file_url(validator_source_zip) "problem.version.problem_output_validator.validator_source_zip.url",
    validator_language_id "problem_output_validator.language_id"`
	}
	return fmt.Sprintf(query, fields, joins, filter), params
}

func includeTests(ctx context.Context, pv *models.ProblemVersion, opt TestOpt) error {
	filter := "WHERE problem_version_id = $1"
	if opt == TestsSamples {
		filter = filter + " AND public_visibility = true"
	}
	query := "SELECT problem_version_id, problem_testgroup_id, testgroup_name, public_visibility FROM problem_testgroup " + filter + " ORDER BY testgroup_name"
	var groups TestGroupList
	if err := db.Conn().SelectContext(ctx, &groups, query, pv.ProblemVersionId); err != nil {
		return err
	}
	pv.TestGroups = groups
	query = `
	SELECT
	problem_testgroup_id,
	problem_testcase_id,
	testcase_name,
	file_hash(input_file_hash) "input_file.hash",
	file_url(input_file_hash) "input_file.url",
	file_hash(output_file_hash) "output_file.hash",
	file_url(output_file_hash) "output_file.url"
	FROM problem_testcase
	NATURAL JOIN problem_testgroup ` + filter
	var tests []*models.TestCase
	if err := db.Conn().SelectContext(ctx, &tests, query, pv.ProblemVersionId); err != nil {
		return err
	}
	groupMap := groups.AsMap()
	for _, t := range tests {
		g := groupMap[t.TestGroupId]
		g.Tests = append(g.Tests, t)
	}
	return nil
}

func includeStatements(ctx context.Context, ps ProblemMap, arg StmtOpt) error {
	if len(ps) == 0 || arg == StmtNone {
		return nil
	}
	pids := ps.Ids()
	var cols string
	if arg == StmtAll {
		cols = ", html"
	}
	query, args, err := sqlx.In(fmt.Sprintf("SELECT problem_id, title, language %s FROM problem_statement WHERE problem_id IN (?);", cols), pids)
	if err != nil {
		return err
	}
	conn := db.Conn()
	query = conn.Rebind(query)
	var statements StatementList
	if err := conn.SelectContext(ctx, &statements, query, args...); err != nil {
		return err
	}
	statements.AddTo(ps)
	return nil
}
