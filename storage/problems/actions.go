// Database actions related to problem storage.
package problems

import (
  "context"
  "errors"
  "database/sql"
  "html/template"

  "github.com/google/logger"

  "github.com/jsannemo/omogenjudge/storage/db"
  "github.com/jsannemo/omogenjudge/storage/files"
)

var ErrDuplicateProblemName = errors.New("Duplicate problem name")
var ErrNoSuchProblem = errors.New("No such problem")

func scanGroup(sc db.Scannable, testGroups TestGroupMap) error {
  group := &TestCaseGroup{}
  err := sc.Scan(&group.TestCaseGroupId, &group.Name, &group.PublicVisibility)
  if err != nil {
    return err
  }
  testGroups[group.TestCaseGroupId] = group
  return nil
}


func scanGroups(rows *sql.Rows) (TestGroupMap, error) {
  testGroups := make(TestGroupMap)
  for rows.Next() {
    err := scanGroup(rows, testGroups)
    if err != nil {
      return nil, err
    }
  }
  if rows.NextResultSet() {
    return nil, rows.Err()
  }
  return testGroups, nil
}

func scanTests(rows *sql.Rows, groups TestGroupMap) error {
  for rows.Next() {
    var groupId int32
    var infile, outfile files.StoredFile
    tc := &TestCase{}
    err := rows.Scan(&groupId, &tc.TestCaseId, &tc.Name, &infile.Hash, &infile.Url, &outfile.Hash, &outfile.Url)
    tc.InputFile = &infile
    tc.OutputFile = &outfile
    if err != nil {
      return err
    }
    groups[groupId].Tests = append(groups[groupId].Tests, tc)
  }
  if rows.NextResultSet() {
    return rows.Err()
  }
  return nil
}

