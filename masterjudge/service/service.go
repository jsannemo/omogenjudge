package service

import (
	"context"

	"google.golang.org/grpc"

	masterpb "github.com/jsannemo/omogenjudge/masterjudge/api"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
	rclient "github.com/jsannemo/omogenjudge/runner/client"
)

type masterServer struct {
	run runpb.RunServiceClient
}

func (s *masterServer) GetLanguages(ctx context.Context, _ *masterpb.GetLanguagesRequest) (*masterpb.GetLanguagesResponse, error) {
	langs, err := s.run.GetLanguages(ctx, &runpb.GetLanguagesRequest{})
	if err != nil {
		return nil, err
	}
	return &masterpb.GetLanguagesResponse{InstalledLanguages: langs.InstalledLanguages}, nil
}

func newServer() *masterServer {
	return &masterServer{
		run: rclient.NewClient(),
	}
}

func Register(grpcServer *grpc.Server) {
	server := newServer()
	masterpb.RegisterMasterServiceServer(grpcServer, server)
}
