// Database actions related to problem storage.
package problems

import (
  "context"
  "database/sql"
  "fmt"
  "strings"

  "github.com/jsannemo/omogenjudge/storage/db"
  "github.com/jsannemo/omogenjudge/storage/files"
)

type ListFilter struct {
  ShortName string
  ProblemId int32
}

type ListArgs struct {
  WithTitles bool
  WithStatements bool
  WithTests bool
  WithSamples bool
}

func problemQuery(args ListArgs, filter ListFilter) (string, []interface{}) {
  var joins []string
  cols := []string{"problem_id", "short_name"}
  if args.WithTitles || args.WithStatements {
    cols = append(cols, "title", "language")
    joins = append(joins, "problem_statement")
  }
  if args.WithStatements {
    cols = append(cols, "html")
  }

  filterStr := ""
  var params []interface{}
  if filter.ShortName != "" {
    filterStr = "WHERE short_name = $1"
    params = append(params, filter.ShortName)
  } else {
  if filter.ProblemId != 0 {
    filterStr = "WHERE problem_id = $1"
    params = append(params, filter.ProblemId)
  }
  }

  joinStr := ""
  if len(joins) != 0 {
    joinStr = fmt.Sprintf("NATURAL JOIN %s", strings.Join(joins, ","))
  }
  return fmt.Sprintf("SELECT %s FROM problem %s %s ORDER BY short_name", strings.Join(cols, ","), joinStr, filterStr), params
}

func ListProblems(ctx context.Context, args ListArgs, filter ListFilter) ([]*Problem, error) {
  conn := db.GetPool()
  query, params := problemQuery(args, filter)
  rows, err := conn.Query(query, params...)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  problems := make(map[int32]*Problem)
  for rows.Next() {
    err := scanProblem(rows, problems, args)
    if err != nil {
      return nil, err
    }
  }
	if err := rows.Err(); err != nil {
    return nil, err
	}
  var pList []*Problem
  for _, problem := range problems {
    if args.WithTests || args.WithSamples {
      if err := includeTestData(conn, problem, args); err != nil {
        return nil, err
      }
    }
    pList = append(pList, problem)
  }
  return pList, nil
}

func scanProblem(sc db.Scannable, problems ProblemMap, args ListArgs) error {
  var id int32
  statement := ProblemStatement{}
  var shortName string

  ptrs := []interface{}{&id, &shortName}
  if args.WithTitles || args.WithStatements {
    ptrs = append(ptrs, &statement.Title, &statement.Language)
  }
  if args.WithStatements {
    ptrs = append(ptrs, &statement.Html)
  }
  if err := sc.Scan(ptrs...); err != nil {
    return err
  }

  p, ok := problems[id]
  if !ok {
    p = &Problem{ProblemId: id, ShortName: shortName}
  }
  if args.WithTitles || args.WithStatements {
    p.Statements = append(p.Statements, &statement)
  }
  if !ok {
    problems[id] = p
  }
  return nil
}

func includeTestData(conn *sql.DB, problem *Problem, args ListArgs) error {
  filter := "WHERE problem_id = $1"
  if args.WithSamples && !args.WithTests {
    filter = filter + " AND public_visibility = true"
  }
  query := fmt.Sprintf("SELECT problem_testgroup_id, testgroup_name, public_visibility FROM problem_testgroup %s ORDER BY testgroup_name", filter)
  rows, err := conn.Query(query, problem.ProblemId)
  if err != nil {
    return err
  }
  defer rows.Close()
  testGroups := make(TestGroupMap)
  for rows.Next() {
    err := scanGroup(rows, testGroups)
    if err != nil {
      return err
    }
  }
	if err := rows.Err(); err != nil {
    return err
	}

  if err := includeTests(conn, testGroups, problem, args); err != nil {
    return err
  }
  for _, group := range testGroups {
    problem.TestGroups = append(problem.TestGroups, group)
  }
  return nil
}

func scanGroup(sc db.Scannable, testGroups TestGroupMap) error {
  group := &TestCaseGroup{}
  err := sc.Scan(&group.TestCaseGroupId, &group.Name, &group.PublicVisibility)
  if err != nil {
    return err
  }
  testGroups[group.TestCaseGroupId] = group
  return nil
}

func includeTests(conn *sql.DB, testGroups TestGroupMap, problem *Problem, args ListArgs) error {
  filter := "WHERE problem_id = $1"
  if args.WithSamples && !args.WithTests {
    filter = filter + " AND public_visibility = true"
  }
  rows, err := conn.Query(
    fmt.Sprintf(
    `
    SELECT
      problem_testgroup_id,
      problem_testcase_id,
      testcase_name,
      file_hash(input_file_hash),
      file_url(input_file_hash),
      file_hash(output_file_hash),
      file_url(output_file_hash)
    FROM problem_testcase
    NATURAL JOIN problem_testgroup
    %s
    ORDER BY testcase_name`, filter), problem.ProblemId)
  if err != nil {
    return err
  }
  defer rows.Close()
  for rows.Next() {
    var groupId int32
    var infile, outfile files.StoredFile
    tc := TestCase{}
    err := rows.Scan(&groupId, &tc.TestCaseId, &tc.Name, &infile.Hash, &infile.Url, &outfile.Hash, &outfile.Url)
    tc.InputFile = &infile
    tc.OutputFile = &outfile
    if err != nil {
      return err
    }
    testGroups[groupId].Tests = append(testGroups[groupId].Tests, &tc)
  }
	if err := rows.Err(); err != nil {
    return err
	}
  return nil
}
