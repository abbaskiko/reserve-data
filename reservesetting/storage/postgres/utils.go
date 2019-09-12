package postgres

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

// TODO: replace this with common/postgres/transaction.go:10
func rollbackUnlessCommitted(tx *sqlx.Tx) {
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		log.Printf("failed to roll back transaction")
	}
}
