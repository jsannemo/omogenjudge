package util

import "fmt"

type Reporter interface {
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Err(string, ...interface{})
	HasError() bool
	Infos() []string
	Warnings() []string
	Errors() []string
	Merge(Reporter)
}

type reporter struct {
	info     []string
	warnings []string
	errors   []string
}

func (r *reporter) HasError() bool {
	return len(r.errors) != 0
}

func (r *reporter) addTo(msgs *[]string, msg string, args ...interface{}) {
	if len(args) != 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	*msgs = append(*msgs, msg)
}

func (r *reporter) Info(msg string, args ...interface{}) {
	r.addTo(&r.info, msg, args...)
}

func (r *reporter) Warn(msg string, args ...interface{}) {
	r.addTo(&r.warnings, msg, args...)
}

func (r *reporter) Err(msg string, args ...interface{}) {
	r.addTo(&r.errors, msg, args...)
}

func (r *reporter) Infos() []string {
	return r.info
}

func (r *reporter) Warnings() []string {
	return r.warnings
}

func (r *reporter) Errors() []string {
	return r.errors
}

func (r *reporter) Merge(other Reporter) {
	r.info = append(r.info, other.Infos()...)
	r.warnings = append(r.warnings, other.Warnings()...)
	r.errors = append(r.errors, other.Errors()...)
}

func NewReporter() Reporter {
	return &reporter{}
}
