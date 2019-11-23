package contests

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

var ErrAccountInTeam = errors.New("an account was already registered in the contest")

// CreateTeam persists a team with its members in the database.
func CreateTeam(ctx context.Context, team *models.Team) error {
	return db.InTransaction(ctx, func(tx *sqlx.Tx) error {
		if len(team.Members) == 0 {
			return fmt.Errorf("team was empty")
		}
		query := `
			INSERT INTO
			  team(contest_id, team_name, virtual, unofficial, team_data)
			VALUES($1, $2, false, false, '{}')
			RETURNING team_id`
		if err := tx.QueryRowContext(ctx, query, team.ContestID, team.TeamName).Scan(&team.TeamID); err != nil {
			return fmt.Errorf("failed create contest query: %v", err)
		}
		for _, tm := range team.Members {
			tm.TeamID = team.TeamID
			if err := insertMember(ctx, team.ContestID, tm, tx); err != nil {
				return fmt.Errorf("failed insert contest problem: %v", err)
			}
		}
		return nil
	})
}

func insertMember(ctx context.Context, contestID int32, member *models.TeamMember, tx *sqlx.Tx) error {
	rows := 0
	if err := tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM team_member LEFT JOIN team USING(team_id) WHERE contest_id = $1 AND account_id = $2`,
		contestID, member.AccountID).Scan(&rows); err != nil {
		return err
	}
	if rows != 0 {
		return ErrAccountInTeam
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO team_member(team_id, account_id) VALUES ($1, $2)`,
		member.TeamID, member.AccountID); err != nil {
		return err
	}
	return nil
}

