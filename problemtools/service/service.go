package service

import (
	"context"

	"google.golang.org/grpc"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/courses"
	"github.com/jsannemo/omogenjudge/problemtools/problems"
)

type toolServer struct {
}

func (s *toolServer) ParseProblem(ctx context.Context, req *toolspb.ParseProblemRequest) (*toolspb.ParseProblemResponse, error) {
	return problems.ParseProblem(req.ProblemPath)
}

func (s *toolServer) VerifyProblem(ctx context.Context, req *toolspb.VerifyProblemRequest) (*toolspb.VerifyProblemResponse, error) {
	return problems.VerifyProblem(req.ProblemToVerify)
}

func (s *toolServer) ParseCourse(ctx context.Context, req *toolspb.ParseCourseRequest) (*toolspb.ParseCourseResponse, error) {
	return courses.ParseCourse(req.CoursePath)
}

func Register(grpcServer *grpc.Server) {
	server := &toolServer{}
	toolspb.RegisterToolServiceServer(grpcServer, server)
}
