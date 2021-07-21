package main

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/google/logger"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/jsannemo/omogenhost/storage"
	apipb "github.com/jsannemo/omogenhost/webapi/proto"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type dbConfig struct {
	Server string
	Port   int
}

type apiConfig struct {
	Server string
	Port   int
}

type config struct {
	Database dbConfig
	Webapi   apiConfig
}

type ApiServer struct {
}

func (a ApiServer) ViewProblem(ctx context.Context, request *apipb.ViewProblemRequest) (*apipb.ViewProblemResponse, error) {
	return &apipb.ViewProblemResponse{
		Statement: &apipb.ProblemStatement{
			Language: "en",
			Title:    "Hello World",
			Html:     "<p>test statement</p>",
			License:  "cc by-sa",
			Authors:  []string{"Johan Sannemo", "Simon Lindholm"},
		},
		Limits: &apipb.ProblemLimits{
			TimeLimitMs:   1000,
			MemoryLimitKb: 1000 * 1000,
		},
	}, nil
}

func main() {
	defer logger.Init("webapi", true, false, ioutil.Discard).Close()
	data, err := ioutil.ReadFile("/etc/omogenjudge/webapi.toml")
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
	grpcServer := grpc.NewServer()
	if err != nil {
		logger.Fatalf("failed to create server: %v", err)
	}
	wrappedGrpc := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			return origin == "http://localhost:8000"
		}))
	httpServer := &http.Server{
		Addr: fmt.Sprintf("%s:%d", conf.Webapi.Server, conf.Webapi.Port),
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			if req.Method == "OPTIONS" {
				resp.Header().Add("Access-Control-Allow-Origin", "*")
				resp.Header().Add("Access-Control-Allow-Headers", "*")
				resp.Header().Add("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
				resp.WriteHeader(http.StatusOK)
				return
			}
			if wrappedGrpc.IsGrpcWebRequest(req) {
				wrappedGrpc.ServeHTTP(resp, req)
				return
			}
			// Fall back to other servers.
			http.DefaultServeMux.ServeHTTP(resp, req)
		}),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	apiServer := &ApiServer{}
	apipb.RegisterProblemServiceServer(grpcServer, apiServer)
	log.Fatal(httpServer.ListenAndServe())
}
