package contests

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"

	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

// ErrShortNameExists is returned when the given shortname already exists.
var ErrShortNameExists = errors.New("the shortname is in use")

// CreateContest persists a contest in the database.
func CreateContest(ctx context.Context, contest *models.Contest) error {
	return db.InTransaction(ctx, func(tx *sqlx.Tx) error {
		query := `
    INSERT INTO
      contest(short_name, host_name, start_time, selection_window_end_time, duration, title, hidden_scoreboard)
    VALUES($1, $2, $3, $4, $5 * '1 microsecond'::interval, $6, $7)
    RETURNING contest_id`
		if err := tx.QueryRowContext(ctx, query, contest.ShortName, contest.HostName, contest.StartTime, contest.FlexibleEndTime, contest.Duration/time.Microsecond, contest.Title, contest.HiddenScoreboard).Scan(&contest.ContestID); err != nil {
			if db.PgErrCode(err) == db.UniquenessViolation {
				return ErrShortNameExists
			} else {
				return fmt.Errorf("failed create contest query: %v", err)
			}
		}
		for _, p := range contest.Problems {
			p.ContestID = contest.ContestID
		}
		if err := insertProblems(ctx, contest, tx); err != nil {
			return fmt.Errorf("failed insert contest problem: %v", err)
		}
		return nil
	})
}

func insertProblems(ctx context.Context, contest *models.Contest, tx *sqlx.Tx) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM contest_problem WHERE contest_id = $1`, contest.ContestID); err != nil {
		return fmt.Errorf("failed clearing old problems: %v", err)
	}
	for _, problem := range contest.Problems {
		if _, err := tx.ExecContext(ctx, `INSERT INTO contest_problem(contest_id, problem_id, label) VALUES ($1, $2 ,$3)`,
			contest.ContestID, problem.ProblemID, problem.Label); err != nil {
			return err
		}
	}
	return nil
}
