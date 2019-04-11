package language

import (
	"bytes"
  "log"
  "os/exec"
  "strings"
  "path/filepath"

	"github.com/jsannemo/omogenexec/run/compilers"
	"github.com/jsannemo/omogenexec/run/runners"

	execpb "github.com/jsannemo/omogenexec/exec/api"
	runpb "github.com/jsannemo/omogenexec/run/api"
)

func init() {
  log.Printf("Initializing Python")
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
		log.Fatal(err)
	}
  outLine := head(stdout.String())
  errLine := head(stderr.String())
  if len(outLine) != 0 {
    return outLine
  }
  if len(errLine) != 0 {
    return errLine
  }
  log.Fatalf("Did not find a version for %s", path)
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
  log.Printf("Checking for Python executable %s", executable)
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
