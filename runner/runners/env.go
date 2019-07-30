package runners

import (
  "os"
  "path/filepath"

  "github.com/jsannemo/omogenjudge/util/go/files"
)

type Env struct {
  ReadRoot string
  WriteRoot string
}

// LinkFile hard links the file path into the inside root.
func (e *Env) LinkFile(path, inName string, writeable bool) (string, error) {
  var root string
  if writeable {
    root = e.WriteRoot
  } else {
    root = e.ReadRoot
  }
  newName := filepath.Join(root, inName)
  err := os.Link(path, newName)
  return newName, err
}

func (e *Env) ClearEnv() error {
  if err := files.RemoveContents(e.ReadRoot); err != nil {
    return err
  }
  if err := files.RemoveContents(e.WriteRoot); err != nil {
    return err
  }
  return nil
}

func NewEnv(envRoot string) (*Env, error) {
  if err := os.Mkdir(filepath.Join(envRoot, "read"), 0755); err != nil {
    return nil, err
  }
  if err := os.Mkdir(filepath.Join(envRoot, "write"), 0755); err != nil {
    return nil, err
  }
  return &Env{
    ReadRoot: filepath.Join(envRoot, "read"),
    WriteRoot: filepath.Join(envRoot, "write"),
  }, nil
}
