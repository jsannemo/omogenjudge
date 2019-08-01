package users

import (
  "context"
  "errors"

  "github.com/jsannemo/omogenjudge/storage/db"
  "github.com/jsannemo/omogenjudge/storage/models"
)

var ErrUserExists = errors.New("The username is in use")

func Create(ctx context.Context, account *models.Account) error {
  conn := db.Conn()
  if err := conn.QueryRowContext(ctx, "INSERT INTO account(username, password_hash) VALUES($1, $2) RETURNING account_id", account.Username, account.PasswordHash).Scan(&account.AccountId); err != nil {
    if db.PgErrCode(err) == db.UniquenessViolation {
      return ErrUserExists
    } else {
      panic(err)
    }
  }
  return nil
}
