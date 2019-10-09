package postgres

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func RollbackUnlessCommitted(tx *sqlx.Tx) {
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		l := zap.S()
		l.Errorw("failed to roll back transaction", "err", err)
	}
}
