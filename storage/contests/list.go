package contests

import (
	"context"
	"fmt"

	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

// A ContestFilter controls the filtering behaviour of ListContests.
type ListFilter struct {
	// Filter by hostname.
	HostName string
}

// ListContests returns a list of contests.
func ListContests(ctx context.Context, filter ListFilter) (ContestList, error) {
	conn := db.Conn()
	var clist ContestList
	query, params := listQuery(filter)
	if err := conn.SelectContext(ctx, &clist, query, params...); err != nil {
		return nil, fmt.Errorf("failed contest list query: %v", err)
	}
	return clist, nil
}

func listQuery(filterArgs ListFilter) (string, []interface{}) {
	filter := ""
	var params []interface{}
	if filterArgs.HostName != "" {
		filter = db.SetParam("WHERE host_name = $%d", &params, filterArgs.HostName)
	}
	return fmt.Sprintf(`SELECT contest_id, short_name, host_name, start_time, (EXTRACT(EPOCH FROM duration) * 1000000000)::bigint "duration", title, hidden_scoreboard FROM contest %s`, filter), params
}

// An ContestList is a slice of contests.
type ContestList []*models.Contest

func (lists ContestList) Latest() *models.Contest {
	which := lists[0]
	for _, c := range lists {
		if c.StartTime.Valid && c.StartTime.Time.After(which.StartTime.Time) {
			which = c
		}
	}
	return which
}

