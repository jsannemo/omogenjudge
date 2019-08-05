package runners

import (
	runpb "github.com/jsannemo/omogenjudge/runner/api"
	execpb "github.com/jsannemo/omogenjudge/sandbox/api"
)

type RunFunc func(*runpb.CompiledProgram, execpb.ExecuteService_ExecuteClient) (Program, error)

type ProgramArgs struct {
	InputPath     string
	OutputPath    string
	ErrorPath     string
	TimeLimitMs   int64
	MemoryLimitKb int64
	ExtraArgs     []string
}

type Program interface {
	SetArgs(*ProgramArgs)
	Execute() (*ExecResult, error)
}

type CommandArgs struct {
	Command          string
	Args             []string
	WorkingDirectory string
}

type ArgFunc func(prog *runpb.CompiledProgram) *CommandArgs

type argProgram struct {
	args        *CommandArgs
	programArgs *ProgramArgs
	first       bool
	client      execpb.ExecuteService_ExecuteClient
}

func (a *argProgram) SetArgs(args *ProgramArgs) {
	a.programArgs = args
	a.first = true
}

func (a *argProgram) Execute() (*ExecResult, error) {
	res, err := Execute(a.client,
		&ExecArgs{
			Command:          a.args.Command,
			Args:             append(a.args.Args, a.programArgs.ExtraArgs...),
			WorkingDirectory: a.args.WorkingDirectory,
			InputPath:        a.programArgs.InputPath,
			OutputPath:       a.programArgs.OutputPath,
			ErrorPath:        a.programArgs.ErrorPath,
			TimeLimitMs:      a.programArgs.TimeLimitMs,
			MemoryLimitKb:    a.programArgs.MemoryLimitKb,
			ReuseContainer:   !a.first,
		})
	a.first = false
	return res, err
}

func CommandProgram(argFunc ArgFunc) RunFunc {
	return func(prog *runpb.CompiledProgram, client execpb.ExecuteService_ExecuteClient) (Program, error) {
		return &argProgram{
			args:   argFunc(prog),
			client: client,
		}, nil
	}
}
