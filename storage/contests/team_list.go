package contests

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"

	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

// A TeamFilter controls the filtering behaviour of ListTeams.
type TeamFilter struct {
	ContestID int32
	AccountID int32
}

// ListTeams returns a list of teams.
func ListTeams(ctx context.Context, filter TeamFilter) (models.TeamList, error) {
	var teams models.TeamList
	if err := db.InTransaction(ctx, func(tx *sqlx.Tx) error {
		query, params := teamListQuery(filter)
		if err := tx.SelectContext(ctx, &teams, query, params...); err != nil {
			return fmt.Errorf("failed team list query: %v", err)
		}
		if err := includeTeams(ctx, teams, tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return teams, nil
}

func includeTeams(ctx context.Context, teams models.TeamList, tx *sqlx.Tx) error {
	if len(teams) == 0 {
		return nil
	}
	var teamIDs []int32
	for _, t := range teams {
		teamIDs = append(teamIDs, t.TeamID)
	}
	var params []interface{}
	query := `
		SELECT
			team_id, account_id,
			account.account_id "account.account_id",
			account.username "account.username"
		FROM team
		LEFT JOIN team_member USING(team_id)
		LEFT JOIN account USING(account_id)
		WHERE team.team_id IN (%s)`
	query = db.SetInParamInt(query, &params, teamIDs)
	var memberList []*models.TeamMember
	if err := tx.SelectContext(ctx, &memberList, query, params...); err != nil {
		return fmt.Errorf("failed team member query: %v", err)
	}
	teamMap := make(map[int32]*models.Team)
	for _, t := range teams {
		teamMap[t.TeamID] = t
	}
	for _, m := range memberList {
		teamMap[m.TeamID].Members = append(teamMap[m.TeamID].Members, m)
	}
	return nil
}

func teamListQuery(filterArgs TeamFilter) (string, []interface{}) {
	var filters []string
	var params []interface{}
	if filterArgs.ContestID != 0 {
		filters = append(filters, db.SetParam("contest_id = $%d", &params, filterArgs.ContestID))
	}
	if filterArgs.AccountID != 0 {
		filters = append(filters, db.SetParam("account_id = $%d", &params, filterArgs.AccountID))
	}
	filter := ""
	if len(filters) > 0 {
		filter = "WHERE " + strings.Join(filters, " AND ")
	}
	return fmt.Sprintf(`SELECT team_id, contest_id, team_name, start_time FROM team LEFT JOIN team_member USING(team_id) %s`, filter), params
}
