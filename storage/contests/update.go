package contests

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

// UpdateContest updates a contest within a transaction.
func UpdateContest(ctx context.Context, contest *models.Contest) error {
	if err := db.InTransaction(ctx, func(tx *sqlx.Tx) error {
		query := `
    UPDATE contest
	SET 
	    short_name = $1,
	    host_name = $2,
	    start_time = $3,
	    duration = $4 * '1 microsecond'::interval,
	    title = $5,
	    hidden_scoreboard = $6,
	    selection_window_end_time = $7
    WHERE short_name = $1
	RETURNING contest_id`
		if err := tx.QueryRowContext(ctx, query, contest.ShortName, contest.HostName, contest.StartTime, contest.Duration/time.Microsecond, contest.Title, contest.HiddenScoreboard, contest.FlexibleEndTime).Scan(&contest.ContestID); err != nil {
			return fmt.Errorf("failed update contest query: %v", err)
		}
		for _, p := range contest.Problems {
			p.ContestID = contest.ContestID
		}
		if err := insertProblems(ctx, contest, tx); err != nil {
			return fmt.Errorf("failed insert contest problem: %v", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
