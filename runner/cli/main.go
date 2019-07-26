// A command-line client that can be used ot test the runner service.
// Note that it must be run on the same machine as the runner service, since it
// creates and reads files with paths that must be shared with the service.
package main

import (
  "context"
  "flag"
  "io/ioutil"

  "github.com/google/logger"

  rclient "github.com/jsannemo/omogenjudge/runner/client"
  runpb "github.com/jsannemo/omogenjudge/runner/api"
)

func getLanguages(client runpb.RunServiceClient) {
  response, err := client.GetLanguages(context.Background(), &runpb.GetLanguagesRequest{})
  if err != nil {
    logger.Fatalf("Could not fetch languages: %v", err)
  }
  logger.Infof("Languages: %v", response)
}

func helloWorld(client runpb.RunServiceClient) {
  dir, err := ioutil.TempDir("", "program")
  if err != nil {
    logger.Fatalf("Failed to create a local program", err)
  }
  compileResponse, err := client.Compile(context.Background(), &runpb.CompileRequest{
    Program: &runpb.Program{
      Sources: []*runpb.SourceFile{
        &runpb.SourceFile{
          Path: "helloworld.py",
          Contents: []byte(`print("Hello World!")`),
        },
      },
      LanguageId: "cpython3",
    },
    OutputPath: dir,
  })
  if err != nil {
    logger.Fatalf("Failed compiling program: %v", err)
  }
  logger.Infof("Compilation result: %v", compileResponse)

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
    logger.Fatalf("Run request failed: %v", err)
  }
  logger.Infof("Run result: %v", in)
	stream.CloseSend()
}

func main() {
  flag.Parse()
  client := rclient.NewClient()
  op := flag.Arg(0)
  if op == "langs" {
    getLanguages(client)
  } else if op == "helloworld" {
    helloWorld(client)
  }
}
