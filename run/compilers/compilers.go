package compilers

import (
  "path"
  "io/ioutil"

	execpb "github.com/jsannemo/omogenexec/exec/api"
	runpb "github.com/jsannemo/omogenexec/run/api"
)

func writeFile(path string, contents []byte) error {
  return ioutil.WriteFile(path, contents, 0700)
}

func Copy(req *runpb.Program, outputPath string, _ execpb.ExecuteServiceClient) (*runpb.CompiledProgram, error) {
  compiledPaths := []string{}

  for _, file := range req.Sources {
    writeFile(path.Join(outputPath, file.Path), file.Contents)
    compiledPaths = append(compiledPaths, file.Path)
  }
  return &runpb.CompiledProgram{
    ProgramRoot: outputPath,
    CompiledPaths: compiledPaths,
    LanguageId: req.LanguageId,
  }, nil
}

func PrefixPaths(prefix string, strs []string) []string {
  var prefixedStrs []string
  for _, str := range(strs) {
    prefixedStrs = append(prefixedStrs, path.Join(prefix, str))
  }
  return prefixedStrs
}
