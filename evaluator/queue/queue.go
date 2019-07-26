// The evaluator queue keeps track of all unjudged submissions, including new ones.
package queue

import (
  "strconv"

  "github.com/google/logger"

  "github.com/jsannemo/omogenjudge/storage/db"
  "github.com/jsannemo/omogenjudge/storage/submissions"
)

func StartQueue(judge chan<- int32) error {
  listener := db.NewListener()
  err := listener.Listen("new_submission")
  if err != nil {
    return err
  }
  unjudgedIds, err := submissions.UnjudgedIds()
  logger.Infof("Unjudged IDs: %v", unjudgedIds)
  if err != nil {
    return err
  }
  go func() {
    alreadyJudged := int32(0)
    for _, id := range unjudgedIds {
      judge <- id
      alreadyJudged = id
    }
    for {
        notification := <-listener.Notify
        submissionId, _ := strconv.Atoi(notification.Extra)
        subId := int32(submissionId)
        if subId > alreadyJudged {
          judge <- int32(subId)
        }
    }
  }()
  return nil
}
