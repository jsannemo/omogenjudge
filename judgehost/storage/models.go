package storage

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"time"
)

type JSON json.RawMessage

type StoredFile struct {
	FileHash     string `gorm:"primaryKey"`
	FileContents []byte `gorm:"type:bytea"`
}

type Problem struct {
	ProblemId        int64 `gorm:"primaryKey"`
	CurrentVersionId int64
	CurrentVersion   ProblemVersion `gorm:"foreignKey:CurrentVersionId; References:ProblemVersionId"`
}

type ProblemOutputValidator struct {
	ProblemOutputValidatorId int64          `gorm:"primaryKey"`
	RunCommand               pq.StringArray `gorm:"type:text[]"`
	ValidatorZipId           string
	ValidatorZip             StoredFile `gorm:"foreignKey:ValidatorZipId; References:FileHash"`
	ScoringValidator         bool
}

type ProblemGrader struct {
	ProblemGraderId int64          `gorm:"primaryKey"`
	RunCommand      pq.StringArray `gorm:"type:text[]"`
	GraderZipId     string
	GraderZip       StoredFile `gorm:"foreignKey:GraderZipId; References:FileHash"`
}

type ProblemTestcase struct {
	ProblemTestcaseId  int64 `gorm:"primaryKey"`
	ProblemTestgroupId int64
	TestcaseName       string
	InputFileHash      string
	InputFile          StoredFile `gorm:"foreignKey:InputFileHash; References:FileHash"`
	OutputFileHash     string
	OutputFile         StoredFile `gorm:"foreignKey:OutputFileHash; References:FileHash"`
}

const (
	VerdictModeWorstError   = "worst_error"
	VerdictModeFirstError   = "first_error"
	VerdictModeAlwaysAccept = "always_accept"
)

const (
	ScoringModeSum = "sum"
	ScoringModeAvg = "avg"
	ScoringModeMin = "min"
	ScoringModeMax = "max"
)

type ProblemTestgroup struct {
	ProblemTestgroupId   int64 `gorm:"primaryKey"`
	ParentId             int64
	Parent               *ProblemTestgroup `gorm:"foreignKey:ParentId; References:ProblemTestgroupId"`
	ProblemVersionId     int64
	ProblemVersion       ProblemVersion `gorm:"foreignKey:ProblemVersionId; References:ProblemVersionId"`
	TestgroupName        string
	MinScore             sql.NullFloat64
	MaxScore             sql.NullFloat64
	AcceptScore          sql.NullFloat64
	RejectScore          sql.NullFloat64
	BreakOnReject        bool
	AcceptIfAnyAccepted  bool
	IgnoreSample         bool
	ScoringMode          string
	VerdictMode          string
	GraderFlags          pq.StringArray `gorm:"type:text[]"`
	CustomGrading        bool
	OutputValidatorFlags pq.StringArray    `gorm:"type:text[]"`
	ProblemTestcases     []ProblemTestcase `gorm:"References:ProblemTestgroupId"`
}

type ProblemVersion struct {
	ProblemVersionId  int64 `gorm:"primaryKey"`
	ProblemId         int64
	RootGroupId       int64
	RootGroup         *ProblemTestgroup `gorm:"foreignKey:RootGroupId; References:ProblemTestgroupId"`
	TimeLimitMs       int64
	MemoryLimitKb     int64
	OutputValidatorId int64
	OutputValidator   ProblemOutputValidator
	CustomGraderId    int64
	CustomGrader      ProblemGrader
	IncludedFiles     JSON
	Scoring           bool
	Interactive       bool
	ScoreMaximization sql.NullBool
}

type SubmissionCaseRun struct {
	SubmissionCaseRunId int64 `gorm:"primaryKey"`
	SubmissionRunId     int64
	SubmissionRun       SubmissionRun
	ProblemTestcaseId   int64
	ProblemTestcase     ProblemTestcase
	DateCreated         time.Time `gorm:"autoCreateTime"`
	TimeUsageMs         int64
	Score               float64
	Verdict             Verdict
}

type SubmissionGroupRun struct {
	SubmissionGroupRunId int64 `gorm:"primaryKey"`
	SubmissionRunId      int64
	SubmissionRun        SubmissionRun
	ProblemTestgroupId   int64
	ProblemTestgroup     ProblemTestgroup
	DateCreated          time.Time `gorm:"autoCreateTime"`
	TimeUsageMs          int64
	Score                float64
	Verdict              Verdict
}

const (
	StatusQueued       = "queued"
	StatusCompiling    = "compiling"
	StatusRunning      = "running"
	StatusCompileError = "compile error"
	StatusJudgeError   = "judging error"
	StatusDone         = "done"
)

type Verdict string

const (
	VerdictUnjudged          Verdict = "unjudged"
	VerdictAccepted          Verdict = "accepted"
	VerdictWrongAnswer       Verdict = "wrong answer"
	VerdictTimeLimitExceeded Verdict = "time limit exceeded"
	VerdictRuntimeError      Verdict = "run-time error"
)

type Submission struct {
	SubmissionId    int64 `gorm:"primaryKey"`
	Language        string
	SubmissionFiles JSON
}

type SubmissionRun struct {
	SubmissionRunId  int64 `gorm:"primaryKey"`
	SubmissionId     int64
	Submission       Submission `gorm:"foreignkey:SubmissionId; References:SubmissionId"`
	ProblemVersionId int64
	ProblemVersion   ProblemVersion `gorm:"foreignKey:ProblemVersionId; References:ProblemVersionId"`
	DateCreated      time.Time      `gorm:"autoCreateTime"`
	Status           string
	Verdict          Verdict
	TimeUsageMs      int64
	Score            float64
	CompileError     string
}

func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

func (j *JSON) Value() (driver.Value, error) {
	if len(*j) == 0 {
		return nil, nil
	}
	return json.RawMessage(*j).MarshalJSON()
}
