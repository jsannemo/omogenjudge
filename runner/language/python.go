package language

import (
	"fmt"
	"os/exec"

	"github.com/google/logger"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
	"github.com/jsannemo/omogenjudge/runner/compilers"
	"github.com/jsannemo/omogenjudge/runner/runners"
	"github.com/jsannemo/omogenjudge/util/go/commands"
)

func initPython() error {
	logger.Infof("Initializing Python")
	if err := initPypy2(); err != nil {
		return err
	}
	if err := initPypy3(); err != nil {
		return err
	}
	if err := initPython2(); err != nil {
		return err
	}
	if err := initPython3(); err != nil {
		return err
	}
	return nil
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

func initPythonVersion(executable, name, tag string, languageGroup runpb.LanguageGroup) error {
	logger.Infof("Checking for Python executable %s", executable)
	realPath, err := exec.LookPath(executable)
	if err != nil {
		logger.Infof("Could not find python executable: %v", err)
		return nil
	}
	version, err := commands.FirstLineFromCommand(realPath, []string{"--version"})
	if err != nil {
		return fmt.Errorf("Could not retrieve version for python %s: %v", realPath, err)
	}
	language := &Language{
		Id:            tag,
		Version:       version,
		LanguageGroup: languageGroup,
		Compile:       compilers.Copy,
		Program:       runPython(realPath),
	}
	registerLanguage(language)
	return nil
}

func initPypy2() error {
	return initPythonVersion("pypy", "Python 2 (PyPy)", "pypy2", runpb.LanguageGroup_PYTHON_2_PYPY)
}

func initPypy3() error {
	return initPythonVersion("pypy3", "Python 3 (PyPy)", "pypy3", runpb.LanguageGroup_PYTHON_3_PYPY)
}

func initPython2() error {
	return initPythonVersion("python2", "Python 2 (CPython)", "cpython2", runpb.LanguageGroup_PYTHON_2)
}

func initPython3() error {
	return initPythonVersion("python3", "Python 3 (CPython)", "cpython3", runpb.LanguageGroup_PYTHON_3)
}
