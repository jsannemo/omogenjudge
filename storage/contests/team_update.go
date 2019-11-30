package contests

import (
	"context"

	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

func UpdateTeam(ctx context.Context, team *models.Team) error {
	if _, err := db.Conn().ExecContext(ctx, `
    UPDATE team
	SET 
	    team_name = $1,
	    start_time = $2
    WHERE team_id = $3`,
		team.TeamName,
		team.StartTime,
		team.TeamID); err != nil {
		return err
	}
	return nil
}
