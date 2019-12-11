package problems

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

// ErrDuplicateProblemName is returned when the inserted problem shortname already exists in the database.
var ErrDuplicateProblemName = errors.New("duplicate problem name")

// CreateProblem persists a new problem into the database.
func CreateProblem(ctx context.Context, p *models.Problem) error {
	err := db.InTransaction(ctx, func(tx *sqlx.Tx) error {
		if err := insertProblem(ctx, p, tx); err != nil {
			if db.PgErrCode(err) == db.UniquenessViolation {
				return ErrDuplicateProblemName
			} else {
				return fmt.Errorf("failed persisting problem: %v", err)
			}
		}

		p.CurrentVersion.ProblemID = p.ProblemID
		if err := insertProblemVersion(ctx, p.CurrentVersion, tx); err != nil {
			return fmt.Errorf("failed persisting problem version: %v", err)
		}
		if err := setCurrentVersion(ctx, p, tx); err != nil {
			return fmt.Errorf("could not set current problem version: %v", err)
		}

		for i, _ := range p.Statements {
			p.Statements[i].ProblemID = p.ProblemID
			if err := insertStatement(ctx, p.Statements[i], tx); err != nil {
				return fmt.Errorf("failed persisting statement: %v", err)
			}
		}
		if err := insertStatementFiles(ctx, p, tx); err != nil {
			return fmt.Errorf("failed persisting statement files: %v", err)
		}
		return nil
	})
	return err
}

func insertTest(ctx context.Context, tc *models.TestCase, tx *sqlx.Tx) error {
	return tx.QueryRowContext(
		ctx,
		`
        INSERT INTO
            problem_testcase(problem_testgroup_id, testcase_name, input_file_hash, output_file_hash)
        VALUES($1, $2, $3, $4)
        RETURNING problem_testcase_id`,
		tc.TestGroupID, tc.Name, tc.InputFile.Hash, tc.OutputFile.Hash).Scan(&tc.TestCaseID)
}

func insertTestGroup(ctx context.Context, tg *models.TestGroup, tx *sqlx.Tx) error {
	if err := tx.QueryRowContext(
		ctx,
		`
        INSERT INTO
          problem_testgroup(problem_version_id, testgroup_name, public_visibility, score, output_validator_flags)
        VALUES($1, $2, $3, $4, $5)
        RETURNING problem_testgroup_id`,
		tg.ProblemVersionID, tg.Name, tg.PublicVisibility, tg.Score, tg.OutputValidatorFlags).Scan(&tg.TestGroupID); err != nil {
		return err
	}
	for _, tc := range tg.Tests {
		tc.TestGroupID = tg.TestGroupID
		if err := insertTest(ctx, tc, tx); err != nil {
			return err
		}
	}
	return nil
}

func insertStatement(ctx context.Context, s *models.ProblemStatement, tx *sqlx.Tx) error {
	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO problem_statement(problem_id, language, title, html) VALUES($1, $2, $3, $4)`,
		s.ProblemID, s.Language, s.Title, s.HTML); err != nil {
		return err
	}
	return nil
}

func insertStatementFiles(ctx context.Context, p *models.Problem, tx *sqlx.Tx) error {
	for _, file := range p.StatementFiles {
		file.ProblemID = p.ProblemID
		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO
    					problem_statement_file(problem_id, file_hash, file_path, attachment)
    				VALUES($1, $2, $3, $4)`,
			p.ProblemID, file.Content.Hash, file.Path, file.Attachment); err != nil {
			return err
		}
	}
	return nil
}

func insertOutputValidator(ctx context.Context, version *models.ProblemVersion, tx *sqlx.Tx) error {
	if version.OutputValidator == nil {
		return nil
	}
	if _, err := tx.ExecContext(ctx,
		`
    INSERT INTO
      problem_output_validator(problem_version_id, validator_language_id, validator_source_zip)
    VALUES($1, $2, $3)`, version.ProblemVersionID, version.OutputValidator.ValidatorLanguageID, version.OutputValidator.ValidatorSourceZIP.Hash); err != nil {
		return err
	}
	return nil
}

func insertProblem(ctx context.Context, problem *models.Problem, tx *sqlx.Tx) error {
	return tx.QueryRowContext(ctx,
		`
    INSERT INTO
      problem(short_name, author, source, license)
    VALUES($1, $2, $3, $4)
    RETURNING problem_id`,
		problem.ShortName, problem.Author, problem.Source, problem.License).Scan(&problem.ProblemID)
}

func setCurrentVersion(ctx context.Context, problem *models.Problem, tx *sqlx.Tx) error {
	_, err := tx.ExecContext(ctx,
		`
	UPDATE problem
	SET current_version = $1
	WHERE problem_id = $2`,
		problem.CurrentVersion.ProblemVersionID, problem.ProblemID)
	return err
}

func insertProblemVersion(ctx context.Context, version *models.ProblemVersion, tx *sqlx.Tx) error {
	if err := tx.QueryRowContext(ctx, `
    INSERT INTO
        problem_version(problem_id, time_limit_ms, memory_limit_kb)
    VALUES($1, $2, $3)
    RETURNING problem_version_id`,
		version.ProblemID,
		version.TimeLimMS,
		version.MemLimKB).Scan(&version.ProblemVersionID); err != nil {
		return err
	}
	if err := insertOutputValidator(ctx, version, tx); err != nil {
		return fmt.Errorf("failed persisting output validator: %v", err)
	}
	for _, tg := range version.TestGroups {
		tg.ProblemVersionID = version.ProblemVersionID
		if err := insertTestGroup(ctx, tg, tx); err != nil {
			return fmt.Errorf("failed persisting test group %v", err)
		}
	}
	return nil
}
