package submissions

import (
	"context"
	"fmt"
	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
	"strings"
)

// A RunField represents a field in a run.
type RunField string

const (
	RunFieldVerdict      RunField = "verdict"
	RunFieldTimeUsageMs  RunField = "time_usage_ms"
	RunFieldScore        RunField = "score"
	RunFieldStatus       RunField = "status"
	RunFieldCompileError RunField = "compile_error"
)

// An UpdateRunArgs controls what to update in a run.
type UpdateRunArgs struct {
	// The fields of a SubmissionRun to update.
	Fields []RunField
}

// UpdateRun updates the given submission run in the database.
func UpdateRun(ctx context.Context, run *models.SubmissionRun, args UpdateRunArgs) error {
	if len(args.Fields) == 0 {
		return nil
	}
	conn := db.Conn()
	var params []interface{}
	var updates []string
	for _, field := range args.Fields {
		switch field {
		case RunFieldVerdict:
			updates = append(updates, db.SetParam("verdict = $%d", &params, run.Verdict))
		case RunFieldTimeUsageMs:
			updates = append(updates, db.SetParam("time_usage_ms = $%d", &params, run.TimeUsageMS))
		case RunFieldScore:
			updates = append(updates, db.SetParam("score = $%d", &params, run.Score))
		case RunFieldStatus:
			updates = append(updates, db.SetParam("status = $%d", &params, run.Status))
		case RunFieldCompileError:
			updates = append(updates, db.SetParam("compile_error = $%d", &params, run.CompileError))
		}
	}
	query := `
		UPDATE submission_run ` +
		fmt.Sprintf(` SET %s `, strings.Join(updates, ", ")+
			` WHERE `+db.SetParam(` submission_run_id = $%d `, &params, run.SubmissionRunID))
	if _, err := conn.ExecContext(ctx, query, params...); err != nil {
		return fmt.Errorf("failed writing submission run update: %v", err)
	}
	return nil
}
