package submissions

import (
  "context"
  "fmt"
  "strings"

  "github.com/jsannemo/omogenjudge/storage/db"
  "github.com/jsannemo/omogenjudge/storage/models"
)

type ListArgs struct {
  WithFiles bool
}

type ListFilter struct {
  OnlyUnjudged bool
  SubmissionId int32
  UserId int32
}

func listQuery(args ListArgs, filter ListFilter) (string, []interface{}) {
  var filterSegs []string
  var params []interface{}
  if filter.OnlyUnjudged {
    filterSegs = append(filterSegs, "status != 'successful' AND status != 'compilation_failed'")
  }
  if filter.SubmissionId != 0 {
    filterSegs = append(filterSegs, "submission_id = $1")
    params = append(params, filter.SubmissionId)
  }
  if filter.UserId != 0 {
    filterSegs = append(filterSegs, "account_id = $1")
    params = append(params, filter.UserId)
  }
  filterStr := ""
  if len(filterSegs) != 0 {
    filterStr = fmt.Sprintf("WHERE %s", strings.Join(filterSegs, " AND "))
  }
  return fmt.Sprintf("SELECT * FROM submission %s ORDER BY submission_id DESC", filterStr), params
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
