package runners

import (
	"os"
	"path/filepath"

	"github.com/jsannemo/omogenjudge/util/go/files"
)

// Env is a persistent environment across executions within the same sandbox.
type Env struct {
	ReadRoot  string
	WriteRoot string
}

// PathFor returns the path that a file will get inside the env.
func (e *Env) PathFor(inName string, writeable bool) string {
	var root string
	if writeable {
		root = e.WriteRoot
	} else {
		root = e.ReadRoot
	}
	newName := filepath.Join(root, inName)
	return newName
}

// LinkFile hard links the file path into the inside root.
func (e *Env) LinkFile(path, inName string, writeable bool) error {
	return os.Link(path, e.PathFor(inName, writeable))
}

// Clear resets the environment for a new execution.
func (e *Env) Clear() error {
	if err := files.RemoveContents(e.ReadRoot); err != nil {
		return err
	}
	if err := files.RemoveContents(e.WriteRoot); err != nil {
		return err
	}
	return nil
}

// NewEnv returns a new environment, rooted at the given path.
func NewEnv(envRoot string) (*Env, error) {
	if err := os.MkdirAll(filepath.Join(envRoot, "read"), 0755); err != nil {
		return nil, err
	}
  writePath := filepath.Join(envRoot, "write")
	if err := os.MkdirAll(writePath, 0755); err != nil {
		return nil, err
	}
  if err := os.Chmod(writePath, 0775); err != nil {
    return nil, err
  }
	return &Env{
		ReadRoot:  filepath.Join(envRoot, "read"),
		WriteRoot: filepath.Join(envRoot, "write"),
	}, nil
}
