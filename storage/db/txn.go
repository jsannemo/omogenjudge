package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func InTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx := Conn().MustBeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	err := fn(tx)
	if err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("failed rollback after (%v) due to (%v)", err, rerr)
		}
		return fmt.Errorf("transaction failed: %v", err)
	} else {
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func SetParam(str string, params []interface{}, param interface{}) (string, []interface{}) {
	str = fmt.Sprintf(str, len(params))
	params = append(params, param)
	return str, params
}

func SetInParam(str string, params []interface{}, param []interface{}) (string, []interface{}) {
	arg := ""
	for _, v := range param {
		if arg == "" {
			arg = fmt.Sprintf("$%d", len(params))
		} else {
			arg = fmt.Sprintf(arg+",$%d", len(params))
		}
		params = append(params, v)
	}
	return fmt.Sprintf(str, arg), params
}

func SetInParamInt(str string, params []interface{}, param []int32) (string, []interface{}) {
	var sl []interface{}
	for _, v := range param {
		sl = append(sl, v)
	}
	return SetInParam(str, params, sl)
}
