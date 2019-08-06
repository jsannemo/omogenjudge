package models

import (
	"database/sql"
	filepb "github.com/jsannemo/omogenjudge/filehandler/api"
	"github.com/jsannemo/omogenjudge/util/go/filestore"
)

type NilableStoredFile struct {
	Hash sql.NullString

	Url []byte
}

func (s *NilableStoredFile) NotNil() bool {
	return s.Hash.Valid
}

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

func (s *StoredFile) ToNilable() *NilableStoredFile {
	return &NilableStoredFile{
		sql.NullString{String: s.Hash, Valid: true},
		s.Url,
	}
}
