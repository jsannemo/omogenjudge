package service

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

func validateProgram(program *runpb.CompiledProgram) error {
	if program.ProgramRoot == "" {
		return status.Errorf(codes.InvalidArgument, "No submission program path")
	}
	if len(program.CompiledPaths) == 0 {
		return status.Errorf(codes.InvalidArgument, "Empty submission program")
	}
	if program.LanguageId == "" {
		return status.Errorf(codes.InvalidArgument, "Submission program has no language")
	}
	return nil
}

func ValidateEvaluateRequest(req *runpb.EvaluateRequest) error {
	if req.SubmissionId <= 0 {
		return status.Errorf(codes.InvalidArgument, "Non-positive submission ID")
	}
	if req.Program == nil {
		return status.Errorf(codes.InvalidArgument, "No submission program")
	}
	if err := validateProgram(req.Program); err != nil {
		return err
	}
	if len(req.Cases) == 0 {
		return status.Errorf(codes.InvalidArgument, "Evaluation had no test cases")
	}
	if req.TimeLimitMs <= 0 {
		return status.Errorf(codes.InvalidArgument, "Non-positive time limit")
	}
	if req.MemLimitKb <= 0 {
		return status.Errorf(codes.InvalidArgument, "Non-positive memory limit")
	}
	if req.Validator != nil {
		if err := validateProgram(req.Validator.Program); err != nil {
			return err
		}
	}
	return nil
}
