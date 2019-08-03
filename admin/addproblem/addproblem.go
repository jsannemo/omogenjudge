package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/logger"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	ptclient "github.com/jsannemo/omogenjudge/problemtools/client"
	"github.com/jsannemo/omogenjudge/storage/files"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/problems"
	"github.com/jsannemo/omogenjudge/util/go/cli"
	futil "github.com/jsannemo/omogenjudge/util/go/files"
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
		Name:       tc.Name,
		InputFile:  inputFile,
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
		Name:             tc.Name,
		PublicVisibility: tc.PublicSamples,
		Tests:            tests}, nil
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

func installProblem(path string) error {
	logger.Infof("Installing problem %s", path)

	tmp, err := ioutil.TempDir("/tmp", "omogeninstall")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	if err := os.Chmod(tmp, 0755); err != nil {
		return err
	}
	npath := filepath.Join(tmp, filepath.Base(path))
	if err := futil.CopyDirectory(path, npath); err != nil {
		return err
	}

	ctx := context.Background()

	client := ptclient.NewClient()
	parsed, err := client.ParseProblem(ctx, &toolspb.ParseProblemRequest{
		ProblemPath: npath,
	})
	if err != nil {
		return err
	}
	for _, warnMsg := range parsed.Warnings {
		logger.Warningln(warnMsg)
	}
	for _, errMsg := range parsed.Errors {
		logger.Errorln(errMsg)
	}
	if len(parsed.Errors) != 0 {
		return fmt.Errorf("Problem had errors; will not install")
	}
	hasWarnings := len(parsed.Warnings) > 0
	problem := parsed.ParsedProblem
	verified, err := client.VerifyProblem(ctx, &toolspb.VerifyProblemRequest{
		ProblemToVerify: problem,
	})
	if err != nil {
		return err
	}
	hasWarnings = hasWarnings || len(verified.Warnings) > 0
	for _, warnMsg := range verified.Warnings {
		logger.Warningln(warnMsg)
	}
	for _, errMsg := range verified.Errors {
		logger.Errorln(errMsg)
	}
	if len(verified.Errors) != 0 {
		return fmt.Errorf("Problem had errors; will not install")
	}
	if hasWarnings {
		if !cli.RequestConfirmation("The problem had warnings; do you still want to install it?") {
			return nil
		}
	}

	storageTestGroups, err := toStorageTestGroups(ctx, problem.TestGroups)
	if err != nil {
		return err
	}
	if err := problems.Create(ctx, &models.Problem{
		ShortName:  problem.Metadata.ProblemId,
		Statements: toStorageStatements(problem.Statements),
		TestGroups: storageTestGroups,
	}); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	path := flag.Arg(0)
	path, err := filepath.Abs(path)
	if err != nil {
		logger.Fatal(err)
	}
	if err := installProblem(path); err != nil {
		logger.Fatalf("Failed installing problem: %v", err)
	}
}
