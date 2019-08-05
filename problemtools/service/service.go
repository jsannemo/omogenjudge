package service

import (
	"context"

	"google.golang.org/grpc"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/courses"
	"github.com/jsannemo/omogenjudge/problemtools/problems"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
	rclient "github.com/jsannemo/omogenjudge/runner/client"
)

type toolServer struct {
	runner runpb.RunServiceClient
}

func (s *toolServer) ParseProblem(ctx context.Context, req *toolspb.ParseProblemRequest) (*toolspb.ParseProblemResponse, error) {
	return problems.ParseProblem(req.ProblemPath)
}

func (s *toolServer) VerifyProblem(ctx context.Context, req *toolspb.VerifyProblemRequest) (*toolspb.VerifyProblemResponse, error) {
	return problems.VerifyProblem(ctx, req, s.runner)
}

func (s *toolServer) ParseCourse(ctx context.Context, req *toolspb.ParseCourseRequest) (*toolspb.ParseCourseResponse, error) {
	return courses.ParseCourse(req.CoursePath)
}

func Register(grpcServer *grpc.Server) {
	server := &toolServer{
		rclient.NewClient(),
	}
	toolspb.RegisterToolServiceServer(grpcServer, server)
}
