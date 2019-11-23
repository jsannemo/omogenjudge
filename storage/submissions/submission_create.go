package submissions

import (
	"context"

	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"

	"github.com/jmoiron/sqlx"
)

// CreateSubmission persists a new submission.
func CreateSubmission(ctx context.Context, sub *models.Submission, problemVersion int32) error {
	err := db.InTransaction(ctx, func(tx *sqlx.Tx) error {
		if err := CreateSubmissionTx(ctx, tx, sub); err != nil {
			return err
		}
		for _, file := range sub.Files {
			file.SubmissionID = sub.SubmissionID
			if err := createFileTx(ctx, tx, file); err != nil {
				return err
			}
		}
		for _, run := range sub.Runs {
			run.SubmissionID = sub.SubmissionID
			if err := CreateRunTx(ctx, run, tx); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// CreateSubmissionTx persists a new submission within a transaction.
func CreateSubmissionTx(ctx context.Context, tx *sqlx.Tx, sub *models.Submission) error {
	query := `
		INSERT INTO
			submission(account_id, problem_id, language)
		VALUES($1, $2, $3)
		RETURNING submission_id`
	return tx.QueryRowContext(ctx, query, sub.AccountID, sub.ProblemID, sub.Language).Scan(&sub.SubmissionID)
}

func createFileTx(ctx context.Context, tx *sqlx.Tx, file *models.SubmissionFile) error {
	query := `INSERT INTO submission_file(submission_id, file_path, file_contents) VALUES($1, $2, $3)`
	_, err := tx.ExecContext(ctx, query, file.SubmissionID, file.Path, file.Contents)
	return err
}
