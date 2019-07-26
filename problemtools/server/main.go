// Server implementing the ToolsService.
package main

import (
	"context"
	"flag"
	"net"

	"github.com/google/logger"
	"google.golang.org/grpc"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/parser"
	"github.com/jsannemo/omogenjudge/problemtools/verifier"
)

var (
	toolsAddress = flag.String("tools_listen_addr", "127.0.0.1:61812", "The tool server address to listen to")
)

type toolServer struct {
}

func (s *toolServer) ParseProblem(ctx context.Context, req *toolspb.ParseProblemRequest) (*toolspb.ParseProblemResponse, error) {
	return parser.ParseProblem(req.ProblemPath)
}

func (s *toolServer) VerifyProblem(ctx context.Context, req *toolspb.VerifyProblemRequest) (*toolspb.VerifyProblemResponse, error) {
	return verifier.VerifyProblem(req.ProblemToVerify)
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *toolsAddress)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	server := &toolServer{}
	if err != nil {
		logger.Fatalf("failed to create server: %v", err)
	}
	toolspb.RegisterToolServiceServer(grpcServer, server)
	grpcServer.Serve(lis)
}
