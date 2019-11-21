package problems

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

var ErrDuplicateProblemName = errors.New("Duplicate problem name")

func insertTest(ctx context.Context, tc *models.TestCase, tx *sqlx.Tx) error {
	return tx.QueryRowContext(
		ctx,
		`
        INSERT INTO
            problem_testcase(problem_testgroup_id, testcase_name, input_file_hash.hash, output_file_hash.hash)
        VALUES($1, $2, $3, $4)
        RETURNING problem_testcase_id`,
		tc.TestGroupId, tc.Name, tc.InputFile.Hash, tc.OutputFile.Hash).Scan(&tc.TestCaseId)
}

func insertTestGroup(ctx context.Context, tg *models.TestGroup, tx *sqlx.Tx) error {
	if err := tx.QueryRowContext(
		ctx,
		`
        INSERT INTO
          problem_testgroup(problem_version_id, testgroup_name, public_visibility, score, output_validator_flags)
        VALUES($1, $2, $3, $4, $5)
        RETURNING problem_testgroup_id`,
		tg.ProblemVersionId, tg.Name, tg.PublicVisibility, tg.Score, tg.OutputValidatorFlags).Scan(&tg.TestGroupId); err != nil {
		return err
	}
	for _, tc := range tg.Tests {
		tc.TestGroupId = tg.TestGroupId
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
		s.ProblemId, s.Language, s.Title, s.Html); err != nil {
		return err
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
      problem_output_validator(problem_id, validator_language_id, validator_source_zip.hash)
    VALUES($1, $2, $3)`, version.ProblemId, version.OutputValidator.ValidatorLanguageId, version.OutputValidator.ValidatorSourceZip.Hash); err != nil {
		return err
	}
	return nil
}

func insertProblem(ctx context.Context, problem *models.Problem, tx *sqlx.Tx) error {
	return tx.QueryRowContext(ctx,
		`
    INSERT INTO
      problem(short_name, author, license)
    VALUES($1, $2, $3)
    RETURNING problem_id`,
		problem.ShortName, problem.Author, problem.License).Scan(&problem.ProblemId)
}

func insertProblemVersion(ctx context.Context, problem *models.Problem, tx *sqlx.Tx) error {
	return tx.QueryRowContext(ctx,
		`
    INSERT INTO
        problem_version(problem_id, time_limit_ms, memory_limit_kb)
    VALUES($1, $2, $3)
    RETURNING problem_version_id`,
		problem.ProblemId,
		problem.CurrentVersion.TimeLimMs,
		problem.CurrentVersion.MemLimKb).Scan(&problem.CurrentVersion.ProblemVersionId)
}

// TODO(jsannemo): add support for updating
func Create(ctx context.Context, p *models.Problem) error {
	err := db.InTransaction(ctx, func(tx *sqlx.Tx) error {
		if err := insertProblem(ctx, p, tx); err != nil {
			if db.PgErrCode(err) == db.UniquenessViolation {
				return ErrDuplicateProblemName
			} else {
				return err
			}
		}

		if err := insertProblemVersion(ctx, p, tx); err != nil {
			return err
		}

		if err := insertOutputValidator(ctx, p.CurrentVersion, tx); err != nil {
			return err
		}
		for _, s := range p.Statements {
			s.ProblemId = p.ProblemId
			if err := insertStatement(ctx, s, tx); err != nil {
				return err
			}
		}
		for _, tg := range p.CurrentVersion.TestGroups {
			tg.ProblemVersionId = p.CurrentVersion.ProblemVersionId
			if err := insertTestGroup(ctx, tg, tx); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
