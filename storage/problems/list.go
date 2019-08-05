// Database actions related to problem storage.
package problems

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

type ListFilter struct {
	ShortName   string
	ProblemId   int32
	Submissions []*models.Submission
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

func problemQuery(args ListArgs, filter ListFilter) (string, []interface{}) {
	var params []interface{}
	var filterSegs []string

	if filter.ShortName != "" {
		params = append(params, filter.ShortName)
		filterSegs = append(filterSegs, fmt.Sprintf("short_name = $%d", len(params)))
	}
	if filter.ProblemId != 0 {
		params = append(params, filter.ProblemId)
		filterSegs = append(filterSegs, fmt.Sprintf("problem_id = $%d", len(params)))
	}

	filterStr := ""
	if len(filterSegs) != 0 {
		filterStr = "WHERE " + strings.Join(filterSegs, " AND ")
	}

	fields := ""
	joins := ""
	if args.WithTests == TestsAll {
		joins = "LEFT JOIN problem_output_validator USING(problem_id)"
		fields = `, file_hash(validator_source_zip) "problem_output_validator.validator_source_zip.hash", file_url(validator_source_zip) "problem_output_validator.validator_source_zip.url",
    validator_language_id "problem_output_validator.language_id"`
	}
	return fmt.Sprintf("SELECT problem_id, short_name, author, license, time_limit_ms, memory_limit_kb "+fields+" FROM problem %s %s ORDER BY short_name", joins, filterStr), params
}

func List(ctx context.Context, args ListArgs, filter ListFilter) models.ProblemList {
	conn := db.Conn()
	var probs models.ProblemList
	query, params := problemQuery(args, filter)
	if err := conn.SelectContext(ctx, &probs, query, params...); err != nil {
		panic(err)
	}
	if args.WithStatements != StmtNone {
		includeStatements(ctx, probs.AsMap(), args.WithStatements)
	}
	for _, p := range probs {
		if args.WithTests != TestsNone {
			includeTests(ctx, p, args.WithTests)
		}
	}
	return probs
}

func includeTests(ctx context.Context, p *models.Problem, opt TestOpt) {
	filter := "WHERE problem_id = $1"
	if opt == TestsSamples {
		filter = filter + " AND public_visibility = true"
	}
	query := "SELECT problem_id, problem_testgroup_id, testgroup_name, public_visibility FROM problem_testgroup " + filter + " ORDER BY testgroup_name"
	var groups models.TestGroupList
	if err := db.Conn().SelectContext(ctx, &groups, query, p.ProblemId); err != nil {
		panic(err)
	}
	p.TestGroups = groups
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
	if err := db.Conn().SelectContext(ctx, &tests, query, p.ProblemId); err != nil {
		panic(err)
	}
	groupMap := groups.AsMap()
	for _, t := range tests {
		g := groupMap[t.TestGroupId]
		g.Tests = append(g.Tests, t)
	}
}

func includeStatements(ctx context.Context, ps models.ProblemMap, arg StmtOpt) {
	if len(ps) == 0 {
		return
	}
	pids := ps.Ids()
	var cols string
	if arg == StmtAll {
		cols = ", html"
	}
	query, args, err := sqlx.In(fmt.Sprintf("SELECT problem_id, title, language %s FROM problem_statement WHERE problem_id IN (?);", cols), pids)
	if err != nil {
		panic(err)
	}
	conn := db.Conn()
	query = conn.Rebind(query)
	var statements models.StatementList
	if err := conn.SelectContext(ctx, &statements, query, args...); err != nil {
		panic(err)
	}
	statements.AddTo(ps)
}
