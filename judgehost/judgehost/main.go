package main

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/google/logger"
	"github.com/jsannemo/omogenexec/eval"
	apipb "github.com/jsannemo/omogenhost/judgehost/api"
	"github.com/jsannemo/omogenhost/storage"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
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
	Database  dbConfig
	Judgehost hostConfig
}

type JudgehostServer struct {
}

func (j *JudgehostServer) Evaluate(_ context.Context, request *apipb.EvaluateRequest) (*apipb.EvaluateResponse, error) {
	runId := request.RunId
	logger.Infof("Received run %d", runId)
	err := evaluate(runId)
	return &apipb.EvaluateResponse{}, err
}

func main() {
	defer logger.Init("localjudge", true, false, ioutil.Discard).Close()
	eval.InitLanguages()
	data, err := ioutil.ReadFile("/etc/omogen/judgehost.toml")
	if err != nil {
		panic(err)
	}
	var conf config
	if _, err := toml.Decode(string(data), &conf); err != nil {
		panic(err)
	}
	connStr := fmt.Sprintf("postgres://omogenjudge:omogenjudge@%s:%d/omogenjudge", conf.Database.Server, conf.Database.Port)
	if err := storage.Init(connStr); err != nil {
		panic(err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Judgehost.Server, conf.Judgehost.Port))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	if err != nil {
		logger.Fatalf("failed to create server: %v", err)
	}

	judgehostServer := &JudgehostServer{}
	apipb.RegisterJudgehostServiceServer(grpcServer, judgehostServer)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("could not listen: %v", err)
	}
}
