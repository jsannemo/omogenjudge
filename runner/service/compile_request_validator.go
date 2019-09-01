package service

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

func ValidateCompileRequest(req *runpb.CompileRequest) *status.Status {
	if req.Program == nil {
		return status.Newf(codes.InvalidArgument, "Missing program to compile")
	}
	if req.Program.LanguageId == "" {
		return status.Newf(codes.InvalidArgument, "Program has no language ID")
	}
	if len(req.Program.Sources) == 0 {
		return status.Newf(codes.InvalidArgument, "Program has no source files")
	}
	if req.OutputPath == "" {
		return status.Newf(codes.InvalidArgument, "Missing output path")
	}
	return nil
}
