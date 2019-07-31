package submissions

import (
  "context"
  "fmt"
  "strings"

  "github.com/jsannemo/omogenjudge/storage/db"
)

type Field string

const (
  FieldVerdict Field = "verdict"
  FieldStatus Field = "status"
)

type UpdateArgs struct {
  Fields []Field
}

func updateQuery(sub *Submission, args UpdateArgs) (string, []interface{}) {
  params := []interface{}{sub.SubmissionId}
  var updates []string
  for _, field := range args.Fields {
    switch field {
    case FieldVerdict:
      updates = append(updates, fmt.Sprintf("verdict = $%d", len(params) + 1))
      params = append(params, sub.Verdict)
    case FieldStatus:
      updates = append(updates, fmt.Sprintf("status = $%d", len(params) + 1))
      params = append(params, sub.Status)
    default:
    }
  }
  return fmt.Sprintf("UPDATE submission SET %s WHERE submission_id = $1", strings.Join(updates, ", ")), params
}

func Update(ctx context.Context, sub *Submission, args UpdateArgs) error {
  if len(args.Fields) == 0 {
    return nil
  }
  conn := db.GetPool()
  query, params := updateQuery(sub, args)
  _, err := conn.Exec(query, params...)
  if err != nil {
    return err
  }
  return nil
}
