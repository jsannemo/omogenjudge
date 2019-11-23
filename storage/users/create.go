package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

// ErrUserExists is returned when the given username already exists.
var ErrUserExists = errors.New("the username is in use")

// CreateUser persists a new user account in the database.
func CreateUser(ctx context.Context, account *models.Account) error {
	query := `
    INSERT INTO
      account(username, password_hash, email, full_name)
    VALUES($1, $2, $3, $4)
    RETURNING account_id`
	conn := db.Conn()
	if err := conn.QueryRowContext(ctx, query, account.Username, account.PasswordHash, account.Email, account.FullName).Scan(&account.AccountID); err != nil {
		if db.PgErrCode(err) == db.UniquenessViolation {
			return ErrUserExists
		} else {
			return fmt.Errorf("failed create user query: %v", err)
		}
	}
	return nil
}
