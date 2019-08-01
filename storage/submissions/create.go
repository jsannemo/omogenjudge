package submissions

import (
  "context"

  "github.com/jsannemo/omogenjudge/storage/db"
  "github.com/jsannemo/omogenjudge/storage/models"

  "github.com/jmoiron/sqlx"
)

func Create(ctx context.Context, sub *models.Submission) error {
  err := db.InTransaction(ctx, func (tx *sqlx.Tx) error {
    if err := tx.QueryRowContext(ctx, "INSERT INTO submission(account_id, problem_id, status) VALUES($1, $2, 'new') RETURNING submission_id", sub.AccountId, sub.ProblemId).Scan(&sub.SubmissionId); err != nil {
      return err
    }
    for _, file := range sub.Files {
      if _, err := tx.ExecContext(ctx, "INSERT INTO submission_file(submission_id, file_path, file_contents) VALUES($1, $2, $3)", sub.SubmissionId, file.Path, file.Contents); err != nil {
        return err
      }
    }
    return nil
  })
  return err
}
