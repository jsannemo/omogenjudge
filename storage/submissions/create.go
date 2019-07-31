package submissions

import (
  "context"

  "database/sql"

  "github.com/jsannemo/omogenjudge/storage/db"
)

// CreateFile inserts the given file into the database.
// If the file already existed, the URL will be updated instead.
func CreateSubmission(ctx context.Context, sub *Submission) error {
  tx, err := db.GetPool().BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
  if err != nil {
    return err
  }
  err = tx.QueryRow("INSERT INTO submission(account_id, problem_id, status) VALUES($1, $2, 'new') RETURNING submission_id", sub.AccountId, sub.ProblemId).Scan(&sub.SubmissionId)
  if err != nil {
    _ = tx.Rollback()
    return err
  }
  for _, file := range sub.Files {
    _, err = tx.Exec("INSERT INTO submission_file(submission_id, file_path, file_contents) VALUES($1, $2, $3)", sub.SubmissionId, file.Path, file.Contents)
    if err != nil {
      _ = tx.Rollback()
      return err
    }
  }
  err = tx.Commit()
  return err
}
