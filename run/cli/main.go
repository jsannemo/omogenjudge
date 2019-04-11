package main

import (
  "context"
  "flag"
  "log"
  "io/ioutil"

  "google.golang.org/grpc"

  runpb "github.com/jsannemo/omogenexec/run/api"
)

var (
  serverAddr = flag.String("server_addr", "127.0.0.1:61811", "The server address in the format of host:port")
)

func getLanguages(client runpb.RunServiceClient) {
  response, err := client.GetLanguages(context.Background(), &runpb.GetLanguagesRequest{})
  if err != nil {
    log.Fatalf("Could not fetch languages: %v", err)
  }
  log.Printf("Languages: %v", response)
}

func helloWorld(client runpb.RunServiceClient) {
  dir, err := ioutil.TempDir("", "program")
  if err != nil {
    log.Fatal(err)
  }
  compileResponse, err := client.Compile(context.Background(), &runpb.CompileRequest{
    Program: &runpb.Program{
      Sources: []*runpb.SourceFile{
        &runpb.SourceFile{
          Path: "helloworld.py",
          Contents: []byte(`print("Hello World!")`),
        },
      },
      LanguageId: "cpython2",
    },
    OutputPath: dir,
  })
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("Compilation result: %v", compileResponse)

	stream, err := client.Run(context.Background())
  indir, err := ioutil.TempDir("", "input")
  outdir, err := ioutil.TempDir("", "output")
  ioutil.WriteFile(indir + "/in", []byte{}, 0700)
  stream.Send(
    &runpb.RunRequest{
      Program: compileResponse.Program,
      InputPath: indir + "/in",
      OutputPath: outdir + "/out",
      ErrorPath: outdir + "/err",
      TimeLimitMs: 1000,
      MemoryLimitKb: 256 * 1000,
    },
  )
  in, err := stream.Recv()
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("Run result: %v", in)
	stream.CloseSend()
}

func main() {
  flag.Parse()
  var opts []grpc.DialOption
  opts = append(opts, grpc.WithInsecure())
  conn, err := grpc.Dial(*serverAddr, opts...)
  if err != nil {
    log.Fatalf("fail to dial: %v", err)
  }
  defer conn.Close()
  client := runpb.NewRunServiceClient(conn)
  op := flag.Arg(0)
  log.Printf("Op: %s", op)
  if op == "langs" {
    getLanguages(client)
  }
  if op == "helloworld" {
    helloWorld(client)
  }
}
