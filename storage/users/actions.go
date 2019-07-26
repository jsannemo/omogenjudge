// Database actions relating to users.
package users

import (
  "context"
  "errors"
  "database/sql"

  "golang.org/x/crypto/bcrypt"

  "github.com/jsannemo/omogenjudge/storage/db"
)

var ErrUserExists = errors.New("The username is in use")
var ErrInvalidLogin = errors.New("The login details were incorrect")

func hashPassword(password string) (string, error) {
  hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
  if err != nil {
    return "", err
  }
	return string(hash), nil
}

func verifyHashAndPassword(hash, password string) (bool, error) {
  err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
  if err == bcrypt.ErrMismatchedHashAndPassword {
    return false, nil
  } else if err != nil {
    return false, err
  }
  return true, nil
}

// scanAccount reads an account from the given Scannable
// Account columns are given in the order:
// - account.account_id
// - account.username
// - account.password_hash
func scanAccount(sc db.Scannable) (*Account, error) {
  var account Account
  err := sc.Scan(&account.AccountId, &account.Username, &account.PasswordHash)
  return &account, err
}

// AuthenticateUser verifies if the given username and password corresponds to an existing account.
// If not, an ErrInvalidLogin is returned as error.
func AuthenticateUser(ctx context.Context, username, password string) (*Account, error) {
  conn := db.GetPool()
  row := conn.QueryRow("SELECT account_id, username, password_hash FROM account WHERE username = $1", username)
  user, err := scanAccount(row)
  if err == sql.ErrNoRows {
    return nil, ErrInvalidLogin
  } else if err != nil {
    return nil, err
  }
  match, err := verifyHashAndPassword(user.PasswordHash, password)
  if err != nil {
    return nil, err
  }
  if !match {
    return nil, ErrInvalidLogin
  }
  return user, nil
}

// CreateUser creates a new user with the given username and password, returning the account ID of the new user.
func CreateUser(ctx context.Context, username, password string) (int32, error) {
  passwordHash, err := hashPassword(password)
  if err != nil {
    return 0, err
  }
  conn := db.GetPool()
  var id int32
  err = conn.QueryRow("INSERT INTO account(username, password_hash) VALUES($1, $2) RETURNING account_id", username, passwordHash).Scan(&id)
  if err != nil {
    if db.PgErrCode(err) == db.UniquenessViolation {
      return 0, ErrUserExists
    } else {
      return 0, err
    }
  }
  return id, nil
}
