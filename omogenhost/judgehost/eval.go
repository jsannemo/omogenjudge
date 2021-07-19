package main

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/logger"
	apipb "github.com/jsannemo/omogenexec/api"
	"github.com/jsannemo/omogenexec/eval"
	"github.com/jsannemo/omogenexec/util"
	"github.com/jsannemo/omogenhost/storage"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var langMap = map[string]apipb.LanguageGroup{
	"cpp":     apipb.LanguageGroup_CPP,
	"python3": apipb.LanguageGroup_PYTHON_3,
}

type submissionJson struct {
	Files map[string]string
}

var evalMutex sync.Mutex

func evaluate(runId int64) error {
	evalMutex.Lock()
	defer evalMutex.Unlock()

	var run storage.SubmissionRun
	if res := storage.GormDB.Debug().Joins("Submission").Joins("ProblemVersion").Preload("ProblemVersion.OutputValidator").First(&run, runId); res.Error != nil {
		return fmt.Errorf("failed loading run: %v", res.Error)
	}
	logger.Infof("Found run %d of submission %d", run.SubmissionRunId, run.SubmissionId)

	run.Status = storage.StatusCompiling
	if res := storage.GormDB.Select("Status").Save(&run); res.Error != nil {
		logger.Warningf("failed marking run as compiling: %v", res.Error)
	}
	program := &apipb.Program{
		Language: langMap[run.Submission.Language],
	}
	submissionFiles := submissionJson{}
	err := json.Unmarshal(run.Submission.SubmissionFiles, &submissionFiles)
	if err != nil {
		return err
	}
	logger.Infof("Files: %v", submissionFiles)
	for path, content := range submissionFiles.Files {
		content, err := base64.StdEncoding.DecodeString(content)
		if err != nil {
			return err
		}
		program.Sources = append(program.Sources, &apipb.SourceFile{
			Path:     filepath.Base(path),
			Contents: content,
		})
	}
	// In case we retry judging of the run, put it in a new folder instead to aovid collisions
	subRoot := fmt.Sprintf("/var/lib/omogen/submissions/%d-%d", runId, time.Now().Unix())
	compile, err := eval.Compile(program, filepath.Join(subRoot, "compile"))
	if err != nil {
		return err
	}
	if compile.Program == nil {
		run.CompileError = compile.CompilerErrors
		run.Status = storage.StatusCompileError
		if res := storage.GormDB.Select("CompileError", "Status").Save(&run); res.Error != nil {
			return fmt.Errorf("failed marking program as compile error: %v", res.Error)
		}
		return nil
	} else {
		run.Status = storage.StatusRunning
		if res := storage.GormDB.Select("Status").Save(&run); res.Error != nil {
			return fmt.Errorf("failed marking program as running: %v", res.Error)
		}
	}
	logger.Infof("Compiled program runs with: %v", compile.Program.RunCommand)

	evalPlan, err := makeEvalPlan(compile.Program, run.ProblemVersion)
	if err != nil {
		return fmt.Errorf("failed constructing evaluation plan: %v", err)
	}
	resultChan := make(chan *apipb.Result, 1000)
	var lastRes *apipb.Result
	resultWait := sync.WaitGroup{}
	resultWait.Add(1)
	// TODO: insert group / case runs in database
	go func() {
		for result := range resultChan {
			lastRes = result
		}
		resultWait.Done()
	}()
	evaluator, err := eval.NewEvaluator(subRoot, evalPlan, resultChan)
	if err != nil {
		return fmt.Errorf("failed initializing evaluator: %v", err)
	}
	if err := evaluator.Evaluate(); err != nil {
		return fmt.Errorf("failed evaluation: %v", err)
	}
	resultWait.Wait()
	run.Status = storage.StatusDone
	switch lastRes.Verdict {
	case apipb.Verdict_ACCEPTED:
		run.Verdict = storage.VerdictAccepted
	case apipb.Verdict_TIME_LIMIT_EXCEEDED:
		run.Verdict = storage.VerdictTimeLimitExceeded
	case apipb.Verdict_WRONG_ANSWER:
		run.Verdict = storage.VerdictWrongAnswer
	case apipb.Verdict_RUN_TIME_ERROR:
		run.Verdict = storage.VerdictRuntimeError
	}
	run.TimeUsageMs = int64(lastRes.TimeUsageMs)
	run.Score = lastRes.Score
	if res := storage.GormDB.Select("Status", "Verdict", "TimeUsageMs", "Score").Save(&run); res.Error != nil {
		return fmt.Errorf("failed writing submission results: %v", res.Error)
	}
	return nil
}

type validatorConfig struct {
	RunCommand []string `json:"run_command"`
}

