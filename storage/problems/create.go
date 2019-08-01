package problems

import (
  "context"
  "errors"

  "github.com/jsannemo/omogenjudge/storage/db"
  "github.com/jsannemo/omogenjudge/storage/models"
  "github.com/jmoiron/sqlx"
)

var ErrDuplicateProblemName = errors.New("Duplicate problem name")

func insertTest(ctx context.Context, tc *models.TestCase, tx *sqlx.Tx) error {
  return tx.QueryRowContext(
      ctx,
      "INSERT INTO problem_testcase(problem_testgroup_id, testcase_name, input_file_hash.hash, output_file_hash.hash) VALUES($1, $2, $3, $4) RETURNING problem_testcase_id",
      tc.TestGroupId, tc.Name, tc.InputFile.Hash, tc.OutputFile.Hash).Scan(&tc.TestCaseId)
}

func insertTestGroup(ctx context.Context, tg *models.TestGroup, tx *sqlx.Tx) error {
  if err := tx.QueryRowContext(
    ctx,
    "INSERT INTO problem_testgroup(problem_id, testgroup_name, public_visibility) VALUES($1, $2, $3) RETURNING problem_testgroup_id",
    tg.ProblemId, tg.Name, tg.PublicVisibility).Scan(&tg.TestGroupId); err != nil {
    return err
  }
  for _, tc := range tg.Tests {
    tc.TestGroupId = tg.TestGroupId
    if err := insertTest(ctx, tc, tx); err != nil {
      return err
    }
  }
  return nil
}

func insertStatement(ctx context.Context, s *models.ProblemStatement, tx *sqlx.Tx) error {
  if _, err := tx.ExecContext(
    ctx,
    "INSERT INTO problem_statement(problem_id, language, title, html) VALUES($1, $2, $3, $4)",
    s.ProblemId, s.Language, s.Title, s.Html); err != nil {
    return err
  }
  return nil
}

func insertProblem(ctx context.Context, problem *models.Problem, tx *sqlx.Tx) error {
  return tx.QueryRowContext(ctx, "INSERT INTO problem(short_name) VALUES($1) RETURNING problem_id", problem.ShortName).Scan(&problem.ProblemId)
}

func Create(ctx context.Context, p *models.Problem) error {
  err := db.InTransaction(ctx, func (tx *sqlx.Tx) error {
    if err := insertProblem(ctx, p, tx); err != nil {
      if db.PgErrCode(err) == db.UniquenessViolation {
        return ErrDuplicateProblemName
      } else {
        return err
      }
    }
    for _, s := range p.Statements {
      s.ProblemId = p.ProblemId
      if err := insertStatement(ctx, s, tx); err != nil {
        return err
      }
    }
    for _, tg := range p.TestGroups {
      tg.ProblemId = p.ProblemId
      if err := insertTestGroup(ctx, tg, tx); err != nil {
        return err
      }
    }
    return nil
  })
  return err
}
