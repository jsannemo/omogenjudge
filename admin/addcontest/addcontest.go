package main

import (
	"context"
	"flag"
	"io/ioutil"
	"path/filepath"

	"github.com/google/logger"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	ptclient "github.com/jsannemo/omogenjudge/problemtools/client"
)

func main() {
	flag.Parse()
	defer logger.Init("addproblem", true, false, ioutil.Discard).Close()
	path := flag.Arg(0)
	path, err := filepath.Abs(path)
	if err != nil {
		logger.Fatal(err)
	}
	config, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Fatal(err)
	}
	client := ptclient.NewClient()
	response, err := client.InstallContest(context.Background(),&toolspb.InstallContestRequest{ContestYaml: string(config)})
	if err != nil {
		logger.Fatal(err)
	}
	for _, warn := range response.Warnings {
		logger.Warningln(warn)
	}
	for _, errs := range response.Errors {
		logger.Errorln(errs)
	}
}

