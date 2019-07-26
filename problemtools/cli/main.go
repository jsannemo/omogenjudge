package main

import (
	"context"
	"flag"

	"github.com/google/logger"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	pbclient "github.com/jsannemo/omogenjudge/problemtools/client"
)

func parseProblem(path string, client toolspb.ToolServiceClient) {
	compileResponse, err := client.ParseProblem(context.Background(), &toolspb.ParseProblemRequest{
		ProblemPath: path,
	})
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Infof("Result: %v", compileResponse)
}

func main() {
	flag.Parse()
	client := pbclient.NewClient()
	path := flag.Arg(0)
	logger.Infof("Parsing problem at path: %s", path)
	parseProblem(path, client)
}
