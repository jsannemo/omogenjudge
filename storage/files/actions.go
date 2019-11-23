// files contains database actions relating to stored files.
package files

import (
	"context"
	"fmt"

	filepb "github.com/jsannemo/omogenjudge/filehandler/api"
	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

// CreateFile inserts the given file into the database.
// If the file already existed, the URL will be updated instead.
func CreateFile(ctx context.Context, file *models.StoredFile) error {
	query := `
		INSERT INTO
		    stored_file(file_hash, url)
		VALUES($1, $2)
		ON CONFLICT(file_hash) DO
		    UPDATE SET url = $2`
	_, err := db.Conn().ExecContext(ctx, query, file.Hash, file.URL)
	if err != nil {
		return fmt.Errorf("failed persisting file: %v", err)
	}
	return nil
}

// A FileList is a slice of StoredFiles
type FileList []*models.StoredFile

// ToHandlerFiles converts a stored file list into a list of handles for the FileHandlerService.
func (s FileList) ToHandlerFiles() []*filepb.FileHandle {
	var handlerFiles []*filepb.FileHandle
	for _, file := range s {
		handlerFiles = append(handlerFiles, &filepb.FileHandle{Sha256Hash: file.Hash, Url: file.URL})
	}
	return handlerFiles
}
