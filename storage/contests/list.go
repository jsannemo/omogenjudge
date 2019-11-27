package contests

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jsannemo/omogenjudge/storage/problems"

	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

// A ContestFilter controls the filtering behaviour of ListContests.
type ListFilter struct {
	// Filter by hostname.
	HostName string
}

type ListArgs struct {
	WithProblems bool
}

// ListContests returns a list of contests.
func ListContests(ctx context.Context, args ListArgs, filter ListFilter) (ContestList, error) {
	conn := db.Conn()
	var clist ContestList
	query, params := listQuery(filter)
	if err := conn.SelectContext(ctx, &clist, query, params...); err != nil {
		return nil, fmt.Errorf("failed contest list query: %v", err)
	}
	for _, c := range clist {
		if args.WithProblems {
			if err := addProblems(ctx, c, conn); err != nil {
				return nil, err
			}
		}
	}
	return clist, nil
}

func addProblems(ctx context.Context, contest *models.Contest, db *sqlx.DB) error {
	if err := db.SelectContext(ctx, &contest.Problems, `SELECT contest_id, problem_id, label FROM contest_problem WHERE contest_id = $1 ORDER BY label`, contest.ContestID); err != nil {
		return err
	}
	var pids []int32
	for _, k := range contest.Problems {
		pids = append(pids, k.ProblemID)
	}
	probs, err := problems.List(ctx, problems.ListArgs{WithStatements: problems.StmtTitles, WithTests: problems.TestsGroups}, problems.ListFilter{ProblemID: pids})
	if err != nil {
		return err
	}
	probMap := probs.AsMap()
	for _, k := range contest.Problems {
		k.Problem = probMap[k.ProblemID]
	}
	return nil
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
