// Database actions relating to stored files.
package submissions

import (
  "context"

  "database/sql"

  "github.com/jsannemo/omogenjudge/storage/db"
)

func UnjudgedIds() ([]int32, error) {
  conn := db.GetPool()
  rows, err := conn.Query("SELECT submission_id FROM submission WHERE status != 'successful' ORDER BY submission_id ASC")
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var ids []int32
  for rows.Next() {
    var id int32
    if err := rows.Scan(&id); err != nil {
      return nil, err
    }
    ids = append(ids, id)
  }
  return ids, nil
}

func readSubmission(ctx context.Context, conn *sql.DB, sub *Submission) error {
  return conn.QueryRow("SELECT problem_id, account_id FROM submission WHERE submission_id = $1", sub.SubmissionId).Scan(
    &sub.ProblemId, &sub.AccountId)
}

func readSubmissionFiles(ctx context.Context, conn *sql.DB, sub *Submission) error {
  rows, err := conn.Query("SELECT file_path, file_contents FROM submission_file WHERE submission_id = $1", sub.SubmissionId)
  if err != nil {
    return err
  }
  defer rows.Close()

  for rows.Next() {
    subFile := &SubmissionFile{}
    if err := rows.Scan(&subFile.Path, &subFile.Contents); err != nil {
      return err
    }
    sub.Files = append(sub.Files, subFile)
  }
  return nil
}

func GetSubmission(ctx context.Context, id int32) (*Submission, error) {
  conn := db.GetPool()
  submission := &Submission{SubmissionId: id}
  if err := readSubmission(ctx, conn, submission); err != nil {
    return nil, err
  }
  if err := readSubmissionFiles(ctx, conn, submission); err != nil {
    return nil, err
  }
  return submission, nil
}

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
