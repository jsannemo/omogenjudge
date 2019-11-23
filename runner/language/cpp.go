package language

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"

	"github.com/google/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
	"github.com/jsannemo/omogenjudge/runner/compilers"
	"github.com/jsannemo/omogenjudge/runner/runners"
	execpb "github.com/jsannemo/omogenjudge/sandbox/api"
	"github.com/jsannemo/omogenjudge/util/go/commands"
	"github.com/jsannemo/omogenjudge/util/go/files"
	"github.com/jsannemo/omogenjudge/util/go/strings"
	"github.com/jsannemo/omogenjudge/util/go/users"
)

func initCpp() error {
	logger.Infof("Initializing C++")
	if err := initCpp17(); err != nil {
		return err
	}
	return nil
}

func cppCompile(executable, version string) CompileFunc {
	return func(req *runpb.Program, outputPath string, exec execpb.ExecuteServiceClient) (*compilers.Compilation, error) {
		programFiles, err := compilers.WriteProgramToDisc(req, outputPath)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed writing program to disc: %v", err)
		}
		// TODO used a propagated context
		stream, err := exec.Execute(context.TODO())
		defer stream.CloseSend()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed opening exec connection: %v", err)
		}
		fb := files.NewFileBase("/var/lib/omogen/tmps")
		fb.Gid = users.OmogenClientsID()
		inFile, outFile, errFile := strings.RandStr(8), strings.RandStr(8), strings.RandStr(8)
		inf, err := fb.FullPath(inFile)
		if err != nil {
			panic(err)
		}
		outf, err := fb.FullPath(outFile)
		if err != nil {
			panic(err)
		}
		errf, err := fb.FullPath(errFile)
		if err != nil {
			panic(err)
		}
		defer fb.Remove(inf)
		defer fb.Remove(outf)
		defer fb.Remove(errf)
		if err := fb.WriteFile(inFile, []byte{}); err != nil {
			return nil, err
		}
		fb.GroupWritable = true
		if err := fb.WriteFile(outFile, []byte{}); err != nil {
			return nil, err
		}
		if err := fb.WriteFile(errFile, []byte{}); err != nil {
			return nil, err
		}
		termination, err := runners.Execute(stream, &runners.ExecArgs{
			Command:          executable,
			Args:             append(programFiles, version, "-Ofast", "-static"),
			WorkingDirectory: outputPath,
			InputPath:        inf,
			OutputPath:       outf,
			ErrorPath:        errf,
			ExtraReadPaths:   []string{"/usr/include", "/usr/lib/gcc"},
			ExtraWritePaths:  []string{outputPath},
			// TODO: revisit these limits
			TimeLimitMs:   10000,
			MemoryLimitKb: 1000 * 1000,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed executing compilation command: %v", err)
		}

		compileOut, err := ioutil.ReadFile(outf)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed opening compiler output: %v", err)
		}
		compileErr, err := ioutil.ReadFile(errf)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed opening compiler errors: %v", err)
		}

		var program *runpb.CompiledProgram
		if termination.CrashedWith(0) {
			program = &runpb.CompiledProgram{
				ProgramRoot:   outputPath,
				CompiledPaths: []string{"a.out"},
				LanguageId:    req.LanguageId,
			}
		}
		return &compilers.Compilation{
			Program: program,
			Output:  string(compileOut),
			Errors:  string(compileErr),
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

func initCppVersion(executable, name, tag, versionFlag string, languageGroup runpb.LanguageGroup) error {
	logger.Infof("Checking for C++ executable %s", executable)
	realPath, err := exec.LookPath(executable)
	if err != nil {
		logger.Infof("Could not find G++: %v", err)
		return nil
	}
	version, err := commands.FirstLineFromCommand(realPath, []string{"--version"})
	if err != nil {
		return fmt.Errorf("Could not get G++ version: %v", err)
	}
	language := &Language{
		Id:            tag,
		Version:       version,
		LanguageGroup: languageGroup,
		Compile:       cppCompile(realPath, versionFlag),
		Program:       runCpp(realPath),
	}
	registerLanguage(language)
	return nil
}

func initCpp17() error {
	return initCppVersion("g++", "GCC C++17", "gpp17", "--std=gnu++17", runpb.LanguageGroup_CPP_17)
}
