package main

import (
	"flag"
	"net"
  "os"

	"github.com/google/logger"
	"google.golang.org/grpc"

	toolservice "github.com/jsannemo/omogenjudge/problemtools/service"
	runnerservice "github.com/jsannemo/omogenjudge/runner/service"
	fileservice "github.com/jsannemo/omogenjudge/filehandler/service"
)

var (
	listenAddr = flag.String("localjudge_listen_addr", "127.0.0.1:61811", "The local judge server address to listen to")
)

func main() {
	flag.Parse()
  defer logger.Init("localjudge", false, true, os.Stderr).Close()
	lis, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	if err != nil {
		logger.Fatalf("failed to create server: %v", err)
	}
  toolservice.Register(grpcServer)
  runnerservice.Register(grpcServer)
  fileservice.Register(grpcServer)
  logger.Infof("serving on %v", *listenAddr)
	grpcServer.Serve(lis)
}
