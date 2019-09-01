package language

import (
	"os/exec"

	"github.com/google/logger"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
	"github.com/jsannemo/omogenjudge/runner/compilers"
	"github.com/jsannemo/omogenjudge/runner/runners"
	"github.com/jsannemo/omogenjudge/util/go/commands"
)

func init() {
	logger.Infof("Initializing Python")
	initPypy2()
	initPypy3()
	initPython2()
	initPython3()
}

func runPython(executable string) runners.RunFunc {
	argFunc := func(prog *runpb.CompiledProgram) *runners.CommandArgs {
		return &runners.CommandArgs{
			Command:          executable,
			Args:             prog.CompiledPaths,
			WorkingDirectory: prog.ProgramRoot,
		}
	}
	return runners.CommandProgram(argFunc)
}

func initPython(executable, name, tag string, languageGroup runpb.LanguageGroup) {
	logger.Infof("Checking for Python executable %s", executable)
	realPath, err := exec.LookPath(executable)
	if err != nil {
		// TODO: check if error was because of something other than not existing
		return
	}
	version, err := commands.FirstLineFromCommand(realPath, []string{"--version"})
	if err != nil {
		logger.Fatalf("Could not retrieve version for python %v", realPath)
	}
	language := &Language{
		Id:            tag,
		Version:       version,
		LanguageGroup: languageGroup,
		Compile:       compilers.Copy,
		Program:       runPython(realPath),
	}
	registerLanguage(language)
}

func initPypy2() {
	initPython("pypy", "Python 2 (PyPy)", "pypy2", runpb.LanguageGroup_PYTHON_2_PYPY)
}

func initPypy3() {
	initPython("pypy3", "Python 3 (PyPy)", "pypy3", runpb.LanguageGroup_PYTHON_3_PYPY)
}

func initPython2() {
	initPython("python2", "Python 2 (CPython)", "cpython2", runpb.LanguageGroup_PYTHON_2)
}

func initPython3() {
	initPython("python3", "Python 3 (CPython)", "cpython3", runpb.LanguageGroup_PYTHON_3)
}
