package main

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/google/logger"
	apipb "github.com/jsannemo/omogenhost/judgehost/api"
	"github.com/jsannemo/omogenhost/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"strconv"
	"time"
)

type dbConfig struct {
	Server string
	Port   int
}

type hostConfig struct {
	Server string
	Port   int
}

type config struct {
	Database   dbConfig
	Judgehosts hostConfig
}

func NewClient(address string) apipb.JudgehostServiceClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		logger.Fatalf("fail to dial: %v", err)
	}
	return apipb.NewJudgehostServiceClient(conn)
}

func main() {
	defer logger.Init("judgequeue", true, false, ioutil.Discard).Close()
	data, err := ioutil.ReadFile("/etc/omogen/queue.toml")
	if err != nil {
		panic(err)
	}
	var conf config
	if _, err := toml.Decode(string(data), &conf); err != nil {
		panic(err)
	}

	hostClient := NewClient(fmt.Sprintf("%s:%d", conf.Judgehosts.Server, conf.Judgehosts.Port))

	connStr := fmt.Sprintf("postgres://omogenjudge:omogenjudge@%s:%d/omogenjudge", conf.Database.Server, conf.Database.Port)
	if err := storage.Init(connStr); err != nil {
		panic(err)
	}
	logger.Info("Starting judging queue")
	listener := storage.NewListener(connStr)
	if err := listener.Listen("new_run"); err != nil {
		logger.Fatalf("Failed starting database listener: %v", err)
	}
	logger.Infoln("Started database listener")

	var unjudgedRuns []storage.SubmissionRun
	res := storage.GormDB.Select("submission_run_id").Order("submission_run_id asc").Where("status = ?", storage.StatusQueued).Find(&unjudgedRuns)
	if res.Error != nil {
		logger.Fatalf("Failed loading run backlog: %v", res.Error)
	}
	logger.Infof("Had backlog of %d submissions", len(unjudgedRuns))
	judgeChan := make(chan int64, len(unjudgedRuns)+100)
	go func() {
		var alreadyJudged int64 = 0
		for _, sub := range unjudgedRuns {
			id := sub.SubmissionRunId
			judgeChan <- id
			alreadyJudged = id
		}
		for {
			notification := <-listener.Notify
			runId, _ := strconv.ParseInt(notification.Extra, 10, 64)
			// We may have read some of the newly delivered submissions in our list call,
			// so we need to filter out any earlier submissions.
			if runId > alreadyJudged {
				judgeChan <- runId
			}
		}
	}()
	for sub := range judgeChan {
		logger.Infof("Sending submission %d for judging", sub)
		// TODO: give context a deadline to prevent stuck judge hosts...
		ctx := context.Background()
		req := &apipb.EvaluateRequest{RunId: sub}
		for {
			_, err := hostClient.Evaluate(ctx, req)
			errcode := status.Code(err)
			if errcode == codes.Unavailable {
				logger.Infof("Judge host unavailable; retrying in 10s...")
				time.Sleep(time.Second * 10)
				continue
			}

			// TODO: retry failed judging 1 more time
			if err != nil {
				logger.Fatalf("Failed judging %d: %v", sub, err)
				if res := storage.GormDB.Model(
					&storage.SubmissionRun{SubmissionRunId: sub},
				).Update("Status", storage.StatusJudgeError); res.Error != nil {
					logger.Warningf("failed marking run as compiling: %v", res.Error)
				}
			}
			break
		}
		logger.Infof("Done judging run %d", sub)
	}
}
