// Database actions relating to stored files.
package files

import (
  "context"

  "github.com/jsannemo/omogenjudge/storage/db"
)

// CreateFile inserts the given file into the database.
// If the file already existed, the URL will be updated instead.
func CreateFile(ctx context.Context, file *StoredFile) error {
  conn := db.GetPool()
  _, err := conn.ExecContext(ctx, "INSERT INTO stored_file(file_hash.hash, url) VALUES($1, $2) ON CONFLICT(file_hash) DO UPDATE SET url = $2", file.Hash, file.Url)
  return err
}
