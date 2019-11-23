package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// InTransaction executes the given function within a transaction with serializable isolation.
// If the function returns an error, the transaction will be rollbacked, otherwise committed.
// Errors from the function will be propagated upwards unchanged.
func InTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx := Conn().MustBeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	err := fn(tx)
	if err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("failed rollback after (%v) due to (%v)", err, rerr)
		}
		return err
	} else if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %v", err)
	}
	return nil
}

// For example, if the 3rd parameter is being set, the string "id = %s" would become "id = $3"
func SetParam(query string, params *[]interface{}, param interface{}) string {
	*params = append(*params, param)
	query = fmt.Sprintf(query, len(*params))
	return query
}

// SetInParam appends a slice to the parameter list and updates the query string with the correct bindvar.
// For example, if first three parameters are being set, the string "id IN (%s)" would become "id IN ($1,$2,$3)"
func SetInParam(str string, params *[]interface{}, param []interface{}) string {
	var args []string
	for _, v := range param {
		*params = append(*params, v)
		args = append(args, fmt.Sprintf("$%d", len(*params)))
	}
	return fmt.Sprintf(str, strings.Join(args, ","))
}

// SetInParamInt calls SetInParam, but first converts the int32 param slice.
func SetInParamInt(str string, params *[]interface{}, param []int32) string {
	var sl []interface{}
	for _, v := range param {
		sl = append(sl, v)
	}
	return SetInParam(str, params, sl)
}
