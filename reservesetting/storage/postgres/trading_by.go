package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	pgutil "github.com/KyberNetwork/reserve-data/common/postgres"
	"github.com/KyberNetwork/reserve-data/reservesetting/common"
)

func (s *Storage) createTradingBy(tx *sqlx.Tx, assetID, tradingPairID uint64) (uint64, error) {
	var tradingByID uint64
	err := tx.Stmtx(s.stmts.newTradingBy).Get(&tradingByID, assetID, tradingPairID)
	if err != nil {
		pErr, ok := err.(*pq.Error)
		if !ok {
			return 0, fmt.Errorf("unknown error returned err=%s", err.Error())
		}

		if pErr.Code == errCodeUniqueViolation {
			return 0, common.ErrTradingByAlreadyExists
		}

		return 0, fmt.Errorf("failed to create TradingBy, err=%s", pErr)
	}
	return tradingByID, nil
}

// CreateTradingBy create TradingBy for exists Asset and TradingPair
func (s *Storage) CreateTradingBy(assetID, tradingPairID uint64) (uint64, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer pgutil.RollbackUnlessCommitted(tx)
	id, err := s.createTradingBy(tx, assetID, tradingPairID)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	log.Printf("asset trading by #%d has been created successfully\n", id)
	return id, nil
}

// TODO: create endpoint for this function later
// DeleteTradingBy delete a trading by with a given ID
func (s *Storage) DeleteTradingBy(tradingByID uint64) error {
	var returningTradingByID uint64
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer pgutil.RollbackUnlessCommitted(tx)

	err = tx.Stmtx(s.stmts.deleteTradingBy).Get(&returningTradingByID, tradingByID)
	if err != nil {
		if err == sql.ErrNoRows {
			return common.ErrNotFound
		}
		return err
	}
	log.Printf("asset trading by #%d has been deleted successfully\n", tradingByID)
	return nil
}

// GetTradingBy get a trading by with a given ID
func (s *Storage) GetTradingBy(tradingByID uint64) (common.TradingBy, error) {
	var (
		result tradingByDB
	)
	err := s.stmts.getTradingBy.Get(&result, tradingByID)
	switch err {
	case sql.ErrNoRows:
		return common.TradingBy{}, common.ErrNotFound
	case nil:
		return result.ToCommon(), nil
	default:
		return common.TradingBy{}, err
	}
}
