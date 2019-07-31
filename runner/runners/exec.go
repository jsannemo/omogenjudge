package runners

import (
  "errors"
  "os"
  "path/filepath"

  "github.com/google/logger"

  execpb "github.com/jsannemo/omogenjudge/sandbox/api"
)

type ExecArgs struct {
  Command string
  Args []string
  WorkingDirectory string
  InputPath string
  OutputPath string
  ErrorPath string
  ExtraReadPaths []string
  ExtraWritePaths []string
  TimeLimitMs int64
  MemoryLimitKb int64
  ReuseContainer bool
}

type ExitReason int

const (
  Exited ExitReason = iota
  Signaled
  TimedOut
)

type ExecResult struct {
  ExitReason ExitReason
  ExitCode int
  Signal int
  TimeUsageMs int
  MemoryUsageKb int
}


func ensureFolder(path string) {
  os.MkdirAll(path, 0755)
}

func makeStreams(in, out, err string) *execpb.Streams {
  ensureFolder(filepath.Dir(out))
  ensureFolder(filepath.Dir(err))
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

func Execute(exec execpb.ExecuteService_ExecuteClient, args *ExecArgs) (*ExecResult, error) {
  req := &execpb.ExecuteRequest{
    Execution: &execpb.Execution{
      Command: &execpb.Command{
        Command: args.Command,
        Flags: args.Args,
      },
      Environment: &execpb.Environment{
        WorkingDirectory: args.WorkingDirectory,
        StreamRedirections: makeStreams(args.InputPath, args.OutputPath, args.ErrorPath),
      },
      ResourceLimits: &execpb.ResourceAmounts{
        Amounts: []*execpb.ResourceAmount{
          &execpb.ResourceAmount{
            Type: execpb.ResourceType_CPU_TIME,
            Amount: args.TimeLimitMs},
          &execpb.ResourceAmount{
            Type: execpb.ResourceType_WALL_TIME,
            Amount: 2 * args.TimeLimitMs},
          &execpb.ResourceAmount{
            Type: execpb.ResourceType_MEMORY,
            Amount: args.MemoryLimitKb},
          &execpb.ResourceAmount{
            Type: execpb.ResourceType_PROCESSES,
            Amount: 10},
        },
      },
    },
    ContainerSpec: &execpb.ContainerSpec{
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
  logger.Infof("Sending Execute: %v", req)
  err := exec.Send(req)
  if err != nil {
    return nil, err
  }
  res, err := exec.Recv()
  logger.Infof("Received Execute: %v", res)
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
  switch termination.Termination.(type) {
  case *execpb.Termination_Signal_:
    return &ExecResult{
      ExitReason: Signaled,
    }, nil
  case *execpb.Termination_Exit_:
    return &ExecResult{
      ExitReason: Exited,
      ExitCode: int(termination.GetExit().Code),
    }, nil
  case *execpb.Termination_ResourceExceeded:
    if termination.GetResourceExceeded() == execpb.ResourceType_CPU_TIME {
      return &ExecResult{
        ExitReason: TimedOut,
      }, nil
    } else if termination.GetResourceExceeded() == execpb.ResourceType_WALL_TIME {
      return &ExecResult{
        ExitReason: TimedOut,
      }, nil
    } else if termination.GetResourceExceeded() == execpb.ResourceType_MEMORY {
      return &ExecResult{
        ExitReason: Signaled,
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
      PathInsideContainer: path,
      PathOutsideContainer: path,
      Writable: writeable,
    })
  }
}

func makeMounts(readPaths, writePaths []string) []*execpb.DirectoryMount {
  seen := make(map[string]bool)
  var dirs []*execpb.DirectoryMount
  appendPaths(seen, writePaths, true, &dirs)
  appendPaths(seen, readPaths, false, &dirs)
  dirs = append(dirs, &execpb.DirectoryMount{
    PathInsideContainer: "/etc",
    PathOutsideContainer: "/var/lib/omogen/etc",
    Writable: false,
  })
  return dirs
}

