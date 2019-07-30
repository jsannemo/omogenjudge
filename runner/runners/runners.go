package runners

import (
  "os"
  "path/filepath"

  "github.com/google/logger"

  execpb "github.com/jsannemo/omogenjudge/sandbox/api"
  runpb "github.com/jsannemo/omogenjudge/runner/api"
)

func makeMounts(readPaths, writePaths []string) []*execpb.DirectoryMount {
  seen := make(map[string]bool)

  var dirs []*execpb.DirectoryMount
  for _, path := range writePaths {
    if seen[path] {
      continue
    }
    seen[path] = true
    dirs = append(dirs, &execpb.DirectoryMount{
      PathInsideContainer: path,
      PathOutsideContainer: path,
      Writable: true,
    })
  }
  for _, path := range readPaths {
    if seen[path] {
      continue
    }
    seen[path] = true
    dirs = append(dirs, &execpb.DirectoryMount{
      PathInsideContainer: path,
      PathOutsideContainer: path,
      Writable: false,
    })
  }
  dirs = append(dirs, &execpb.DirectoryMount{
    PathInsideContainer: "/etc",
    PathOutsideContainer: "/var/lib/omogen/etc",
    Writable: false,
  })
  return dirs
}

type RunArgs struct {
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

func CommandRunner(exec execpb.ExecuteService_ExecuteClient, args RunArgs) (*execpb.Termination, error) {
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
            Amount: 2 * args.MemoryLimitKb},
          &execpb.ResourceAmount{
            Type: execpb.ResourceType_PROCESSES,
            Amount: 10},
        },
      },
    },
    ContainerSpec: &execpb.ContainerSpec{
      Mounts: makeMounts(args.ExtraReadPaths, args.ExtraWritePaths),
    },
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
  return res.Termination, nil
}

func getUsage(amounts *execpb.ResourceAmounts, resourceType execpb.ResourceType) int64 {
  for _, amount := range amounts.Amounts {
    if amount.Type == resourceType {
      return amount.Amount
    }
  }
  // TODO propagate upwards
  logger.Fatalf("Missing type %v in %v", resourceType, amounts)
  return -1
}

func TerminationToResponse(termination *execpb.Termination) *runpb.RunResponse {
  response := &runpb.RunResponse{
    TimeUsageMs: getUsage(termination.UsedResources, execpb.ResourceType_CPU_TIME),
    MemoryUsageKb: getUsage(termination.UsedResources, execpb.ResourceType_MEMORY),
  }
  switch termination.Termination.(type) {
  case *execpb.Termination_Signal_:
    response.Exit = &runpb.RunResponse_Signaled{
      &runpb.Signaled{Signal: termination.GetSignal().Signal}}
  case *execpb.Termination_Exit_:
    response.Exit = &runpb.RunResponse_Exited{
      &runpb.Exited{ExitCode: termination.GetExit().Code}}
  case *execpb.Termination_ResourceExceeded:
    if termination.GetResourceExceeded() == execpb.ResourceType_CPU_TIME {
      response.Exit = &runpb.RunResponse_TimeLimitExceeded{
        &runpb.TimeLimitExceeded{}}
    } else if termination.GetResourceExceeded() == execpb.ResourceType_MEMORY {
      response.Exit = &runpb.RunResponse_MemoryLimitExceeded{
        &runpb.MemoryLimitExceeded{}}
    }
  default:
    // TODO propagate upwards
    logger.Fatalln("Invalid termination")
  }
  return response
}
