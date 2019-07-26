package main

import (
  "context"
  "flag"
  "io/ioutil"
  "net"
  "path/filepath"
  "os"
  "errors"

	"google.golang.org/grpc"
  "github.com/google/logger"

  "github.com/jsannemo/omogenjudge/util/go/filestore"
	filepb "github.com/jsannemo/omogenjudge/filehandler/api"
)

var (
  address = flag.String("file_listen_addr", "127.0.0.1:61814", "The file server address to listen to in the format host:port")
  cachePath = flag.String("file_cache_path", "/var/lib/omogen/filecache", "The folder to used to store cached files")
)

type fileServer struct {
}

func ensureFile(handle *filepb.FileHandle) (string, error) {
  dir := filepath.Join(*cachePath, handle.Sha256Hash)
  storagePath := filepath.Join(dir, handle.Sha256Hash)

	stat, err := os.Stat(dir)
  if err == nil && !stat.IsDir() {
    return "", errors.New("File directory was an existing file")
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
    return "", errors.New("File cache was directory")
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


// Implementation of FileHandlerServer.EnsureFile.
func (s *fileServer) EnsureFile(ctx context.Context, req *filepb.EnsureFileRequest) (*filepb.EnsureFileResponse, error) {
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

func newServer() (*fileServer, error) {
	s := &fileServer{}
	return s, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *address)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
  server, err := newServer()
  if err != nil {
    logger.Fatalf("failed to create server: %v", err)
  }
	filepb.RegisterFileHandlerServiceServer(grpcServer, server)
	grpcServer.Serve(lis)
}

