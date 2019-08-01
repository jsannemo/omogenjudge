package users

import (
  "context"
  "errors"
  "fmt"
  "strings"

  "github.com/jsannemo/omogenjudge/storage/db"
  "github.com/jsannemo/omogenjudge/storage/models"
)

var ErrNoSuchUser = errors.New("The user ID did not exist")
var ErrInvalidLogin = errors.New("The login details were incorrect")

type ListArgs struct {
}

type FilterArgs struct {
  AccountId int32
  Username string
}

func listQuery(args ListArgs, filter FilterArgs) (string, []interface{}) {
  var filterSegs []string
  var params []interface{}
  if filter.AccountId != 0 {
    params = append(params, filter.AccountId)
    filterSegs = append(filterSegs, fmt.Sprintf("account_id = $%d", len(params)))
  }
  if filter.Username != "" {
    params = append(params, filter.Username)
    filterSegs = append(filterSegs, fmt.Sprintf("username = $%s", len(params)))
  }

  filterStr := ""
  if len(filterSegs) != 0 {
    filterStr = fmt.Sprintf("WHERE %s", strings.Join(filterSegs, " AND "))
  }

  return fmt.Sprintf("SELECT account_id, username, password_hash FROM account %s", filterStr), params
}

func List(ctx context.Context, args ListArgs, filter FilterArgs) models.AccountList {
  conn := db.Conn()
  var accs models.AccountList
  query, params := listQuery(args, filter)
  if err := conn.SelectContext(ctx, &accs, query, params...); err != nil {
    panic(err)
  }
  return accs
}

func Get(ctx context.Context, id int32) (*models.Account, error) {
  accs := List(ctx, ListArgs{}, FilterArgs{AccountId: id})
  if len(accs) == 0 {
    return nil, ErrNoSuchUser
  }
  return accs.Single(), nil
}

func Authenticate(ctx context.Context, username, password string) (*models.Account, error) {
  accs := List(ctx, ListArgs{}, FilterArgs{})
  if len(accs) == 0 {
    return nil, ErrInvalidLogin
  }
  acc := accs.Single()
  if match := acc.MatchesPassword(password); !match {
    return nil, ErrInvalidLogin
  }
  return acc, nil
}

