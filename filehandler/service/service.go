package service

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"google.golang.org/grpc"

	filepb "github.com/jsannemo/omogenjudge/filehandler/api"
	"github.com/jsannemo/omogenjudge/util/go/filestore"
)

var (
	cachePath = flag.String("file_cache_path", "/var/lib/omogen/filecache", "The folder to used to store cached files")
)

type FileServer struct {
}

func ensureFile(handle *filepb.FileHandle) (string, error) {
	dir := filepath.Join(*cachePath, handle.Sha256Hash)
	storagePath := filepath.Join(dir, handle.Sha256Hash)

	stat, err := os.Stat(dir)
	if err == nil && !stat.IsDir() {
		return "", fmt.Errorf("File directory %s was an existing file", dir)
	} else if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(dir, 0707)
			if err != nil && !os.IsExist(err) {
				return "", err
			}
		}
		return "", err
	}
	stat, err = os.Stat(storagePath)
	if err == nil && stat.IsDir() {
		return "", fmt.Errorf("File cache was directory")
	} else if err != nil {
		if os.IsNotExist(err) {
			contents, err := filestore.GetFile(handle.Url)
			if err != nil {
				return "", err
			}
			err = ioutil.WriteFile(storagePath, contents, 0644)
			if err != nil {
				return "", err
			}
		}
	}
	return storagePath, nil
}

func (s *FileServer) EnsureFile(ctx context.Context, req *filepb.EnsureFileRequest) (*filepb.EnsureFileResponse, error) {
	var paths []string
	for _, file := range req.Handles {
		path, err := ensureFile(file)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	ret := &filepb.EnsureFileResponse{Paths: paths}
	return ret, nil
}

func Register(grpcServer *grpc.Server) {
  server := &FileServer{}
  filepb.RegisterFileHandlerServiceServer(grpcServer, server)
}
