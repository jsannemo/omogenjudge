package util

import "fmt"

type Reporter interface {
	Err(string, ...interface{})
	Warn(string, ...interface{})
	HasError() bool
	Errors() []string
	Warnings() []string
	AddFailures(*[]string, *[]string)
}

type reporter struct {
	warnings []string
	errors   []string
}

func (r *reporter) HasError() bool {
	return len(r.errors) != 0
}

func (r *reporter) Err(msg string, args ...interface{}) {
	if len(args) != 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	r.errors = append(r.errors, msg)
}

func (r *reporter) Warn(msg string, args ...interface{}) {
	if len(args) != 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	r.warnings = append(r.warnings, fmt.Sprintf(msg, args))
}

func (r *reporter) Errors() []string {
	return r.errors
}

func (r *reporter) Warnings() []string {
	return r.warnings
}

func NewReporter() Reporter {
	return &reporter{}
}

func (reporter *reporter) AddFailures(errors *[]string, warnings *[]string) {
	*errors = append(*errors, reporter.Errors()...)
	*warnings = append(*warnings, reporter.Warnings()...)
}
