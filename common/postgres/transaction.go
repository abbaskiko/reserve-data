package postgres

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

func RollbackUnlessCommitted(tx *sqlx.Tx) {
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		log.Printf("failed to roll back transaction")
	}
}
