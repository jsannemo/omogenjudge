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
	"ruby":    apipb.LanguageGroup_RUBY,
	"rust":    apipb.LanguageGroup_RUST,
	"java":    apipb.LanguageGroup_JAVA,
	"csharp":  apipb.LanguageGroup_CSHARP,
}

type submissionJson struct {
	Files map[string]string
}

type includedCodeJson struct {
	FilesByLanguage map[string]map[string]string `json:"files_by_language"`
}

var evalMutex sync.Mutex

func toStorageVerdict(verdict apipb.Verdict) storage.Verdict {
	switch verdict {
	case apipb.Verdict_ACCEPTED:
		return storage.VerdictAccepted
	case apipb.Verdict_TIME_LIMIT_EXCEEDED:
		return storage.VerdictTimeLimitExceeded
	case apipb.Verdict_WRONG_ANSWER:
		return storage.VerdictWrongAnswer
	case apipb.Verdict_RUN_TIME_ERROR:
		return storage.VerdictRuntimeError
	}
	panic(fmt.Sprintf("unknown API verdict: %v", verdict))
}

func evaluate(runId int64) error {
	evalMutex.Lock()
	defer evalMutex.Unlock()

	var run storage.SubmissionRun
	if res := storage.GormDB.Debug().Joins("Submission").Joins("ProblemVersion").Preload("ProblemVersion.OutputValidator").Preload("ProblemVersion.CustomGrader").First(&run, runId); res.Error != nil {
		return fmt.Errorf("failed loading run: %v", res.Error)
	}
	logger.Infof("Found run %d of submission %d", run.SubmissionRunId, run.SubmissionId)

	run.Status = storage.StatusCompiling
	if res := storage.GormDB.Select("Status").Save(&run); res.Error != nil {
		logger.Warningf("failed marking run as compiling: %v", res.Error)
	}
	lang, ok := langMap[run.Submission.Language]
	if !ok {
		return fmt.Errorf("run has unknown language %s", run.Submission.Language)
	}
	program := &apipb.Program{
		Language: lang,
	}
	submissionFiles := submissionJson{}
	if err := json.Unmarshal(run.Submission.SubmissionFiles, &submissionFiles); err != nil {
		return err
	}
	includedCode := includedCodeJson{}
	if err := json.Unmarshal(run.ProblemVersion.IncludedFiles, &includedCode); err != nil {
		return err
	}
	logger.Infof("Lang: %v, included code:", run.Submission.Language, includedCode)
	extraFiles := includedCode.FilesByLanguage[run.Submission.Language]

	logger.Infof("Files: %v", submissionFiles)
	for path, content := range submissionFiles.Files {
		if _, hasExtraFile := extraFiles[path]; hasExtraFile {
			continue
		}
		content, err := base64.StdEncoding.DecodeString(content)
		if err != nil {
			return err
		}
		program.Sources = append(program.Sources, &apipb.SourceFile{
			Path:     filepath.Base(path),
			Contents: content,
		})
	}

	logger.Infof("Extra files: %v", extraFiles)
	for path, content := range extraFiles {
		program.Sources = append(program.Sources, &apipb.SourceFile{
			Path:     filepath.Base(path),
			Contents: []byte(content),
		})
	}

	// In case we retry judging of the run, put it in a new folder instead to avoid collisions
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

	var groups []*storage.ProblemTestgroup
	if res := storage.GormDB.Debug().Where("problem_version_id = ?", run.ProblemVersionId).Preload("ProblemTestcases").Order("problem_testgroup_id asc").Find(&groups); res.Error != nil {
		return fmt.Errorf("failed gathering testdata: %v", res.Error)
	}
	subgroupMap := make(map[int64][]*storage.ProblemTestgroup)
	for _, group := range groups {
		if group.ProblemTestgroupId == run.ProblemVersion.RootGroupId {
			run.ProblemVersion.RootGroup = group
		}
		if group.ParentId != 0 {
			subgroupMap[group.ParentId] = append(subgroupMap[group.ParentId], group)
		}
	}

	evalPlan, err := makeEvalPlan(compile.Program, run.ProblemVersion)
	if err != nil {
		return fmt.Errorf("failed constructing evaluation plan: %v", err)
	}
	resultChan := make(chan *apipb.Result, 1000)
	var lastRes *apipb.Result
	resultWait := sync.WaitGroup{}
	resultWait.Add(1)
	var resultError error
	go func() {
		var groupStack []*storage.ProblemTestgroup
		var tcIdx []int
		var groupIdx []int
		groupStack = append(groupStack, run.ProblemVersion.RootGroup)
		tcIdx = append(tcIdx, 0)
		groupIdx = append(groupIdx, 0)

		for result := range resultChan {
			for {
				curIdx := len(groupStack) - 1
				curGroup := groupStack[curIdx]
				subgroups, found := subgroupMap[curGroup.ProblemTestgroupId]
				if !found {
					break
				}
				nextTc := tcIdx[curIdx]
				nextGroup := groupIdx[curIdx]
				// If the next item to be judged is a test case group, transcend down to that group
				if nextGroup < len(subgroups) && (nextTc == len(curGroup.ProblemTestcases) ||
					subgroups[nextGroup].TestgroupName < curGroup.ProblemTestcases[nextTc].TestcaseName) {
					groupIdx[curIdx] = nextGroup + 1
					groupStack = append(groupStack, subgroups[nextGroup])
					tcIdx = append(tcIdx, 0)
					groupIdx = append(groupIdx, 0)
				} else {
					break
				}
			}

			curIdx := len(groupStack) - 1
			curGroup := groupStack[curIdx]
			switch result.Type {
			case apipb.ResultType_TEST_CASE:
				testcase := curGroup.ProblemTestcases[tcIdx[curIdx]]
				tcRun := storage.SubmissionCaseRun{
					SubmissionRunId:   run.SubmissionRunId,
					ProblemTestcaseId: testcase.ProblemTestcaseId,
					TimeUsageMs:       result.TimeUsageMs,
					Score:             result.Score,
					Verdict:           toStorageVerdict(result.Verdict),
				}
				if res := storage.GormDB.Save(&tcRun); res.Error != nil {
					resultError = res.Error
				}
				tcIdx[curIdx] += 1
			case apipb.ResultType_TEST_GROUP:
				tcRun := storage.SubmissionGroupRun{
					SubmissionRunId:    run.SubmissionRunId,
					ProblemTestgroupId: curGroup.ProblemTestgroupId,
					TimeUsageMs:        result.TimeUsageMs,
					Score:              result.Score,
					Verdict:            toStorageVerdict(result.Verdict),
				}
				if res := storage.GormDB.Save(&tcRun); res.Error != nil {
					resultError = res.Error
				}
				groupStack = groupStack[:curIdx]
				tcIdx = tcIdx[:curIdx]
				groupIdx = groupIdx[:curIdx]
			}
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
	if resultError != nil {
		return fmt.Errorf("failed writing sub-submission results: %v", err)
	}
	run.Status = storage.StatusDone
	run.TimeUsageMs = lastRes.TimeUsageMs
	run.Score = lastRes.Score
	run.Verdict = toStorageVerdict(lastRes.Verdict)
	if res := storage.GormDB.Select("Status", "Verdict", "TimeUsageMs", "Score").Save(&run); res.Error != nil {
		return fmt.Errorf("failed writing submission results: %v", res.Error)
	}
	return nil
}

type validatorConfig struct {
	RunCommand []string `json:"run_command"`
}

type graderConfig struct {
	RunCommand []string `json:"run_command"`
}

func makeEvalPlan(program *apipb.CompiledProgram, version storage.ProblemVersion) (*apipb.EvaluationPlan, error) {
	var groups []storage.ProblemTestgroup
	if res := storage.GormDB.Debug().Where("problem_version_id = ?", version.ProblemVersionId).Preload("ProblemTestcases").Order("problem_testgroup_id asc").Find(&groups); res.Error != nil {
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
		evalPlan.ScoringValidator = version.OutputValidator.ScoringValidator
		val, err := zipProgram(version.OutputValidator.ValidatorZipId, version.OutputValidator.RunCommand, "validators")
		if err != nil {
			return nil, fmt.Errorf("failed loading zip'ed validator: %v", err)
		}
		evalPlan.Validator = val
	}
	if version.CustomGraderId != 0 {
		grader, err := zipProgram(version.CustomGrader.GraderZipId, version.CustomGrader.RunCommand, "graders")
		if err != nil {
			return nil, fmt.Errorf("failed loading zip'ed grader: %v", err)
		}
		evalPlan.Grader = grader
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

func zipProgram(id string, runCmd []string, programType string) (*apipb.CompiledProgram, error) {
	logger.Infof("Loading validator %s", id)
	valPath := filepath.Join("/var/lib/omogen/", programType, id)
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

func toApiScoringMode(scoringMode string) (apipb.ScoringMode, error) {
	switch scoringMode {
	case storage.ScoringModeAvg:
		return apipb.ScoringMode_AVG, nil
	case storage.ScoringModeMax:
		return apipb.ScoringMode_MAX, nil
	case storage.ScoringModeMin:
		return apipb.ScoringMode_MIN, nil
	case storage.ScoringModeSum:
		return apipb.ScoringMode_SUM, nil
	}
	return apipb.ScoringMode_SCORING_MODE_UNSPECIFIED, fmt.Errorf("unknown scoring mode: %s", scoringMode)
}

func toApiVerdictMode(verdictMode string) (apipb.VerdictMode, error) {
	switch verdictMode {
	case storage.VerdictModeWorstError:
		return apipb.VerdictMode_WORST_ERROR, nil
	case storage.VerdictModeFirstError:
		return apipb.VerdictMode_FIRST_ERROR, nil
	case storage.VerdictModeAlwaysAccept:
		return apipb.VerdictMode_ALWAYS_ACCEPT, nil
	}
	return apipb.VerdictMode_VERDICT_MODE_UNSPECIFIED, fmt.Errorf("unknown verdict mode: %s", verdictMode)
}

func toGroup(testgroup storage.ProblemTestgroup) (*apipb.TestGroup, error) {
	scoringMode, err := toApiScoringMode(testgroup.ScoringMode)
	if err != nil {
		return nil, err
	}
	verdictMode, err := toApiVerdictMode(testgroup.VerdictMode)
	if err != nil {
		return nil, err
	}
	group := &apipb.TestGroup{
		Name:                 testgroup.TestgroupName,
		AcceptScore:          testgroup.AcceptScore.Float64,
		RejectScore:          testgroup.RejectScore.Float64,
		OutputValidatorFlags: testgroup.OutputValidatorFlags,
		BreakOnFail:          testgroup.BreakOnReject,
		AcceptIfAnyAccepted:  testgroup.AcceptIfAnyAccepted,
		IgnoreSample:         testgroup.IgnoreSample,
		ScoringMode:          scoringMode,
		VerdictMode:          verdictMode,
		CustomGrading:        testgroup.CustomGrading,
		GraderFlags:          testgroup.GraderFlags,
	}

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
