package submissions

import (
	"context"
	"fmt"
	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

// A ListArgs controls what data to include in submissions.
type ListArgs struct {
	// Whether to include submission files in the query.
	WithFiles bool
}

// A ListFilter controls what submissions to search for.  Only one of SubmissionID and UserID may be set.
type ListFilter struct {
	SubmissionID []int32
	UserID       int32
}

// A SubmissionList is a slice of Submissions.
type SubmissionList []*models.Submission

// ListSubmissions searches for a list of submissions.
func ListSubmissions(ctx context.Context, args ListArgs, filterArgs ListFilter) (SubmissionList, error) {
	conn := db.Conn()
	filter := ""
	var params []interface{}
	if len(filterArgs.SubmissionID) != 0 {
		filter = db.SetInParamInt(`WHERE submission_id IN (%s)`, &params, filterArgs.SubmissionID)
	} else if filterArgs.UserID != 0 {
		filter = db.SetParam(`WHERE account_id = %s`, &params, filterArgs.UserID)
	}
	query := fmt.Sprintf(`SELECT * FROM submission %s ORDER BY submission_id DESC`, filter)
	var subs SubmissionList
	if err := conn.SelectContext(ctx, &subs, query, params...); err != nil {
		return nil, fmt.Errorf("could not read submissions: %v", err)
	}
	if args.WithFiles {
		for _, sub := range subs {
			if err := conn.SelectContext(ctx, &sub.Files, `SELECT * FROM submission_file WHERE submission_id = $1`, sub.SubmissionID); err != nil {
				return nil, fmt.Errorf("could not read submission files: %v", err)
			}
		}
	}
	return subs, nil
}
