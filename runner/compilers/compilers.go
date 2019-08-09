package compilers

import (
	"io/ioutil"
	"os"
	"path/filepath"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
	execpb "github.com/jsannemo/omogenjudge/sandbox/api"
)

type Compilation struct {
	Program *runpb.CompiledProgram
	Output  string
	Errors  string
}

func writeFile(path string, contents []byte) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}
  if err := os.Chmod(filepath.Dir(path), 0775); err != nil {
    return err
  }
	return ioutil.WriteFile(path, contents, 0755)
}

func WriteProgramToDisc(req *runpb.Program, outputPath string) ([]string, error) {
	compiledPaths := []string{}

	for _, file := range req.Sources {
		err := writeFile(filepath.Join(outputPath, file.Path), []byte(file.Contents))
		if err != nil {
			return nil, err
		}
		compiledPaths = append(compiledPaths, file.Path)
	}
	return compiledPaths, nil

}

// Copy produces a compiled progam that is simply the input program but with all files copied into the output path.
func Copy(req *runpb.Program, outputPath string, _ execpb.ExecuteServiceClient) (*Compilation, error) {
	compiledPaths, err := WriteProgramToDisc(req, outputPath)
	if err != nil {
		return nil, err
	}
	return &Compilation{
		Program: &runpb.CompiledProgram{
			ProgramRoot:   outputPath,
			CompiledPaths: compiledPaths,
			LanguageId:    req.LanguageId,
		}}, nil
}

func Noop(req *runpb.Program, outputPath string, _ execpb.ExecuteServiceClient) (*Compilation, error) {
	return nil, nil
}

// TODO move this to /util/go
func PrefixPaths(prefix string, strs []string) []string {
	var prefixedStrs []string
	for _, str := range strs {
		prefixedStrs = append(prefixedStrs, filepath.Join(prefix, str))
	}
	return prefixedStrs
}