func readTests(conn *sql.DB, problemId int32, testGroups TestGroupMap) error {
  rows, err := conn.Query(
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
    WHERE problem_id = $1
    ORDER BY testcase_name`, problemId)
  if err != nil {
    return err
  }
  defer rows.Close()
  err = scanTests(rows, testGroups)
  if err != nil {
    return err
  }
  return nil
}

func readTestGroups(conn *sql.DB, problemId int32) (TestGroupMap, error) {
  rows, err := conn.Query("SELECT problem_testgroup_id, testgroup_name, public_visibility FROM problem_testgroup WHERE problem_id = $1 ORDER BY testgroup_name", problemId)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  testGroups, err := scanGroups(rows)
  if err != nil {
    return nil, err
  }
  return testGroups, nil
}

func readTestData(conn *sql.DB, problemId int32) ([]*TestCaseGroup, error) {
  testGroups, err := readTestGroups(conn, problemId)
  if err != nil {
    return nil, err
  }
  err = readTests(conn, problemId, testGroups)
  if err != nil {
    return nil, err
  }
  var groups []*TestCaseGroup
  for _, group := range testGroups {
    groups = append(groups, group)
  }
  return groups, nil
}

// scanProblem scans a problem and appends it to the given problem map.
// If the problem already exists in the map, the statement loaded in the read is appended to it.
// The fields are assumed to be in the order:
// - problem.problem_id
// - problem.short_name
// - problem_statement.title
// - problem_statement.language
// - problem_statement.html (optional, if withStatements is true)
func scanProblem(sc db.Scannable, withStatement bool, problems ProblemMap) error {
  var id int32
  var name, title, language, html string
  var problem *Problem
  var err error
  if withStatement {
    err = sc.Scan(&id, &name, &title, &language, &html)
  } else {
    err = sc.Scan(&id, &name, &title, &language)
  }
  if err != nil {
    return err
  }
  problem, ok := problems[id]
  if !ok {
    problem = &Problem{id, name, nil, nil}
  }
  problem.Statements = append(problem.Statements, &ProblemStatement{language, title, template.HTML(html)})
  if !ok {
    problems[id] = problem
  }
  return nil
}

// scanProblems reads problems from rows until there are no more.
// Columns are assumed to be in the same order as required by scanProblem
func scanProblems(rows *sql.Rows, withStatements bool) (ProblemMap, error) {
  problems := make(map[int32]*Problem)
  for rows.Next() {
    err := scanProblem(rows, withStatements, problems)
    if err != nil {
      return nil, err
    }
  }
  if rows.NextResultSet() {
    return nil, rows.Err()
  }
  return problems, nil
}

// TODO: don't duplicate this method....
func GetProblemForJudging(ctx context.Context, id int32) (*Problem, error) {
  conn := db.GetPool()
  rows, err := conn.Query("SELECT problem_id, short_name, title, language, html FROM problem NATURAL JOIN problem_statement WHERE problem_id = $1", id)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  problems, err := scanProblems(rows, true)
  if err != nil {
    return nil, err
  }
  if len(problems) == 0 {
    return nil, ErrNoSuchProblem
  }

  for _, problem := range problems {
    groups, err := readTestData(conn, problem.ProblemId)
    if err != nil {
      return nil, err
    }
    problem.TestGroups = groups
    return problem, nil
  }
  logger.Fatalln("Invalid state")
  return nil, nil
}

// GetProblem reads a problem from the database with a given short name.
// Since short names are unique, this returns at most one problem.
// All statements are loaded, including their HTML.
func GetProblem(ctx context.Context, shortName string, includeTests bool) (*Problem, error) {
  conn := db.GetPool()
  rows, err := conn.Query("SELECT problem_id, short_name, title, language, html FROM problem NATURAL JOIN problem_statement WHERE short_name = $1", shortName)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  problems, err := scanProblems(rows, true)
  if err != nil {
    return nil, err
  }
  if len(problems) == 0 {
    return nil, ErrNoSuchProblem
  }

  for _, problem := range problems {
    if includeTests {
      groups, err := readTestData(conn, problem.ProblemId)
      if err != nil {
        return nil, err
      }
      problem.TestGroups = groups
    }
    return problem, nil
  }
  logger.Fatalln("Invalid state")
  return nil, nil
}

// ListProblems reads all problems from the database.
// Titles of all languages are read from the problems as well.
func ListProblems(ctx context.Context) (ProblemMap, error) {
  conn := db.GetPool()
  rows, err := conn.Query("SELECT problem_id, short_name, title, language FROM problem NATURAL JOIN problem_statement ORDER BY short_name")
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  return scanProblems(rows, false)
}

func insertTests(ctx context.Context, group *TestCaseGroup, tx *sql.Tx) error {
  for _, tc := range group.Tests {
    err := tx.QueryRow(
      "INSERT INTO problem_testcase(problem_testgroup_id, testcase_name, input_file_hash.hash, output_file_hash.hash) VALUES($1, $2, $3, $4) RETURNING problem_testcase_id",
      group.TestCaseGroupId, tc.Name, tc.InputFile.Hash, tc.OutputFile.Hash).Scan(&tc.TestCaseId)
    if err != nil {
      return err
    }
  }
  return nil
}

func insertTestGroups(ctx context.Context, problem *Problem, tx *sql.Tx) error {
  for _, group := range problem.TestGroups {
    err := tx.QueryRow(
      "INSERT INTO problem_testgroup(problem_id, testgroup_name, public_visibility) VALUES($1, $2, $3) RETURNING problem_testgroup_id",
      problem.ProblemId, group.Name, group.PublicVisibility).Scan(&group.TestCaseGroupId)
    if err != nil {
      return err
    }
    err = insertTests(ctx, group, tx)
    if err != nil {
      return err
    }
  }
  return nil
}

func insertStatements(ctx context.Context, problem *Problem, tx *sql.Tx) error {
  for _, statement := range problem.Statements {
    _, err := tx.Exec(
      "INSERT INTO problem_statement(problem_id, language, title, html) VALUES($1, $2, $3, $4)",
      problem.ProblemId, statement.Language, statement.Title, statement.Html)
    if err != nil {
      return err
    }
  }
  return nil
}

func insertProblem(ctx context.Context, problem *Problem, tx *sql.Tx) error {
  return tx.QueryRow("INSERT INTO problem(short_name) VALUES($1) RETURNING problem_id", problem.ShortName).Scan(&problem.ProblemId)
}

func errorWithRollback(err error, tx *sql.Tx) error {
  if err != nil {
    if err := tx.Rollback(); err != nil {
      return err
    }
  }
  return err
}

// CreateProblem inserts the given problem into the database.
// Insertion is atomic, so that if insertion of any aspect of the problem fails, none of it is inserted.
func CreateProblem(ctx context.Context, problem *Problem) error {
  tx, err := db.GetPool().BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
  if err != nil {
    return err
  }
  if err := errorWithRollback(insertProblem(ctx, problem, tx), tx); err != nil {
    if db.PgErrCode(err) == db.UniquenessViolation {
      return ErrDuplicateProblemName
    } else {
      return err
    }
  }
  if err := errorWithRollback(insertStatements(ctx, problem, tx), tx); err != nil {
    return err
  }
  if err := errorWithRollback(insertTestGroups(ctx, problem, tx), tx); err != nil {
    return err
  }
  if err := tx.Commit(); err != nil {
    return err
  }
  return nil
}