func makeEvalPlan(program *apipb.CompiledProgram, version storage.ProblemVersion) (*apipb.EvaluationPlan, error) {
	var groups []storage.ProblemTestgroup
	if res := storage.GormDB.Debug().Where("problem_version_id = ?", version.ProblemVersionId).Preload("ProblemTestcases").Order("problem_testgroup_id asc").Find(&groups)
		res.Error != nil {
		return nil, fmt.Errorf("failed gathering testdata: %v", res.Error)
	}

	evalPlan := &apipb.EvaluationPlan{
		Program:              program,
		TimeLimitMs:          int32(version.TimeLimitMs),
		MemLimitKb:           int32(version.MemoryLimitKb),
		ValidatorTimeLimitMs: 60_000,
		ValidatorMemLimitKb:  1_000_000,
	}
	if version.Interactive {
		evalPlan.PlanType = apipb.EvaluationType_INTERACTIVE
	} else {
		evalPlan.PlanType = apipb.EvaluationType_SIMPLE
	}
	if version.OutputValidatorId != 0 {
		evalPlan.ScoringValidator = version.OutputValidator.Scoringvalidator
		conf := validatorConfig{}
		if err := json.Unmarshal(version.OutputValidator.ValidatorRunConfig, &conf); err != nil {
			return nil, err
		}
		val, err := zipProgram(version.OutputValidator.ValidatorSourceZipId, conf.RunCommand)
		if err != nil {
			return nil, fmt.Errorf("failed loading zip'ed validator: %v", err)
		}
		evalPlan.Validator = val
	}

	apigroups := make(map[int64]*apipb.TestGroup)
	for _, group := range groups {
		apigroup, err := toGroup(group)
		if err != nil {
			return nil, fmt.Errorf("failed loading test group %d: %v", group.TestgroupName, err)
		}
		apigroups[group.ProblemTestgroupId] = apigroup
		if group.ParentId != 0 {
			apigroups[group.ParentId].Groups = append(apigroups[group.ParentId].Groups, apigroup)
		}
	}
	evalPlan.RootGroup = apigroups[version.RootGroupId]
	return evalPlan, nil
}

func zipProgram(id string, runCmd []string) (*apipb.CompiledProgram, error) {
	logger.Infof("Loading validator %s", id)
	valPath := filepath.Join("/var/lib/omogen/validator/", id)
	if _, err := os.Stat(valPath); err != nil {
		if os.IsNotExist(err) {
			if err := syncFiles([]string{id}); err != nil {
				return nil, err
			}
			zipPath, _ := findPath(id)
			r, err := zip.OpenReader(zipPath)
			if err != nil {
				return nil, err
			}
			defer r.Close()

			fb := util.NewFileBase(valPath)
			fb.OwnerGid = util.OmogenexecGroupId()
			if err := fb.Mkdir("."); err != nil {
				return nil, err
			}
			for _, f := range r.File {
				subPath := f.Name
				if f.FileInfo().IsDir() {
					if err := fb.Mkdir(subPath); err != nil {
						return nil, err
					}
					continue
				}
				if err := fb.Mkdir(filepath.Dir(subPath)); err != nil {
					return nil, err
				}
				fileReader, err := f.Open()
				if err != nil {
					return nil, err
				}
				content, err := ioutil.ReadAll(fileReader)
				if err != nil {
					return nil, err
				}
				if err := fb.WriteFile(subPath, content); err != nil {
					return nil, err
				}
				if err := fb.FixModeExec(subPath); err != nil {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	}
	return &apipb.CompiledProgram{
		ProgramRoot: valPath,
		RunCommand:  runCmd,
	}, nil
}

func toGroup(testgroup storage.ProblemTestgroup) (*apipb.TestGroup, error) {
	group := &apipb.TestGroup{
		Name:                 testgroup.TestgroupName,
		AcceptScore:          testgroup.AcceptScore.Float64,
		RejectScore:          testgroup.RejectScore.Float64,
		OutputValidatorFlags: testgroup.OutputValidatorFlags,
		BreakOnFail:          testgroup.BreakOnReject,
		// TODO: support scoring problems
		AcceptIfAnyAccepted: false,
	}
	// TODO: support scoring problems
	group.ScoringMode = apipb.ScoringMode_SUM
	group.VerdictMode = apipb.VerdictMode_FIRST_ERROR

	var missingFiles []string
	for _, testcase := range testgroup.ProblemTestcases {
		inpath, hasin := findPath(testcase.InputFileHash)
		outpath, hasout := findPath(testcase.OutputFileHash)
		group.Cases = append(group.Cases, &apipb.TestCase{
			Name:       testcase.TestcaseName,
			InputPath:  inpath,
			OutputPath: outpath,
		})
		if !hasin {
			missingFiles = append(missingFiles, testcase.InputFileHash)
		}
		if !hasout {
			missingFiles = append(missingFiles, testcase.OutputFileHash)
		}
	}
	if err := syncFiles(missingFiles); err != nil {
		return nil, err
	}
	return group, nil
}

func syncFiles(fileIds []string) error {
	if len(fileIds) == 0 {
		return nil
	}
	var files []storage.StoredFile
	res := storage.GormDB.Find(&files, fileIds)
	if res.Error != nil {
		return fmt.Errorf("failed loading stored files: %v", res.Error)
	}
	fb := util.NewFileBase("/var/lib/omogen/cache")
	fb.OwnerGid = util.OmogenexecGroupId()
	for _, file := range files {
		if err := fb.WriteFile(file.FileHash, file.FileContents); err != nil {
			return err
		}
	}
	return nil
}

func findPath(id string) (string, bool) {
	path := fmt.Sprintf("/var/lib/omogen/cache/%s", id)
	if _, err := os.Stat(path); err == nil {
		return path, true
	}
	return path, false
}
