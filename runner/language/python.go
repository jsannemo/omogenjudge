// A Language implementation of Python, both version 2 and 3 and CPython/pypy runtimes.
package language

import (
	"bytes"
  "os/exec"
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
  initPypy2()
  initPypy3()
  initPython2()
  initPython3()
}

func head(output string) string {
  temp := strings.Split(output,"\n")
  return temp[0]
}

func getVersion(path string) string {
  cmd := exec.Command(path, "--version")
	var stderr, stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
    logger.Fatalf("Failed retreiving python version: %v", err)
	}
  outLine := head(stdout.String())
  errLine := head(stderr.String())
  if len(outLine) != 0 {
    return outLine
  }
  if len(errLine) != 0 {
    return errLine
  }
  logger.Fatalf("Could not find a version for python %s", path)
  return ""
}

func runFunc(executable string) RunFunc {
  return func(req *runpb.RunRequest, exec execpb.ExecuteService_ExecuteClient) (*runpb.RunResponse, error) {
    result, err := runners.CommandRunner(exec, runners.RunArgs{
      Command: executable,
      Args: req.Program.CompiledPaths,
      WorkingDirectory: req.Program.ProgramRoot,
      InputPath: req.InputPath,
      OutputPath: req.OutputPath,
      ErrorPath: req.ErrorPath,
      ExtraReadPaths: []string{filepath.Dir(req.InputPath), req.Program.ProgramRoot,},
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

func initPython(executable, name, tag string, languageGroup runpb.LanguageGroup) {
  logger.Infof("Checking for Python executable %s", executable)
  realPath, err := exec.LookPath(executable)
  if err != nil {
    return
  }
  version := getVersion(realPath)
  language := &Language{
    Id: tag,
    Version: version,
    LanguageGroup: languageGroup,
    Compile: compilers.Copy,
    Run: runFunc(realPath),
  }
  registerLanguage(language)
}

func initPypy2() {
  initPython("pypy", "Python 2 (PyPy)", "pypy2", runpb.LanguageGroup_PYTHON_2)
}

func initPypy3() {
  initPython("pypy3", "Python 3 (PyPy)", "pypy3", runpb.LanguageGroup_PYTHON_3)
}

func initPython2() {
  initPython("python2", "Python 2 (CPython)", "cpython2", runpb.LanguageGroup_PYTHON_2)
}

func initPython3() {
  initPython("python3", "Python 3 (CPython)", "cpython3", runpb.LanguageGroup_PYTHON_3)
}
