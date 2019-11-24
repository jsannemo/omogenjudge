// Package main contains a binary for installing contests.
package main

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/logger"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	ptclient "github.com/jsannemo/omogenjudge/problemtools/client"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
	"github.com/jsannemo/omogenjudge/storage/files"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/problems"
	"github.com/jsannemo/omogenjudge/util/go/cli"
	futil "github.com/jsannemo/omogenjudge/util/go/files"
	"github.com/jsannemo/omogenjudge/util/go/filestore"
	"github.com/jsannemo/omogenjudge/util/go/users"
)

func main() {
	flag.Parse()
	defer logger.Init("addproblem", true, false, ioutil.Discard).Close()
	path := flag.Arg(0)
	path, err := filepath.Abs(path)
	if err != nil {
		logger.Fatal(err)
	}
	if err := installProblem(path); err != nil {
		logger.Fatalf("Failed installing problem: %v", err)
	}
}

func getOutputValidator(ctx context.Context, val *runpb.Program) (*models.OutputValidator, error) {
	if val == nil {
		return nil, nil
	}
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	for _, file := range val.Sources {
		f, err := w.Create(file.Path)
		if err != nil {
			return nil, err
		}
		if _, err = f.Write([]byte(file.Contents)); err != nil {
			return nil, err
		}
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	hash, url, err := filestore.StoreFile(buf.Bytes())
	if err != nil {
		return nil, err
	}
	storedFile := &models.StoredFile{Hash: hash, URL: url}
	if err := files.CreateFile(ctx, storedFile); err != nil {
		return nil, err
	}

	return &models.OutputValidator{
		ValidatorLanguageID: sql.NullString{String: val.LanguageId, Valid: true},
		ValidatorSourceZIP:  storedFile.ToNilable(),
	}, nil
}

func getStatements(statements []*toolspb.ProblemStatement) []*models.ProblemStatement {
	var storage []*models.ProblemStatement
	for _, s := range statements {
		storage = append(storage,
			&models.ProblemStatement{
				Language: s.LanguageCode,
				Title:    s.Title,
				HTML:     s.StatementHtml,
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
	if err := files.CreateFile(ctx, storedFile); err != nil {
		return nil, err
	}
	return storedFile, nil
}

func getTestCase(ctx context.Context, tc *toolspb.TestCase) (*models.TestCase, error) {
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

func getTestGroup(ctx context.Context, tc *toolspb.TestGroup) (*models.TestGroup, error) {
	var tests []*models.TestCase
	for _, test := range tc.Tests {
		storageTest, err := getTestCase(ctx, test)
		if err != nil {
			return nil, err
		}
		tests = append(tests, storageTest)
	}
	return &models.TestGroup{
		Name:             tc.Name,
		Score:            tc.Score,
		PublicVisibility: tc.PublicSamples,
		Tests:            tests}, nil
}

func getTestGroups(ctx context.Context, testGroups []*toolspb.TestGroup) ([]*models.TestGroup, error) {
	var groups []*models.TestGroup
	for _, group := range testGroups {
		storageGroup, err := getTestGroup(ctx, group)
		if err != nil {
			return nil, err
		}
		groups = append(groups, storageGroup)
	}
	return groups, nil
}

func makeInstallationDirectory(tmpDir, problemPath, problemName string) (string, error) {
	fb := futil.NewFileBase(tmpDir)
	fb.Gid = users.OmogenClientsID()
	fb.GroupWritable = true
	if err := fb.FixOwners("."); err != nil {
		return "", err
	}
	if err := fb.FixModeExec("."); err != nil {
		return "", err
	}
	if err := fb.Mkdir(problemName); err != nil {
		return "", fmt.Errorf("could not create installation problem directory: %v", err)
	}
	installfb, err := fb.SubBase(problemName)
	if err != nil {
		return "", err
	}
	if err := installfb.CopyInto(problemPath); err != nil {
		return "", fmt.Errorf("could not clone into installation problem directory: %v", err)
	}
	return installfb.Path(), nil
}

func installProblem(path string) error {
	logger.Infof("Installing problem %s", path)
	tmp, err := ioutil.TempDir("/tmp", "omogeninstall")
	if err != nil {
		return fmt.Errorf("could not create installation directory: %v", err)
	}
	defer os.RemoveAll(tmp)
	problemName := filepath.Base(path)
	npath, err := makeInstallationDirectory(tmp, path, problemName)

	ctx := context.Background()
	client := ptclient.NewClient()
	parsed, err := client.ParseProblem(ctx, &toolspb.ParseProblemRequest{
		ProblemPath: npath,
	})
	if err != nil {
		return fmt.Errorf("ParseProblem failed: %v", err)
	}
	for _, infoMsg := range parsed.Infos {
		logger.Infoln(infoMsg)
	}
	for _, warnMsg := range parsed.Warnings {
		logger.Warningln(warnMsg)
	}
	for _, errMsg := range parsed.Errors {
		logger.Errorln(errMsg)
	}
	if len(parsed.Errors) != 0 {
		return fmt.Errorf("problem had errors; will not install")
	}
	hasWarnings := len(parsed.Warnings) > 0
	problem := parsed.ParsedProblem
	verified, err := client.VerifyProblem(ctx, &toolspb.VerifyProblemRequest{
		ProblemToVerify: problem,
		ProblemPath:     npath,
	})
	if err != nil {
		return fmt.Errorf("VerifyProblem failed: %v", err)
	}
	problem = verified.VerifiedProblem
	hasWarnings = hasWarnings || len(verified.Warnings) > 0
	for _, infoMsg := range verified.Infos {
		logger.Infoln(infoMsg)
	}
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

	storageTestGroups, err := getTestGroups(ctx, problem.TestGroups)
	if err != nil {
		return err
	}
	outputValidator, err := getOutputValidator(ctx, problem.OutputValidator)
	if err != nil {
		return err
	}
	problemVersion := &models.ProblemVersion{
		TestGroups:      storageTestGroups,
		TimeLimMS:       problem.Metadata.Limits.TimeLimitMs,
		MemLimKB:        problem.Metadata.Limits.MemoryLimitKb,
		OutputValidator: outputValidator,
	}
	storageProblem := &models.Problem{
		ShortName:      problem.Metadata.ProblemId,
		Statements:     getStatements(problem.Statements),
		License:        models.License(problem.Metadata.License.String()),
		Author:         problem.Metadata.Author,
		CurrentVersion: problemVersion,
	}
	if err := problems.CreateProblem(ctx, storageProblem); err != nil {
		if err == problems.ErrDuplicateProblemName {
			if cli.RequestConfirmation(fmt.Sprintf("A problem named %s is already installed; update?", problem.Metadata.ProblemId)) {
				return problems.UpdateProblem(ctx, storageProblem)
			}
		}
		return err
	}
	return nil
}
