package models

import (
	"fmt"

	"github.com/google/logger"
	"golang.org/x/crypto/bcrypt"

	"github.com/jsannemo/omogenjudge/frontend/paths"
)

type AccountList []*Account

func (accs AccountList) Single() *Account {
	if len(accs) == 0 {
		panic(fmt.Errorf("Request single account from empty list"))
	}
	if len(accs) > 1 {
		panic(fmt.Errorf("Request single account from multi-entry list"))
	}
	return accs[0]
}

type Account struct {
	// This should not be exposed externally.
	AccountId int32 `db:"account_id"`

	// This can be exposed externally and used in URLs.
	Username string `db:"username"`

	// A SHA-256 hash of the user's password.
	PasswordHash string `db:"password_hash"`
}

func (a *Account) SetPassword(password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		logger.Fatalf("Could not hash password: %v", hash)
	}
	a.PasswordHash = string(hash)
}

func (a Account) MatchesPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false
	} else if err != nil {
		logger.Fatalf("Could not match password: %v", err)
	}
	return true
}

func (a Account) Link() string {
	return paths.Route(paths.User, paths.UserNameArg, a.Username)
}
