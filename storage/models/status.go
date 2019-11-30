package models

import (
	"fmt"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

// A Status is the execution status of a submission run.
type Status string

const (
	// The status of a submission that has not yet been processed.
	StatusNew Status = "new"
	// The status of a submission that is currently compiling.
	StatusCompiling Status = "compiling"
	// The status of a submission that failed compiling.
	StatusCompilationFailed Status = "compilation_failed"
	// The status of a submission that is currently being judged on test data.
	StatusRunning Status = "running"
	// The status of a submission that is been successfully processed.
	StatusSuccessful Status = "successful"
)

// String returns a user-friendly name of a status.
func (s Status) String() string {
	switch s {
	case StatusNew:
		return "I kön"
	case StatusCompiling:
		return "Kompilerar"
	case StatusCompilationFailed:
		return "Kompileringsfel"
	case StatusRunning:
		return "Kör"
	case StatusSuccessful:
		return "Färdig"
	}
	panic(fmt.Errorf("unknown status: %v", s))
}

// Waiting returns whether the status represents a submission run that is still being
func (s Status) Waiting() bool {
	return s == StatusNew || s == StatusCompiling || s == StatusRunning
}

// A Verdict is the result of a submission on a test case, test group or the entire test data.
type Verdict string

const (
	VerdictUnjudged          Verdict = "unjudged"
	VerdictAccepted          Verdict = "accepted"
	VerdictTimeLimitExceeded Verdict = "time_limit_exceeded"
	VerdictRunTimeError      Verdict = "run_time_error"
	VerdictWrongAnswer       Verdict = "wrong_answer"
)

// VerdictFromRunVerdict converts a verdict returned from the RunService to a Verdict.
func VerdictFromRunVerdict(verdict runpb.Verdict) Verdict {
	switch verdict {
	case runpb.Verdict_ACCEPTED:
		return VerdictAccepted
	case runpb.Verdict_RUN_TIME_ERROR:
		return VerdictRunTimeError
	case runpb.Verdict_TIME_LIMIT_EXCEEDED:
		return VerdictTimeLimitExceeded
	case runpb.Verdict_WRONG_ANSWER:
		return VerdictWrongAnswer
	}
	panic(fmt.Errorf("invalid verdict: %v", verdict))
}

// String returns a user-friendly name of a verdict.
func (v Verdict) String() string {
	switch v {
	case VerdictUnjudged:
		return "Inte bedömd"
	case VerdictAccepted:
		return "Godkänd"
	case VerdictWrongAnswer:
		return "Fel svar"
	case VerdictRunTimeError:
		return "Körningsfel"
	case VerdictTimeLimitExceeded:
		return "Tidsgräns överskriden"
	}
	panic(fmt.Errorf("unknown verdict: %s", string(v)))
}

func (v Verdict) Accepted() bool {
	return v == VerdictAccepted
}

func (v Verdict) Waiting() bool {
	return v == VerdictUnjudged
}
