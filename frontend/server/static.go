// Handles serving static resources
package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// registerStaticHandler maps the /static path to the resources folder.
// Only files are mapped (not folders) to disallow directory listings.
func registerStaticHandler(mux *mux.Router) {
	fs := directoryFilteringFileSystem{http.Dir("frontend/resources/")}
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(fs)))
}

// A file system that restricts opening files to only files rather than folders.
type directoryFilteringFileSystem struct {
	fs http.FileSystem
}

func (fs directoryFilteringFileSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil || s.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}
