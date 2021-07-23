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

type StoredFile struct {
	FileHash     string `gorm:"primaryKey"`
	FileContents []byte `gorm:"type:bytea"`
}

type Problem struct {
	ProblemId int64 `gorm:"primaryKey"`
	ShortName string
	Author pq.StringArray `gorm:"type:text[]"`
	Source string
	License string
	CurrentVersionId int64
	CurrentVersion ProblemVersion `gorm:"foreignKey:CurrentVersionId; References:ProblemVersionId"`
	ProblemStatements []ProblemStatement `gorm:"References:ProblemId"`

}

type ProblemStatement struct {
	Id int64 `gorm:"primaryKey"`
	ProblemId int64
	Problem  Problem `gorm:"foreignKey:ProblemId; References:ProblemId"`
	Language string
	Title string
	Html string
}

type ProblemStatementFile struct {
	Id int64 `gorm:"primaryKey"`
	ProblemId int64
	Problem  Problem `gorm:"foreignKey:ProblemId; References:Problem"`
	FilePath string
	StatementFileHash string
	StatementFile StoredFile `gorm:"foreignKey:StatementFileHash; References:StoredFile"`
	Attachment bool
}

const (
	LicensePublicDomain = "public domain"
	LicenseCC0 = "cc0"
	LicenseCCBY = "cc by"
	LicenseCCBYSA = "cc by-sa"
	LicenseEducational = "educational"
	LicensePermission = "permission"
)

type ProblemOutputValidator struct {
	Id                   int64 `gorm:"primaryKey"`
	ValidatorRunConfig   JSON
	ValidatorSourceZipId string
	ValidatorSourceZip   StoredFile `gorm:"foreignKey:ValidatorSourceZipId; References:FileHash"`
	ScoringValidator     bool
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
	IncludedFiles     JSON
	Scoring           bool
	Interactive       bool
	ScoreMaximization sql.NullBool
}

type SubmissionCaseRun struct {
	Id                int64 `gorm:"primaryKey"`
	SubmissionRunId   int64
	SubmissionRun     SubmissionRun
	ProblemTestcaseId int64
	ProblemTestcase   ProblemTestcase
	DateCreated       time.Time `gorm:"autoCreateTime"`
	TimeUsageMs       int64
	Score             float64
	Verdict           string
}

type SubmissionGroupRun struct {
	Id                 int64 `gorm:"primaryKey"`
	SubmissionRunId    int64
	SubmissionRun      SubmissionRun
	ProblemTestgroupId int64
	ProblemTestgroup   ProblemTestgroup
	DateCreated        time.Time `gorm:"autoCreateTime"`
	TimeUsageMs        int64
	Score              float64
	Verdict            string
}

const (
	StatusQueued       = "queued"
	StatusCompiling    = "compiling"
	StatusRunning      = "running"
	StatusCompileError = "compile error"
	StatusJudgeError   = "judging error"
	StatusDone         = "done"
)

const (
	VerdictUnjudged          = "unjudged"
	VerdictAccepted          = "accepted"
	VerdictWrongAnswer       = "wrong answer"
	VerdictTimeLimitExceeded = "time limit exceeded"
	VerdictRuntimeError      = "run-time error"
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
	Verdict          string
	TimeUsageMs      int64
	Score            float64
	CompileError     string
}

type JSON json.RawMessage

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

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}
