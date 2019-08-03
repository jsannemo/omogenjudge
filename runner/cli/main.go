// A command-line client that can be used ot test the runner service.
// Note that it must be run on the same machine as the runner service, since it
// creates and reads files with paths that must be shared with the service.
package main

import (
	"context"
	"flag"

	"github.com/google/logger"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
	rclient "github.com/jsannemo/omogenjudge/runner/client"
)

func getLanguages(client runpb.RunServiceClient) {
	response, err := client.GetLanguages(context.Background(), &runpb.GetLanguagesRequest{})
	if err != nil {
		logger.Fatalf("Could not fetch languages: %v", err)
	}
	logger.Infof("Languages: %v", response)
}

func main() {
	flag.Parse()
	client := rclient.NewClient()
	op := flag.Arg(0)
	if op == "langs" {
		getLanguages(client)
	}
}
