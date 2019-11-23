package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/logger"
	"github.com/jsannemo/omogenjudge/storage/db"
	"github.com/jsannemo/omogenjudge/storage/models"
)

// ErrInvalidLogin is returned when the given login details were incorrect.
var ErrInvalidLogin = errors.New("incorrect login details")

// A ListArgs controls the behaviour of ListUsers.
type ListArgs struct {
}

// A ListFilter controls the filtering behaviour of ListUsers. At most a single filter variable may be set.
type ListFilter struct {
	// Filter by a list of account IDs.
	AccountID []int32
	// Filter by a username.
	Username string
}

// ListUsers returns a list of user accounts.
func ListUsers(ctx context.Context, args ListArgs, filter ListFilter) (AccountList, error) {
	if len(filter.AccountID) != 0 && filter.Username != "" {
		return nil, fmt.Errorf("can only use at most one filter")
	}
	conn := db.Conn()
	var accs AccountList
	query, params := listQuery(args, filter)
	if err := conn.SelectContext(ctx, &accs, query, params...); err != nil {
		return nil, fmt.Errorf("failed user list query: %v", err)
	}
	return accs, nil
}

func listQuery(args ListArgs, filterArgs ListFilter) (string, []interface{}) {
	filter := ""
	var params []interface{}
	if len(filterArgs.AccountID) != 0 {
		filter = db.SetInParamInt("WHERE account_id IN (%s)", &params, filterArgs.AccountID)
	} else if filterArgs.Username != "" {
		filter = db.SetParam("WHERE username = $%d", &params, filterArgs.Username)
	}
	return fmt.Sprintf("SELECT account_id, username, password_hash, email, full_name FROM account %s", filter), params
}

// Authenticate returns a user matching the given login details, if one exists.
func Authenticate(ctx context.Context, username, password string) (*models.Account, error) {
	accs, err := ListUsers(ctx, ListArgs{}, ListFilter{Username: username})
	if err != nil {
		return nil, err
	}
	if len(accs) == 0 {
		return nil, ErrInvalidLogin
	}
	acc := accs.Single()
	if match := acc.MatchesPassword(password); !match {
		return nil, ErrInvalidLogin
	}
	return acc, nil
}

// An AccountList is a slice of user accounts.
type AccountList []*models.Account

// Single returns a single account from an account list.
// If the list contains zero or multiple accounts, it panics.
func (accs AccountList) Single() *models.Account {
	if len(accs) == 0 {
		logger.Fatalf("Requested single account from empty list")
	}
	if len(accs) > 1 {
		logger.Fatalf("Requested single account from multi-entry list")
	}
	return accs[0]
}
