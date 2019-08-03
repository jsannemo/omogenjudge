// Database actions relating to stored files.
package files

import (
	"context"

	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

// CreateFile inserts the given file into the database.
// If the file already existed, the URL will be updated instead.
func Create(ctx context.Context, file *models.StoredFile) {
	db.Conn().MustExecContext(ctx, "INSERT INTO stored_file(file_hash.hash, url) VALUES($1, $2) ON CONFLICT(file_hash) DO UPDATE SET url = $2", file.Hash, file.Url)
}
