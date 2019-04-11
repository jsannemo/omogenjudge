package main

import (
  "flag"
  "context"
  "net"
  "io"
  "log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jsannemo/omogenexec/run/language"
	execpb "github.com/jsannemo/omogenexec/exec/api"
	runpb "github.com/jsannemo/omogenexec/run/api"
)

var (
  address = flag.String("listen", "127.0.0.1:61811", "The server address")
  execAddress = flag.String("exec_server", "127.0.0.1:61810", "The exec server address")
)

type runServer struct {
	languages []*runpb.Language
  exec execpb.ExecuteServiceClient 
}


func (s *runServer) GetLanguages(ctx context.Context, _ *runpb.GetLanguagesRequest) (*runpb.GetLanguagesResponse, error) {
	return &runpb.GetLanguagesResponse{InstalledLanguages: s.languages}, nil
}

func (s *runServer) Compile(ctx context.Context, req *runpb.CompileRequest) (*runpb.CompileResponse, error) {
  log.Printf("Request(Compile): %v", req)
  // TODO: add request validation
  language, exists := language.GetLanguage(req.Program.LanguageId)
  if !exists {
    return nil, status.Errorf(codes.InvalidArgument, "Language %v does not exist", req.Program.LanguageId)
  }
  compiledProgram, err := language.Compile(req.Program, req.OutputPath, s.exec)
  if err != nil {
    log.Printf("Response(Compile): ERROR: %v", err)
    return nil, err
  }
  response := &runpb.CompileResponse{
    Program: compiledProgram,
  }
  log.Printf("Response(Compile): %v", response)
  return response, nil
}

func (s *runServer) Run(stream runpb.RunService_RunServer) error {
  execStream, err := s.exec.Execute(context.Background())
  defer execStream.CloseSend()
  if err != nil {
    log.Fatalf("Could not open exec stream: %v", err)
  }
	for {
		req, err := stream.Recv()
		if err == io.EOF {
      return nil
		}
		if err != nil {
      return err
		}
    language, exists := language.GetLanguage(req.Program.LanguageId)
    if !exists {
      return status.Errorf(codes.InvalidArgument, "Language %v does not exist", req.Program.LanguageId)
    }
    response, err := language.Run(req, execStream)
    if err != nil {
      return err
    }
		stream.Send(response)
  }
}

func newExecClient() execpb.ExecuteServiceClient {
  var conn *grpc.ClientConn
  conn, err := grpc.Dial(*execAddress, grpc.WithInsecure())
  if err != nil {
    log.Fatalf("Could not create ExecuteService client: %s", err)
  }
  client := execpb.NewExecuteServiceClient(conn)
  return client
}

func newServer() (*runServer, error) {
  apiLanguages := make([]*runpb.Language, 0)
  for _, language := range language.GetLanguages() {
    apiLanguages = append(apiLanguages, language.ToApiLanguage())
  }
	s := &runServer{
		languages: apiLanguages,
    exec: newExecClient(),
	}
	return s, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
  server, err := newServer()
  if err != nil {
    log.Fatalf("failed to create server: %v", err)
  }
	runpb.RegisterRunServiceServer(grpcServer, server)
	grpcServer.Serve(lis)
}

