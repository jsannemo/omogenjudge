package models

import (
  "fmt"
)

type Status string

const (
  StatusNew Status = "new"
  StatusCompiling Status = "compiling"
  StatusCompilationFailed Status = "compilation_failed"
  StatusRunning Status = "running"
  StatusSuccessful Status = "successful"
)

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
  panic(fmt.Errorf("Unknown status: %v", s))
}

type Verdict string

const (
  VerdictUnjudged Verdict = "UNJUDGED"
  VerdictAccepted Verdict = "ACCEPTED"
  VerdictTimeLimitExceeded Verdict = "TIME_LIMIT_EXCEEDED"
  VerdictRunTimeError Verdict = "RUN_TIME_ERROR"
  VerdictWrongAnswer Verdict = "WRONG_ANSWER"
)

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
  panic(fmt.Errorf("Unknown verdict: %v", v))
}
