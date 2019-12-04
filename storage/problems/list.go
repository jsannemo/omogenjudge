package problems

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

// A ListFilter filters the problems to search for.
// Only one filter may be set.
type ListFilter struct {
	ShortName string
	Problems  *ProblemFilter
	ContestID int32
}

type ProblemFilter struct {
	ProblemID []int32
}

func Problems(problemIDs ...int32) ListFilter {
	return ListFilter{Problems: &ProblemFilter{problemIDs}}
}

type TestOpt byte
type StmtOpt byte

const (
	TestsNone TestOpt = iota
	// Only load sample test groups.
	TestsSamples
	// Load only the test data structure.
	TestsGroups
	// Load test data and validators.
	TestsAll

	StmtNone StmtOpt = iota
	// Include only titles.
	StmtTitles
	// Include titles and HTML statement.
	StmtAll
)

type ListArgs struct {
	WithStatements StmtOpt
	WithTests      TestOpt
}

func List(ctx context.Context, args ListArgs, filter ListFilter) (models.ProblemList, error) {
	filters := 0
	if filter.ShortName != "" {
		filters++
	}
	if filter.Problems != nil {
		if len(filter.Problems.ProblemID) == 0 {
			return nil, nil
		}
		filters++
	}
	if filter.ContestID != 0 {
		filters++
	}
	if filters > 1 {
		return nil, fmt.Errorf("only one filter is allowed when listing problems")
	}
	conn := db.Conn()
	var probs models.ProblemList
	query, params := problemQuery(args, filter)
	if err := conn.SelectContext(ctx, &probs, query, params...); err != nil {
		return nil, fmt.Errorf("list query failed: %v", err)
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
        problem.problem_id, short_name, author, license,
        problem_version.problem_version_id "problem_version.problem_version_id",
   	    problem_version.time_limit_ms "problem_version.time_limit_ms",
		problem_version.memory_limit_kb "problem_version.memory_limit_kb",
		problem_version.problem_id "problem_version.problem_id"
        %s
    FROM problem
    LEFT JOIN problem_version ON current_version = problem_version.problem_version_id
    %s
    %s
    `
	if filterArgs.ShortName != "" {
		filter = db.SetParam("WHERE short_name = $%d", &params, filterArgs.ShortName)
	} else if filterArgs.Problems != nil {
		filter = db.SetInParamInt("WHERE problem.problem_id IN (%s)", &params, filterArgs.Problems.ProblemID)
	} else if filterArgs.ContestID != 0 {
		joins += " LEFT JOIN contest_problem USING (problem_id) "
		filter = db.SetParam("WHERE contest_id = $%d", &params, filterArgs.ContestID)
	}

	if args.WithTests == TestsAll {
		joins += " LEFT JOIN problem_output_validator USING(problem_version_id) "
		fields = `, validator_source_zip "problem_version.validator.validator_source_zip.hash",
                    file_url(validator_source_zip) "problem_version.validator.validator_source_zip.url",
                    validator_language_id "problem_version.validator.language_id"`
	}
	return fmt.Sprintf(query, fields, joins, filter), params
}

func includeTests(ctx context.Context, pv *models.ProblemVersion, opt TestOpt) error {
	if opt == TestsNone {
		return nil
	}
	if err := includeTestGroups(ctx, pv, opt); err != nil {
		return err
	}
	if opt == TestsGroups {
		return nil
	}
	filter := "WHERE problem_version_id = $1"
	if opt == TestsSamples {
		filter = filter + " AND public_visibility = true"
	}
	query := `
	SELECT
	problem_testgroup_id,
	problem_testcase_id,
	testcase_name,
	input_file_hash "input_file.hash",
	file_url(input_file_hash) "input_file.url",
	output_file_hash "output_file.hash",
	file_url(output_file_hash) "output_file.url"
	FROM problem_testcase
	NATURAL JOIN problem_testgroup ` + filter
	var tests []*models.TestCase
	if err := db.Conn().SelectContext(ctx, &tests, query, pv.ProblemVersionID); err != nil {
		return err
	}
	groupMap := pv.TestGroups.AsMap()
	for _, t := range tests {
		g := groupMap[t.TestGroupID]
		g.Tests = append(g.Tests, t)
	}
	return nil
}

func includeTestGroups(ctx context.Context, pv *models.ProblemVersion, opt TestOpt) error {
	filter := "WHERE problem_version_id = $1"
	if opt == TestsSamples {
		filter = filter + " AND public_visibility = true"
	}
	query := "SELECT problem_version_id, problem_testgroup_id, testgroup_name, score, public_visibility, output_validator_flags FROM problem_testgroup " +
		filter +
		" ORDER BY public_visibility DESC, testgroup_name ASC"
	var groups models.TestGroupList
	if err := db.Conn().SelectContext(ctx, &groups, query, pv.ProblemVersionID); err != nil {
		return err
	}
	pv.TestGroups = groups
	return nil
}

func includeStatements(ctx context.Context, ps models.ProblemMap, arg StmtOpt) error {
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
	var statements models.StatementList
	if err := conn.SelectContext(ctx, &statements, query, args...); err != nil {
		return err
	}
	for _, s := range statements {
		p := ps[s.ProblemID]
		p.Statements = append(p.Statements, s)
	}
	return nil
}
