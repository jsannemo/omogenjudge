package models

import (
  "github.com/jsannemo/omogenjudge/util/go/filestore"
  filepb "github.com/jsannemo/omogenjudge/filehandler/api"
)

type StoredFile struct {
  Hash string

  Url []byte
}

// FileData reads the data of the file from the file store.
func (s *StoredFile) FileData() ([]byte, error) {
  contents, err := filestore.GetFile(s.Url)
  return contents, err
}

// FileString returns the data of the file as a string.
func (s *StoredFile) FileString() (string, error) {
  contents, err := s.FileData()
  return string(contents), err
}

type FileList []*StoredFile

func (s FileList) ToHandlerFiles() []*filepb.FileHandle {
  var handlerFiles []*filepb.FileHandle
  for _, file := range s {
    handlerFiles = append(handlerFiles, &filepb.FileHandle{Sha256Hash: file.Hash, Url: file.Url})
  }
  return handlerFiles
}