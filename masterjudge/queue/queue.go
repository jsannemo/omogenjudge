// The evaluator queue keeps track of all unjudged submissions, including new ones.
package queue

import (
	"context"
	"strconv"

	"github.com/google/logger"

	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/submissions"
)

func StartQueue(ctx context.Context, judge chan<- *models.Submission) error {
	logger.Infoln("Starting submission listener")
	listener := db.NewListener()
	err := listener.Listen("new_submission")
	if err != nil {
		return err
	}
	logger.Infoln("Started listener")
	unjudged := submissions.List(ctx, submissions.ListArgs{WithFiles: true}, submissions.ListFilter{OnlyUnjudged: true})
	logger.Infof("Had backlog of %d submissions", len(unjudged))
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
			subs := submissions.List(ctx, submissions.ListArgs{WithFiles: true}, submissions.ListFilter{SubmissionId: int32(submissionId)})
			if len(subs) == 0 {
				logger.Errorf("requested %d but was not present in DB", submissionId)
			} else if len(subs) > 1 {
				logger.Errorf("requested %d but got %d submissions", submissionId, len(subs))
			} else {
				sub := subs[0]
				// We may have read some of the newly delivered submissions in our list call,
				// so we need to filter out any earlier submissions.
				if sub.SubmissionId > alreadyJudged {
					judge <- sub
				}
			}
		}
	}()
	return nil
}
