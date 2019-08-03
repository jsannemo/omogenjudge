// filestore handles storing and reading arbitrary binary files.
package filestore

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/google/logger"
)

var ErrInvalidUrl = errors.New("The file storage URL could not be recognized")

var binPrefix = []byte("bin:")

// StoreFile stores the given byte strings, returning a hash and URL where the file can be loaded from using GetFile.
func StoreFile(contents []byte) (string, []byte, error) {
	// TODO use something more reasonable for this
	hash := sha256.Sum256(contents)
	url := []byte("bin:")
	url = append(url, contents...)
	return hex.EncodeToString(hash[:]), url, nil
}

// GetFile retrieves the contents of a file given a URL generated by StoreFile.
func GetFile(url []byte) ([]byte, error) {
	logger.Infof("url %v, prefix %v", url, binPrefix)
	if eq(url[:4], binPrefix, 4) {
		return url[4:], nil
	}
	return nil, ErrInvalidUrl
}

func eq(a []byte, b []byte, l int) bool {
	for i := 0; i < l; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
