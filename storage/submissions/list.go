package submissions

import (
  "context"
  "fmt"
  "strings"

  "database/sql"

  "github.com/jsannemo/omogenjudge/storage/db"
)

type ListArgs struct {
  WithFiles bool
}

type ListFilter struct {
  OnlyUnjudged bool
  SubmissionId int
}

func listQuery(args ListArgs, filter ListFilter) (string, []interface{}) {
  var filterSegs []string
  var params []interface{}
  if filter.OnlyUnjudged {
    filterSegs = append(filterSegs, "status != 'successful'")
  }
  if filter.SubmissionId != 0 {
    filterSegs = append(filterSegs, "submission_id = $1")
    params = append(params, filter.SubmissionId)
  }
  filterStr := ""
  if len(filterSegs) != 0 {
    filterStr = fmt.Sprintf("WHERE %s", strings.Join(filterSegs, " AND "))
  }
  return fmt.Sprintf("SELECT problem_id, account_id, submission_id FROM submission %s", filterStr), params
}

func ListSubmissions(ctx context.Context, args ListArgs, filter ListFilter) ([]*Submission, error) {
  conn := db.GetPool()
  query, params := listQuery(args, filter)
  rows, err := conn.Query(query, params...)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  var submissions []*Submission
  for rows.Next() {
    submission := &Submission{}
    if err := rows.Scan(&submission.ProblemId, &submission.AccountId, &submission.SubmissionId); err != nil {
      return nil, err
    }
    if args.WithFiles {
      readSubmissionFiles(ctx, conn, submission)
    }
    submissions = append(submissions, submission)
  }
	if err := rows.Err(); err != nil {
    return nil, err
	}
  return submissions, nil
}

func readSubmissionFiles(ctx context.Context, conn *sql.DB, sub *Submission) error {
  rows, err := conn.Query("SELECT file_path, file_contents FROM submission_file WHERE submission_id = $1", sub.SubmissionId)
  if err != nil {
    return err
  }
  defer rows.Close()

  for rows.Next() {
    subFile := &SubmissionFile{}
    if err := rows.Scan(&subFile.Path, &subFile.Contents); err != nil {
      return err
    }
    sub.Files = append(sub.Files, subFile)
  }
  return nil
}

