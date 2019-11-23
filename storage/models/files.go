package models

import (
	"database/sql"
	"github.com/jsannemo/omogenjudge/util/go/filestore"
)

// A StoredFile represents a stored file in the StoredFile table.
type StoredFile struct {
	// Hash is a hash of the binary data of the stored file.
	Hash string
	// URL is a descriptor of how to retrieve the file, understood by the filestore service.
	URL []byte
}

// FileData reads the data of the file from the file store.
func (s *StoredFile) FileData() ([]byte, error) {
	return filestore.GetFile(s.URL)
}

// FileString returns the data of the file as a string.
func (s *StoredFile) FileString() (string, error) {
	contents, err := s.FileData()
	return string(contents), err
}

// ToNilable converts a StoredFile to a NilableFile.
func (s *StoredFile) ToNilable() *NilableStoredFile {
	return &NilableStoredFile{
		sql.NullString{String: s.Hash, Valid: true},
		s.URL,
	}
}

// A NilableStoredFile is used to represent a file that may or may not be present.
type NilableStoredFile struct {
	Hash sql.NullString
	Url  []byte
}

// Nil checks if the stored file is absent.
func (s *NilableStoredFile) Nil() bool {
	return !s.Hash.Valid
}
