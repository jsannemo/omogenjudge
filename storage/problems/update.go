package problems

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

// UpdateProblem updates a contest.
func UpdateProblem(ctx context.Context, problem *models.Problem) error {
	if err := db.InTransaction(ctx, func(tx *sqlx.Tx) error {
		query := `
			UPDATE problem
			SET 
			    author = $2,
			    source = $3,
			    license = $4
			WHERE short_name = $1
			RETURNING problem_id`
		if err := tx.QueryRowContext(ctx, query, problem.ShortName, problem.Author, problem.Source, problem.License).Scan(&problem.ProblemID); err != nil {
			return fmt.Errorf("failed update problem query: %v", err)
		}

		// We insert this after updating the problem to make sure we know the problem ID.
		if problem.CurrentVersion != nil && problem.CurrentVersion.ProblemVersionID == 0 {
			problem.CurrentVersion.ProblemID = problem.ProblemID
			if err := insertProblemVersion(ctx, problem.CurrentVersion, tx); err != nil {
				return err
			}
		}
		if err := setCurrentVersion(ctx, problem, tx); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "DELETE FROM problem_statement WHERE problem_id = $1", problem.ProblemID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "DELETE FROM problem_statement_file WHERE problem_id = $1", problem.ProblemID); err != nil {
			return err
		}
		if err := insertStatementFiles(ctx, problem, tx); err != nil {
			return err
		}
		for _, s := range problem.Statements {
			s.ProblemID = problem.ProblemID
			if err := insertStatement(ctx, s, tx); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
