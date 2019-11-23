package main

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	//"os"
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

func toStorageOutputValidator(ctx context.Context, val *runpb.Program) (*models.OutputValidator, error) {
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

func toStorageStatements(statements []*toolspb.ProblemStatement) []*models.ProblemStatement {
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
		Score: tc.Score,
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
	problemName := filepath.Base(path)

	tmp, err := ioutil.TempDir("/tmp", "omogeninstall")
	if err != nil {
		return fmt.Errorf("Could not create installation directory: %v", err)
	}
	//defer os.RemoveAll(tmp)
	fb := futil.NewFileBase(tmp)
	fb.Gid = users.OmogenClientsID()
	fb.GroupWritable = true
	if err := fb.FixOwners("."); err != nil {
		return err
	}
	if err := fb.FixModeExec("."); err != nil {
		return err
	}
	if err := fb.Mkdir(problemName); err != nil {
		return fmt.Errorf("Could not create installation problem directory: %v", err)
	}
	installfb, err := fb.SubBase(problemName)
	if err != nil {
		return err
	}
	if err := installfb.CopyInto(path); err != nil {
		return fmt.Errorf("Could not clone into installation problem directory: %v", err)
	}
	npath := installfb.Path()

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
		return fmt.Errorf("Problem had errors; will not install")
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

	storageTestGroups, err := toStorageTestGroups(ctx, problem.TestGroups)
	if err != nil {
		return err
	}
	outputValidator, err := toStorageOutputValidator(ctx, problem.OutputValidator)
	if err != nil {
		return err
	}
	problemVersion := &models.ProblemVersion{
		TestGroups:      storageTestGroups,
		TimeLimMS:       problem.Metadata.Limits.TimeLimitMs,
		MemLimKB:        problem.Metadata.Limits.MemoryLimitKb,
		OutputValidator: outputValidator,
	}
	if err := problems.CreateProblem(ctx, &models.Problem{
		ShortName:      problem.Metadata.ProblemId,
		Statements:     toStorageStatements(problem.Statements),
		License:        models.License(problem.Metadata.License.String()),
		Author:         problem.Metadata.Author,
		CurrentVersion: problemVersion,
	}); err != nil {
		return err
	}
	return nil
}
