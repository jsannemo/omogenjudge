// A Language implementation of C++
package language

import (
	"context"
	"errors"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/logger"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
	"github.com/jsannemo/omogenjudge/runner/compilers"
	"github.com/jsannemo/omogenjudge/runner/runners"
	execpb "github.com/jsannemo/omogenjudge/sandbox/api"
)

func init() {
	logger.Infof("Initializing Python")
	initCpp17()
}

func first(output string) string {
	temp := strings.Split(output, "\n")
	return temp[0]
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
		// TODO clear files
		inf := tmp()
		outf := tmp()
		errf := tmp()
		termination, err := runners.Execute(stream, &runners.ExecArgs{
			Command:          executable,
			Args:             append(files, version, "-Ofast", "-static"),
			WorkingDirectory: outputPath,
			InputPath:        inf,
			OutputPath:       outf,
			ErrorPath:        errf,
			ExtraReadPaths:   []string{"/usr/include", "/usr/lib/gcc"},
			ExtraWritePaths:  []string{outputPath},
			// TODO: revisit limits
			TimeLimitMs:   10000,
			MemoryLimitKb: 500 * 1024,
		})
		if err != nil {
			return nil, err
		}
		if termination.ExitReason != runners.Exited || termination.ExitCode != 0 {
			return nil, errors.New("Compiler crashed :(")
		}
		return &runpb.CompiledProgram{
			ProgramRoot:   outputPath,
			CompiledPaths: []string{"a.out"},
			LanguageId:    req.LanguageId,
		}, nil
	}
}

func runCpp(executable string) runners.RunFunc {
	argFunc := func(prog *runpb.CompiledProgram) *runners.CommandArgs {
		return &runners.CommandArgs{
			Command:          filepath.Join(prog.ProgramRoot, prog.CompiledPaths[0]),
			WorkingDirectory: prog.ProgramRoot,
		}
	}
	return runners.CommandProgram(argFunc)
}

func initCpp(executable, name, tag, versionFlag string, languageGroup runpb.LanguageGroup) {
	logger.Infof("Checking for C++ executable %s", executable)
	realPath, err := exec.LookPath(executable)
	// TODO check why this failed
	if err != nil {
		return
	}
	version, err := runners.FirstLineFromCommand(realPath, []string{"--version"})
	if err != nil {
		logger.Fatalf("Failed retreiving C++ version: %v", err)
	}
	language := &Language{
		Id:            tag,
		Version:       version,
		LanguageGroup: languageGroup,
		Compile:       cppCompile(realPath, versionFlag),
		Program:       runCpp(realPath),
	}
	registerLanguage(language)
}

func initCpp17() {
	initCpp("g++", "GCC C++17", "gpp17", "--std=gnu++17", runpb.LanguageGroup_CPP_17)
}
