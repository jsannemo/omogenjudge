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
