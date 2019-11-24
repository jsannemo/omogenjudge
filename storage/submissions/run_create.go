package submissions

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

// CreateRun creates a new run for a submission.
func CreateRun(ctx context.Context, run *models.SubmissionRun) error {
	err := db.InTransaction(ctx, func(tx *sqlx.Tx) error {
		return CreateRunTx(ctx, run, tx)
	})
	return err
}

// CreateRun creates a new submission run in a transaction.
func CreateRunTx(ctx context.Context, run *models.SubmissionRun, tx *sqlx.Tx) error {
	query := `
			INSERT INTO
			    submission_run(submission_id, problem_version_id, status, time_usage_ms, score, verdict)
			VALUES($1, $2, $3, $4, $5, $6)
			RETURNING submission_run_id`
	if err := tx.QueryRowContext(ctx, query, run.SubmissionID, run.ProblemVersionID, run.Status, run.TimeUsageMS, run.Score, run.Verdict).Scan(&run.SubmissionRunID); err != nil {
		return fmt.Errorf("failed inserting submission run: %v", err)
	}
	return nil
}
