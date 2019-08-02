package service

import (
	"context"

	"google.golang.org/grpc"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/parser"
	"github.com/jsannemo/omogenjudge/problemtools/verifier"
)

type toolServer struct {
}

func (s *toolServer) ParseProblem(ctx context.Context, req *toolspb.ParseProblemRequest) (*toolspb.ParseProblemResponse, error) {
	return parser.ParseProblem(req.ProblemPath)
}

func (s *toolServer) VerifyProblem(ctx context.Context, req *toolspb.VerifyProblemRequest) (*toolspb.VerifyProblemResponse, error) {
	return verifier.VerifyProblem(req.ProblemToVerify)
}

func Register(grpcServer *grpc.Server) {
	server := &toolServer{}
	toolspb.RegisterToolServiceServer(grpcServer, server)
}

