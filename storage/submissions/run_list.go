package submissions

import (
	"context"
	"fmt"
	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
	"strings"
)

type RunListArgs struct {
}

// A RunListFilter controls what runs to search for.
type RunListFilter struct {
	// Whether only unprocessed runs should be included.
	OnlyUnjudged bool
	RunID []int32
}

// ListRuns searches for a list of runs.
func ListRuns(ctx context.Context, args RunListArgs, filterArgs RunListFilter) ([]*models.SubmissionRun, error) {
	conn := db.Conn()
	var filters []string
	var params []interface{}
	if len(filterArgs.RunID) != 0 {
		filters = append(filters, db.SetInParamInt(`submission_run_id IN (%s)`, &params, filterArgs.RunID))
	}
	if filterArgs.OnlyUnjudged {
		filters = append(filters, `status = 'new'`)
	}
	filter := strings.Join(filters, ", ")
	if filter != "" {
		filter = "WHERE = " + filter
	}
	query := fmt.Sprintf(`SELECT * FROM submission_run %s ORDER BY submission_run_id ASC`, filter)
	var runs []*models.SubmissionRun
	if err := conn.SelectContext(ctx, &runs, query, params...); err != nil {
		return nil, fmt.Errorf("could not read runs: %v", err)
	}
	return runs, nil
}
