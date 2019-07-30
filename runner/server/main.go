package main

import (
  "context"
  "flag"
  "io"
  "io/ioutil"
  "net"
  "os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
  "github.com/google/logger"

	"github.com/jsannemo/omogenjudge/runner/diff"
	"github.com/jsannemo/omogenjudge/runner/runners"
	"github.com/jsannemo/omogenjudge/runner/language"
	execpb "github.com/jsannemo/omogenjudge/sandbox/api"
	eclient "github.com/jsannemo/omogenjudge/sandbox/client"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

var (
  address = flag.String("run_listen_addr", "127.0.0.1:61811", "The run server address to listen to in the format host:port")
)

type runServer struct {
	languages []*runpb.Language
  exec execpb.ExecuteServiceClient
}


// Implementation of RunServer.GetLanguages.
func (s *runServer) GetLanguages(ctx context.Context, _ *runpb.GetLanguagesRequest) (*runpb.GetLanguagesResponse, error) {
	return &runpb.GetLanguagesResponse{InstalledLanguages: s.languages}, nil
}

// Implementation of RunServer.Compile.
func (s *runServer) Compile(ctx context.Context, req *runpb.CompileRequest) (*runpb.CompileResponse, error) {
  logger.Infof("RunService.Compile: %v", req)
  // TODO: add request validation
  language, exists := language.GetLanguage(req.Program.LanguageId)
  if !exists {
    return nil, status.Errorf(codes.InvalidArgument, "Language %v does not exist", req.Program.LanguageId)
  }
  compiledProgram, err := language.Compile(req.Program, req.OutputPath, s.exec)
  if err != nil {
    logger.Errorf("Failed program compilation: %v", err)
    return nil, err
  }
  response := &runpb.CompileResponse{
    Program: compiledProgram,
  }
  return response, nil
}

// Implementation of RunServer.Run.
func (s *runServer) Run(stream runpb.RunService_RunServer) error {
  execStream, err := s.exec.Execute(context.Background())
  defer execStream.CloseSend()
  if err != nil {
    logger.Fatalf("Could not open exec stream: %v", err)
  }

  envDir, err := ioutil.TempDir("/var/lib/omogen/tmps", "env")
  if err != nil {
    return err
  }
  defer os.RemoveAll(envDir)
  env, err := runners.NewEnv(envDir)
  if err != nil {
    return err
  }

  var runFunc language.RunFunc
	for {
		req, err := stream.Recv()
		if err == io.EOF {
      return nil
		}
		if err != nil {
      return err
		}
    lang, exists := language.GetLanguage(req.Program.LanguageId)
    if !exists {
      return status.Errorf(codes.InvalidArgument, "Language %v does not exist", req.Program.LanguageId)
    }
    if runFunc == nil {
      runFunc = lang.Run()
    }

    req.InputPath, err = env.LinkFile(req.InputPath, "input", false)
    if err != nil {
      return err
    }
    req.OutputPath, err = env.LinkFile(req.OutputPath, "output", true)
    if err != nil {
      return err
    }
    req.ErrorPath, err = env.LinkFile(req.ErrorPath, "error", true)
    if err != nil {
      return err
    }
    response, err := runFunc(req, execStream)
    env.ClearEnv()
    if err != nil {
      return err
    }
		stream.Send(response)
  }
}

func (s *runServer) Diff(ctx context.Context, req *runpb.DiffRequest) (*runpb.DiffResponse, error) {
  refFile, err := os.Open(req.ReferenceOutputPath)
  if err != nil {
    if os.IsNotExist(err) {
      return nil, status.Error(codes.NotFound, "Reference output did not exist")
    }
    return nil, status.Errorf(codes.Internal, "Failed opening reference output: %v", err)
  }
  outFile, err:= os.Open(req.OutputPath)
  if err != nil {
    if os.IsNotExist(err) {
      return nil, status.Error(codes.NotFound, "Output did not exist")
    }
    return nil, status.Errorf(codes.Internal, "Failed opening output: %v", err)
  }
  diffRes, err := diff.Diff(refFile, outFile)
  if err != nil {
    return nil, status.Errorf(codes.Internal, "Failed diffing: %v", err)
  }
  return &runpb.DiffResponse{
    Matching: diffRes.Match,
    DiffDescription: diffRes.Description,
  }, nil
}

func newServer() (*runServer, error) {
  apiLanguages := make([]*runpb.Language, 0)
  for _, language := range language.GetLanguages() {
    apiLanguages = append(apiLanguages, language.ToApiLanguage())
  }
	s := &runServer{
		languages: apiLanguages,
    exec: eclient.NewClient(),
	}
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
	runpb.RegisterRunServiceServer(grpcServer, server)
	grpcServer.Serve(lis)
}

