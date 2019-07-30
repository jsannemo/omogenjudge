// A Language implementation of C++
package language

import (
	"bytes"
  "errors"
  "context"
  "os/exec"
  "io/ioutil"
  "strings"
  "path/filepath"

  "github.com/google/logger"

	"github.com/jsannemo/omogenjudge/runner/compilers"
	"github.com/jsannemo/omogenjudge/runner/runners"
	execpb "github.com/jsannemo/omogenjudge/sandbox/api"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

func init() {
  logger.Infof("Initializing Python")
  initCpp17()
}

func first(output string) string {
  temp := strings.Split(output,"\n")
  return temp[0]
}

func cppVersion(path string) string {
  cmd := exec.Command(path, "--version")
	var stderr, stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
    logger.Fatalf("Failed retreiving cpp version: %v", err)
	}
  outLine := head(stdout.String())
  errLine := head(stderr.String())
  if len(outLine) != 0 {
    return outLine
  }
  if len(errLine) != 0 {
    return errLine
  }
  logger.Fatalf("Could not find a cpp for python %s", path)
  return ""
}

func tmp() string {
  file, err := ioutil.TempFile("/tmp", "omogendummy")
  if err != nil {
    logger.Fatal(err)
  }
  return file.Name()
}

func cppCompile(executable, version string) CompileFunc {
  return func(req *runpb.Program, outputPath string, exec execpb.ExecuteServiceClient) (*runpb.CompiledProgram, error) {
    files, err := compilers.WriteProgramToDisc(req, outputPath)
    stream, err := exec.Execute(context.TODO())
    defer stream.CloseSend()
    if err != nil {
      return nil, err
    }
    inf := tmp()
    outf := tmp()
    errf := tmp()
    termination, err := runners.CommandRunner(stream, runners.RunArgs{
      Command: executable,
      Args: append(files, version, "-Ofast", "-static"),
      WorkingDirectory: outputPath,
      InputPath: inf,
      OutputPath: outf, 
      ErrorPath: errf,
      ExtraReadPaths: []string{filepath.Dir(inf), outputPath, "/usr/include", "/usr/lib/gcc"},
      ExtraWritePaths: []string{filepath.Dir(outf), outputPath},
      TimeLimitMs: 10000,
      MemoryLimitKb: 500 * 1024,
    })
    if err != nil {
      return nil, err
    }
    switch termination.Termination.(type) {
    case *execpb.Termination_Exit_:
      if termination.GetExit().Code != 0 {
        return nil, errors.New("Compiler crashed :(")
      }
		default:
      return nil, errors.New("Invalid exit for compiler :(")
		}
    return &runpb.CompiledProgram{
      ProgramRoot: outputPath,
      CompiledPaths: []string{"a.out"},
      LanguageId: req.LanguageId,
    }, nil
  }
}

func runSubmission() RunFunc {
  first := true
  return func(req *runpb.RunRequest, exec execpb.ExecuteService_ExecuteClient) (*runpb.RunResponse, error) {
    result, err := runners.CommandRunner(exec, runners.RunArgs{
      Command: filepath.Join(req.Program.ProgramRoot, req.Program.CompiledPaths[0]),
      WorkingDirectory: req.Program.ProgramRoot,
      InputPath: req.InputPath,
      OutputPath: req.OutputPath,
      ErrorPath: req.ErrorPath,
      ExtraReadPaths: []string{filepath.Dir(req.InputPath), req.Program.ProgramRoot,},
      ExtraWritePaths: []string{filepath.Dir(req.OutputPath), filepath.Dir(req.ErrorPath),},
      TimeLimitMs: req.TimeLimitMs,
      MemoryLimitKb: req.MemoryLimitKb,
      ReuseContainer: !first,
    })
    first = false
    if err != nil {
      return nil, err
    }
    return runners.TerminationToResponse(result), nil
  }
}

func initCpp(executable, name, tag string, languageGroup runpb.LanguageGroup) {
  logger.Infof("Checking for C++ executable %s", executable)
  realPath, err := exec.LookPath(executable)
  if err != nil {
    return
  }
  version := getVersion(realPath)
  language := &Language{
    Id: tag,
    Version: version,
    LanguageGroup: languageGroup,
    Compile: cppCompile(realPath, "--std=gnu++17"),
    Run: runSubmission,
  }
  registerLanguage(language)
}

func initCpp17() {
  initCpp("g++", "GCC C++17", "gpp17", runpb.LanguageGroup_CPP_17)
}
