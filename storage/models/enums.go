package models

import (
	"fmt"
)

type Status string

const (
	StatusNew               Status = "new"
	StatusCompiling         Status = "compiling"
	StatusCompilationFailed Status = "compilation_failed"
	StatusRunning           Status = "running"
	StatusSuccessful        Status = "successful"
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

func (s Status) Accepted() bool {
	return false
}

func (s Status) Rejected() bool {
	return s == StatusCompilationFailed
}

func (s Status) Waiting() bool {
	return s == StatusNew || s == StatusCompiling || s == StatusRunning
}

type Verdict string

const (
	VerdictUnjudged          Verdict = "UNJUDGED"
	VerdictAccepted          Verdict = "ACCEPTED"
	VerdictTimeLimitExceeded Verdict = "TIME_LIMIT_EXCEEDED"
	VerdictRunTimeError      Verdict = "RUN_TIME_ERROR"
	VerdictWrongAnswer       Verdict = "WRONG_ANSWER"
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

func (v Verdict) Accepted() bool {
	return v == VerdictAccepted
}

func (v Verdict) Rejected() bool {
	return v != VerdictAccepted
}

func (v Verdict) Waiting() bool {
	return false
}

type License string

const (
	LicenseCcBySa3      License = "CC_BY_SA_3"
	LicensePermission   License = "BY_PERMISSION"
	LicensePublicDomain License = "PUBLIC_DOMAIN"
	LicensePrivate      License = "PRIVATE"
)

func (l License) String() string {
	switch l {
	case LicenseCcBySa3:
		return "CC BY-SA 3.0"
	case LicensePublicDomain:
		return "Fri användning"
	case LicensePermission:
		return "Används med tillåtelse"
	case LicensePrivate:
		return "Endast för privat användning"
	}
	panic(fmt.Errorf("Unknown license: %v", l))
}
