package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	pgutil "github.com/KyberNetwork/reserve-data/common/postgres"
	"github.com/KyberNetwork/reserve-data/reservesetting/common"
	"github.com/KyberNetwork/reserve-data/reservesetting/storage"
)

func (s *Storage) createTradingPair(tx *sqlx.Tx, exchangeID, baseID, quoteID, pricePrecision, amountPrecision uint64,
	amountLimitMin, amountLimitMax, priceLimitMin, priceLimitMax, minNotional float64) (uint64, error) {

	var tradingPairID uint64
	err := tx.NamedStmt(s.stmts.newTradingPair).Get(
		&tradingPairID,
		struct {
			ExchangeID      uint64  `db:"exchange_id"`
			Base            uint64  `db:"base_id"`
			Quote           uint64  `db:"quote_id"`
			PricePrecision  uint64  `db:"price_precision"`
			AmountPrecision uint64  `db:"amount_precision"`
			AmountLimitMin  float64 `db:"amount_limit_min"`
			AmountLimitMax  float64 `db:"amount_limit_max"`
			PriceLimitMin   float64 `db:"price_limit_min"`
			PriceLimitMax   float64 `db:"price_limit_max"`
			MinNotional     float64 `db:"min_notional"`
		}{
			ExchangeID:      exchangeID,
			Base:            baseID,
			Quote:           quoteID,
			PricePrecision:  pricePrecision,
			AmountPrecision: amountPrecision,
			AmountLimitMin:  amountLimitMin,
			AmountLimitMax:  amountLimitMax,
			PriceLimitMin:   priceLimitMin,
			PriceLimitMax:   priceLimitMax,
			MinNotional:     minNotional,
		},
	)
	if err != nil {
		pErr, ok := err.(*pq.Error)
		if !ok {
			return 0, fmt.Errorf("unknown error returned err=%s", err.Error())
		}

		switch pErr.Code {
		case errAssertFailure, errForeignKeyViolation:
			log.Printf("failed to create trading pair as assertion failed base=%d quote=%d exchange_id=%d err=%s",
				baseID,
				quoteID,
				exchangeID,
				pErr.Message)
			return 0, common.ErrBadTradingPairConfiguration
		}

		return 0, fmt.Errorf("failed to create trading pair base=%d quote=%d exchange_id=%d err=%s",
			baseID,
			quoteID,
			exchangeID,
			pErr.Message,
		)
	}
	log.Printf("trading pair created id=%d", tradingPairID)
	return tradingPairID, nil
}

// UpdateTradingPair update a trading pair information
func (s *Storage) UpdateTradingPair(id uint64, updateOpts storage.UpdateTradingPairOpts) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer pgutil.RollbackUnlessCommitted(tx)

	err = s.updateTradingPair(tx, id, updateOpts)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	log.Printf("trading pair update successfully id=%d", id)
	return nil
}

func (s *Storage) updateTradingPair(tx *sqlx.Tx, id uint64, opts storage.UpdateTradingPairOpts) error {
	var updatedID uint64
	err := tx.NamedStmt(s.stmts.updateTradingPair).Get(&updatedID, struct {
		ID              uint64   `db:"id"`
		PricePrecision  *uint64  `db:"price_precision"`
		AmountPrecision *uint64  `db:"amount_precision"`
		AmountLimitMin  *float64 `db:"amount_limit_min"`
		AmountLimitMax  *float64 `db:"amount_limit_max"`
		PriceLimitMin   *float64 `db:"price_limit_min"`
		PriceLimitMax   *float64 `db:"price_limit_max"`
		MinNotional     *float64 `db:"min_notional"`
	}{
		ID:              id,
		PricePrecision:  opts.PricePrecision,
		AmountPrecision: opts.AmountPrecision,
		AmountLimitMin:  opts.AmountLimitMin,
		AmountLimitMax:  opts.AmountLimitMax,
		PriceLimitMin:   opts.PriceLimitMin,
		PriceLimitMax:   opts.PriceLimitMax,
		MinNotional:     opts.MinNotional,
	})
	if err == sql.ErrNoRows {
		return common.ErrNotFound
	} else if err != nil {
		return err
	}
	log.Printf("trading pair configuration %d is updated", id)
	return nil
}

func (s *Storage) deleteTradingPair(tx *sqlx.Tx, id uint64) error {
	var returnedID uint64
	row := tx.Stmt(s.stmts.deleteTradingPair.Stmt).QueryRow(id)
	err := row.Scan(&returnedID)
	switch err {
	case nil:
		return nil
	case sql.ErrNoRows:
		return common.ErrNotFound
	default:
		return err
	}
}
