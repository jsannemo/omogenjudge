package problems

import (
  "context"
  "errors"
  "database/sql"

  "github.com/jsannemo/omogenjudge/storage/db"
)

var ErrDuplicateProblemName = errors.New("Duplicate problem name")

func insertTests(ctx context.Context, group *TestCaseGroup, tx *sql.Tx) error {
  for _, tc := range group.Tests {
    if err := tx.QueryRowContext(
      ctx,
      "INSERT INTO problem_testcase(problem_testgroup_id, testcase_name, input_file_hash.hash, output_file_hash.hash) VALUES($1, $2, $3, $4) RETURNING problem_testcase_id",
      group.TestCaseGroupId, tc.Name, tc.InputFile.Hash, tc.OutputFile.Hash).Scan(&tc.TestCaseId); err != nil {
      return err
    }
  }
  return nil
}

func insertTestGroups(ctx context.Context, problem *Problem, tx *sql.Tx) error {
  for _, group := range problem.TestGroups {
    if err := tx.QueryRowContext(
      ctx,
      "INSERT INTO problem_testgroup(problem_id, testgroup_name, public_visibility) VALUES($1, $2, $3) RETURNING problem_testgroup_id",
      problem.ProblemId, group.Name, group.PublicVisibility).Scan(&group.TestCaseGroupId); err != nil {
      return err
    }
    if err := insertTests(ctx, group, tx); err != nil {
      return err
    }
  }
  return nil
}

func insertStatements(ctx context.Context, problem *Problem, tx *sql.Tx) error {
  for _, statement := range problem.Statements {
    if _, err := tx.ExecContext(
      ctx,
      "INSERT INTO problem_statement(problem_id, language, title, html) VALUES($1, $2, $3, $4)",
      problem.ProblemId, statement.Language, statement.Title, statement.Html); err != nil {
      return err
    }
  }
  return nil
}

func insertProblem(ctx context.Context, problem *Problem, tx *sql.Tx) error {
  return tx.QueryRow("INSERT INTO problem(short_name) VALUES($1) RETURNING problem_id", problem.ShortName).Scan(&problem.ProblemId)
}

// CreateProblem inserts the given problem into the database.
// Insertion is atomic, so that if insertion of any aspect of the problem fails, none of it is inserted.
func CreateProblem(ctx context.Context, problem *Problem) error {
  tx, err := db.GetPool().BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
  if err != nil {
    return err
  }
  if err := insertProblem(ctx, problem, tx); err != nil {
    tx.Rollback()
    if db.PgErrCode(err) == db.UniquenessViolation {
      return ErrDuplicateProblemName
    } else {
      return err
    }
  }
  if err := insertStatements(ctx, problem, tx); err != nil {
    tx.Rollback()
    return err
  }
  if err := insertTestGroups(ctx, problem, tx); err != nil {
    tx.Rollback()
    return err
  }
  if err := tx.Commit(); err != nil {
    tx.Rollback()
    return err
  }
  return nil
}
