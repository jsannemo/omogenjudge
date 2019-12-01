package runners

import (
	"errors"
	"path/filepath"

	execpb "github.com/jsannemo/omogenjudge/sandbox/api"
)

// An ExitType describes why a program exited.
type ExitType int

const (
	// Exited means the program exited normally with an exit code.
	Exited ExitType = iota

	// Signaled means the program was killed by a signal.
	Signaled

	// TimedOut means the program was killed due to exceeding its time limit.
	TimedOut

	// MemoryExceeded means the program was killed due to exceeding its allocated memory.
	MemoryExceeded
)

// An ExecResult describes the result of a single execution.
type ExecResult struct {
	// How how the program exited.
	ExitType ExitType

	// The exit code. Only set if the program exited with a code.
	ExitCode int32

	// The termination singal. Only set if the program exited with a signal.
	Signal int32

	// The time the execution used.
	TimeUsageMs int32

	// The memory the execution used.
	MemoryUsageKb int32
}

// CrashedWith checks whether the program exited normally with the given code.
func (res ExecResult) CrashedWith(code int32) bool {
	return res.ExitType == Exited && res.ExitCode == code
}

// Crashed checks whether the program crashed.
func (res ExecResult) Crashed() bool {
	return (res.ExitType == Exited && res.ExitCode != 0) || res.ExitType == Signaled
}

// TimedOut checks whether the program exceeded its time limit or not.
func (res ExecResult) TimedOut() bool {
	return res.ExitType == TimedOut
}

// MemoryExceeded checks whether the program exceeded its memory limit or not.
func (res ExecResult) MemoryExceeded() bool {
	return res.ExitType == MemoryExceeded
}

func makeStreams(in, out, err string) *execpb.Streams {
	return &execpb.Streams{
		Mappings: []*execpb.Streams_Mapping{
			&execpb.Streams_Mapping{
				Mapping: &execpb.Streams_Mapping_File_{
					&execpb.Streams_Mapping_File{PathInsideContainer: in}},
				Write: false,
			},
			&execpb.Streams_Mapping{
				Mapping: &execpb.Streams_Mapping_File_{
					&execpb.Streams_Mapping_File{PathInsideContainer: out}},
				Write: true,
			},
			&execpb.Streams_Mapping{
				Mapping: &execpb.Streams_Mapping_File_{
					&execpb.Streams_Mapping_File{PathInsideContainer: err}},
				Write: true,
			},
		},
	}
}

type ExecArgs struct {
	Command          string
	Args             []string
	WorkingDirectory string
	InputPath        string
	OutputPath       string
	ErrorPath        string
	ExtraReadPaths   []string
	ExtraWritePaths  []string
	TimeLimitMs      int32
	MemoryLimitKb    int32
	ReuseContainer   bool
}

func Execute(exec execpb.ExecuteService_ExecuteClient, args *ExecArgs) (*ExecResult, error) {
	req := &execpb.ExecuteRequest{
		Execution: &execpb.Execution{
			Command: &execpb.Command{
				Command: args.Command,
				Flags:   args.Args,
			},
			Environment: &execpb.Environment{
				WorkingDirectory:   args.WorkingDirectory,
				StreamRedirections: makeStreams(args.InputPath, args.OutputPath, args.ErrorPath),
			},
			ResourceLimits: &execpb.ResourceAmounts{
				Amounts: []*execpb.ResourceAmount{
					&execpb.ResourceAmount{
						Type:   execpb.ResourceType_CPU_TIME,
						Amount: int64(args.TimeLimitMs)},
					&execpb.ResourceAmount{
						Type:   execpb.ResourceType_WALL_TIME,
						Amount: 2 * int64(args.TimeLimitMs)},
					&execpb.ResourceAmount{
						Type:   execpb.ResourceType_MEMORY,
						Amount: int64(args.MemoryLimitKb)},
					&execpb.ResourceAmount{
						Type:   execpb.ResourceType_PROCESSES,
						Amount: 10},
				},
			},
		},
		ContainerSpec: &execpb.ContainerSpec{
			MaxDiskKb: 1000 * 1000, // 1 GB
			Mounts: makeMounts(
				append(args.ExtraReadPaths,
					filepath.Dir(args.InputPath),
					args.WorkingDirectory),
				append(args.ExtraWritePaths,
					filepath.Dir(args.OutputPath),
					filepath.Dir(args.ErrorPath)),
			)},
	}
	if args.ReuseContainer {
		req.ContainerSpec = nil
	}
	err := exec.Send(req)
	if err != nil {
		return nil, err
	}
	res, err := exec.Recv()
	if err != nil {
		return nil, err
	}
	result, err := toResult(res.Termination)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func toResult(termination *execpb.Termination) (*ExecResult, error) {
	time := int32(0)
	memory := int32(0)
	for _, res := range termination.UsedResources.Amounts {
		switch res.Type {
		case execpb.ResourceType_CPU_TIME:
			time = int32(res.Amount)
		case execpb.ResourceType_MEMORY:
			memory = int32(res.Amount)
		}
	}
	switch termination.Termination.(type) {
	case *execpb.Termination_Signal_:
		return &ExecResult{
			ExitType:      Signaled,
			Signal:        termination.GetSignal().Signal,
			TimeUsageMs:   time,
			MemoryUsageKb: memory,
		}, nil
	case *execpb.Termination_Exit_:
		return &ExecResult{
			ExitType:      Exited,
			ExitCode:      termination.GetExit().Code,
			TimeUsageMs:   time,
			MemoryUsageKb: memory,
		}, nil
	case *execpb.Termination_ResourceExceeded:
		if termination.GetResourceExceeded() == execpb.ResourceType_CPU_TIME {
			return &ExecResult{
				ExitType:      TimedOut,
				TimeUsageMs:   time,
				MemoryUsageKb: memory,
			}, nil
		} else if termination.GetResourceExceeded() == execpb.ResourceType_WALL_TIME {
			return &ExecResult{
				ExitType:      TimedOut,
				TimeUsageMs:   time,
				MemoryUsageKb: memory,
			}, nil
		} else if termination.GetResourceExceeded() == execpb.ResourceType_MEMORY {
			return &ExecResult{
				ExitType:      MemoryExceeded,
				TimeUsageMs:   time,
				MemoryUsageKb: memory,
			}, nil
		} else {
			return nil, errors.New("unknown resource exceeded")
		}
	default:
		return nil, errors.New("unknown termination")
	}
}

func appendPaths(seen map[string]bool, paths []string, writeable bool, dirs *[]*execpb.DirectoryMount) {
	for _, path := range paths {
		if seen[path] {
			continue
		}
		seen[path] = true
		*dirs = append(*dirs, &execpb.DirectoryMount{
			PathInsideContainer:  path,
			PathOutsideContainer: path,
			Writable:             writeable,
		})
	}
}

func makeMounts(readPaths, writePaths []string) []*execpb.DirectoryMount {
	seen := make(map[string]bool)
	var dirs []*execpb.DirectoryMount
	appendPaths(seen, writePaths, true, &dirs)
	appendPaths(seen, readPaths, false, &dirs)
	dirs = append(dirs, &execpb.DirectoryMount{
		PathInsideContainer:  "/etc",
		PathOutsideContainer: "/var/lib/omogen/etc",
		Writable:             false,
	})
	return dirs
}
