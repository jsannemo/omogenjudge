package files

import (
  "os"
  "io/ioutil"
  "path/filepath"
)

func RemoveContents(dirPath string) error {
  dir, err := ioutil.ReadDir(dirPath)
  if err != nil {
    return err
  }
  for _, d := range dir {
    if err := os.RemoveAll(filepath.Join(dirPath, d.Name())); err != nil {
      return err
    }
  }
  return nil
}
