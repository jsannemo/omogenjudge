// A Language implementation of running simple commands
package language

import (
  "path/filepath"

  "github.com/google/logger"

	"github.com/jsannemo/omogenjudge/runner/compilers"
	"github.com/jsannemo/omogenjudge/runner/runners"
	execpb "github.com/jsannemo/omogenjudge/sandbox/api"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

func init() {
  logger.Infof("Initializing Command running")
  initCmd()
}

func cmdRunFunc() RunFunc {
  return func(req *runpb.RunRequest, exec execpb.ExecuteService_ExecuteClient) (*runpb.RunResponse, error) {
    result, err := runners.CommandRunner(exec, runners.RunArgs{
      Command: req.Program.ProgramRoot,
      Args: req.Args,
      WorkingDirectory: filepath.Dir(req.InputPath),
      InputPath: req.InputPath,
      OutputPath: req.OutputPath,
      ErrorPath: req.ErrorPath,
      ExtraReadPaths: append(req.ExtraReadPaths, filepath.Dir(req.InputPath)),
      ExtraWritePaths: []string{filepath.Dir(req.OutputPath), filepath.Dir(req.ErrorPath),},
      TimeLimitMs: req.TimeLimitMs,
      MemoryLimitKb: req.MemoryLimitKb,
    })
    if err != nil {
      return nil, err
    }
    return runners.TerminationToResponse(result), nil
  }
}

func initCmd() {
  language := &Language{
    Id: "cmd",
    Version: "",
    LanguageGroup: runpb.LanguageGroup_CMD,
    Compile: compilers.Noop,
    Run: cmdRunFunc(),
  }
  registerLanguage(language)
}
