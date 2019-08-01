// addproblem parses, verifies and installs a given problem on a judging system.
package main

import (
	"context"
	"flag"
	"path/filepath"
	"io/ioutil"

	"github.com/google/logger"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	ptclient "github.com/jsannemo/omogenjudge/problemtools/client"
	"github.com/jsannemo/omogenjudge/storage/problems"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/files"
	"github.com/jsannemo/omogenjudge/util/go/cli"
	"github.com/jsannemo/omogenjudge/util/go/filestore"
)

func toStorageStatements(statements []*toolspb.ProblemStatement) []*models.ProblemStatement {
	var storage []*models.ProblemStatement
	for _, s := range statements {
		storage = append(storage,
			&models.ProblemStatement{
				Language: s.LanguageCode,
				Title:    s.Title,
				Html:     s.StatementHtml,
			})
	}
	return storage
}

func insertFile(ctx context.Context, path string) (*models.StoredFile, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
  hash, url, err := filestore.StoreFile(dat)
  if err != nil {
    return nil, err
  }
  storedFile := &models.StoredFile{hash, url}
  files.Create(ctx, storedFile)
  if err != nil {
    return nil, err
  }
  return storedFile, nil
}


func toStorageTest(ctx context.Context, tc *toolspb.TestCase) (*models.TestCase, error) {
  inputFile, err := insertFile(ctx, tc.InputPath)
  if err != nil {
    return nil, err
  }
  outputFile, err := insertFile(ctx, tc.OutputPath)
  if err != nil {
    return nil, err
  }
  return &models.TestCase{
    Name: tc.Name,
    InputFile: inputFile,
    OutputFile: outputFile,
  }, nil
}

func toStorageTestGroup(ctx context.Context, tc *toolspb.TestGroup) (*models.TestGroup, error) {
  var tests []*models.TestCase
  for _, test := range tc.Tests {
    storageTest, err := toStorageTest(ctx, test)
    if err != nil {
      return nil, err
    }
    tests = append(tests, storageTest)
  }
  return &models.TestGroup{
    Name: tc.Name,
    PublicVisibility: tc.PublicSamples,
    Tests: tests}, nil
}

func toStorageTestGroups(ctx context.Context, testGroups []*toolspb.TestGroup) ([]*models.TestGroup, error) {
  var groups []*models.TestGroup
  for _, group := range testGroups {
    storageGroup, err := toStorageTestGroup(ctx, group)
    if err != nil {
      return nil, err
    }
    groups = append(groups, storageGroup)
  }
  return groups, nil
}

func main() {
	flag.Parse()
	path := flag.Arg(0)
	path, err := filepath.Abs(path)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Installing problem %s", path)

	ctx := context.Background()

	client := ptclient.NewClient()
	parsed, err := client.ParseProblem(ctx, &toolspb.ParseProblemRequest{
		ProblemPath: path,
	})
	if err != nil {
		logger.Fatalf("Failed parsing problem: %v", err)
	}
	for _, warnMsg := range parsed.Warnings {
		logger.Warningln(warnMsg)
	}
	for _, errMsg := range parsed.Errors {
		logger.Errorln(errMsg)
	}
	if len(parsed.Errors) != 0 {
		logger.Errorf("Problem had errors; will not install")
		return
	}
  if len(parsed.Warnings) != 0 {
    if !cli.RequestConfirmation("The problem had warnings; do you still want to install it?") {
      return
    }
  }

	problem := parsed.ParsedProblem
  storageTestGroups, err := toStorageTestGroups(ctx, problem.TestGroups)
  if err != nil {
    logger.Fatalf("Failed converting test groups: %v", err)
  }
	err = problems.Create(ctx, &models.Problem{
		ShortName:  problem.Metadata.ProblemId,
		Statements: toStorageStatements(problem.Statements),
    TestGroups: storageTestGroups,
	})
	if err != nil {
		logger.Fatalf("Failed inserting problem: %v", err)
	}
}
