// Package compilers provides utilities for compiling programs.
package compilers

import (
	"fmt"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
	execpb "github.com/jsannemo/omogenjudge/sandbox/api"
	"github.com/jsannemo/omogenjudge/util/go/files"
	"github.com/jsannemo/omogenjudge/util/go/users"
)

// A Compilation represents the output of a program compilation.
type Compilation struct {
	// The compiled program. This is unset if the compilation failed.
	Program *runpb.CompiledProgram

	// The output printed by the compiler to stdout.
	Output string

	// The output printed by the compiler to stderr.
	Errors string
}

// WriteProgramToDisc writes the source files in the given program to disc
func WriteProgramToDisc(req *runpb.Program, outputPath string) ([]string, error) {
	compiledPaths := []string{}
	fb := files.NewFileBase(outputPath)
	fb.Gid = users.OmogenClientsID()
	fb.GroupWritable = true
	if err := fb.Mkdir("."); err != nil {
		return nil, fmt.Errorf("failed mkdir %s: %v", outputPath, err)
	}
	for _, file := range req.Sources {
		err := fb.WriteFile(file.Path, []byte(file.Contents))
		if err != nil {
			return nil, fmt.Errorf("failed writing %s: %v", file.Path, err)
		}
		compiledPaths = append(compiledPaths, file.Path)
	}
	return compiledPaths, nil

}

// Copy produces a compiled progam that is simply the input program but with all files copied into the output path.
func Copy(req *runpb.Program, outputPath string, _ execpb.ExecuteServiceClient) (*Compilation, error) {
	compiledPaths, err := WriteProgramToDisc(req, outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed writing program to disc: %v", err)
	}
	return &Compilation{
		Program: &runpb.CompiledProgram{
			ProgramRoot:   outputPath,
			CompiledPaths: compiledPaths,
			LanguageId:    req.LanguageId,
		}}, nil
}
