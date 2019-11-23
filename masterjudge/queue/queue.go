// Package queue provides an evaluator queue that keeps track of all unjudged submission runs, including new ones.
package queue

import (
	"context"
	"strconv"

	"github.com/google/logger"

	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/submissions"
)

func StartQueue(ctx context.Context, judge chan<- *models.SubmissionRun) error {
	logger.Infoln("Starting run listener")
	listener := db.NewListener()
	err := listener.Listen("new_run")
	if err != nil {
		return err
	}
	logger.Infoln("Started listener")
	unjudged, err := submissions.ListRuns(ctx, submissions.RunListArgs{}, submissions.RunListFilter{OnlyUnjudged: true})
	if err != nil {
		return err
	}
	logger.Infof("Had backlog of %d submissions", len(unjudged))
	go func() {
		alreadyJudged := int32(0)
		for _, sub := range unjudged {
			logger.Infof("Unjudged: %v", sub)
			judge <- sub
			alreadyJudged = sub.SubmissionRunID
		}
		for {
			notification := <-listener.Notify
			submissionId, _ := strconv.Atoi(notification.Extra)
			unjudged, err := submissions.ListRuns(ctx, submissions.RunListArgs{}, submissions.RunListFilter{RunID: []int32{int32(submissionId)}})
			if err != nil {
				panic(err)
			}
			if len(unjudged) == 0 {
				logger.Errorf("requested %d but was not present in DB", submissionId)
			} else if len(unjudged) > 1 {
				logger.Errorf("requested %d but got %d submissions", submissionId, len(unjudged))
			} else {
				sub := unjudged[0]
				// We may have read some of the newly delivered submissions in our list call,
				// so we need to filter out any earlier submissions.
				if sub.SubmissionRunID > alreadyJudged {
					judge <- sub
				}
			}
		}
	}()
	return nil
}
