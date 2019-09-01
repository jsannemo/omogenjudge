package runners

import (
	runpb "github.com/jsannemo/omogenjudge/runner/api"
	execpb "github.com/jsannemo/omogenjudge/sandbox/api"
)

// A RunFunc takes a compiled program and an executor and returns an executable program.
type RunFunc func(*runpb.CompiledProgram, execpb.ExecuteService_ExecuteClient) (Program, error)

// A ProgramArgs represents the arguments for a particular execution of a Program.
type ProgramArgs struct {
	// The file path that should be mapped to stdin in the execution.
	InputPath string

	// The file path that should be mapped to stdout in the execution.
	OutputPath string

	// The file path that should be mapped to stderr in the execution.
	ErrorPath string

	// The time limit to enforce on the execution.
	TimeLimitMs int64

	// The memory limit to enforce on the execution.
	MemoryLimitKb int64

	// Potential extra arguments that should be provided to the program.
	ExtraArgs []string
}

// A Program is an abstraction around a program that can be run.
type Program interface {
	// SetArgs sets the arguments for the next execution.
	SetArgs(*ProgramArgs)

	// Execute executes the program with the given arguments.
	Execute() (*ExecResult, error)
}

// A CommandArgs represents how a program should be executed; the command used, what arguments are necessery and so on.
type CommandArgs struct {
	// The executable file that should be run.
	Command string

	// The arguments that should be provided to the program.
	Args []string

	// The working directory the command should be executed.
	WorkingDirectory string
}

// An ArgFunc takes provides the CommandArgs necessary to run a compiled program.
type ArgFunc func(prog *runpb.CompiledProgram) *CommandArgs

// An argProgram is a Program implementation that uses a CommandArgs to execute programs.
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
	if err != nil {
		return nil, err
	}
	a.first = false
	// TODO: this should ch
	if res.TimedOut() {
		a.first = true
	}
	return res, err
}

// CommandProgram returns a RunFunc that creates programs based on the given ArgFunc.
func CommandProgram(argFunc ArgFunc) RunFunc {
	return func(prog *runpb.CompiledProgram, client execpb.ExecuteService_ExecuteClient) (Program, error) {
		return &argProgram{
			args:   argFunc(prog),
			client: client,
		}, nil
	}
}
