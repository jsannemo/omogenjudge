// The evaluator queue keeps track of all unjudged submissions, including new ones.
package queue

import (
  "context"
  "strconv"

  "github.com/google/logger"

  "github.com/jsannemo/omogenjudge/storage/db"
  "github.com/jsannemo/omogenjudge/storage/submissions"
)

func StartQueue(judge chan<- *submissions.Submission) error {
  listener := db.NewListener()
  err := listener.Listen("new_submission")
  if err != nil {
    return err
  }
  unjudged, err := submissions.ListSubmissions(context.TODO(), submissions.ListArgs{WithFiles: true}, submissions.ListFilter{OnlyUnjudged: true})
  if err != nil {
    return err
  }
  go func() {
    alreadyJudged := int32(0)
    for _, sub := range unjudged {
      logger.Infof("Unjudged: %v", sub)
      judge <- sub
      alreadyJudged = sub.SubmissionId
    }
    for {
        notification := <-listener.Notify
        submissionId, _ := strconv.Atoi(notification.Extra)
        subs, err := submissions.ListSubmissions(context.TODO(), submissions.ListArgs{WithFiles: true}, submissions.ListFilter{SubmissionId: submissionId})
        if err != nil {
          logger.Errorf("failed listing unjudged submission: %v", err)
        } else if len(subs) == 0 {
          logger.Errorf("requested %d but was not present in DB", submissionId)
        } else if len(subs) > 1 {
          logger.Errorf("requested %d but got %d submissions", submissionId, len(subs))
        } else {
          sub := subs[0]
          if sub.SubmissionId > alreadyJudged {
            judge <- sub
          }
        }
    }
  }()
  return nil
}
