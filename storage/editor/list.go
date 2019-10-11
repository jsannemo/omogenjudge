package editor

import (
	"context"
	"fmt"
	"strings"

	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

type ListArgs struct {
	WithContent bool
}

type ListFilter struct {
	UserId int32
	Name   string
}

func listQuery(args ListArgs, filter ListFilter) (string, []interface{}) {
	var filterSegs []string
	var params []interface{}
  // TODO: require user ID for filtering
	if filter.UserId != 0 {
		params = append(params, filter.UserId)
		filterSegs = append(filterSegs, fmt.Sprintf("account_id = $%d", len(params)))
	}
	if filter.Name != 0 {
		params = append(params, filter.Name)
		filterSegs = append(filterSegs, fmt.Sprintf("file_name = $%d", len(params)))
	}
	filterStr := ""
	if len(filterSegs) != 0 {
		filterStr = fmt.Sprintf("WHERE %s", strings.Join(filterSegs, " AND "))
	}

  columnSegs := []string{"editor_file_id",  "account_id",  "file_name"}
  if args.WithContent {
    columnSegs = append(columnSegs, "file_content")
  }
	return fmt.Sprintf("SELECT %s FROM editor_file %s ORDER BY file_name DESC", filterStr), params
}

func List(ctx context.Context, args ListArgs, filter ListFilter) models.SubmissionList {
	conn := db.Conn()
	query, params := listQuery(args, filter)
	var subs models.SubmissionList
	if err := conn.SelectContext(ctx, &subs, query, params...); err != nil {
		panic(err)
	}
	if args.WithFiles {
		for _, sub := range subs {
			if err := conn.SelectContext(ctx, &sub.Files, "SELECT * FROM submission_file WHERE submission_id = $1", sub.SubmissionId); err != nil {
				panic(err)
			}
		}
	}
	return subs
}
