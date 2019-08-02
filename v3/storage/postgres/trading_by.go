package postgres

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/KyberNetwork/reserve-data/v3/common"
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
	defer rollbackUnlessCommitted(tx)
	id, err := s.createTradingBy(tx, assetID, tradingPairID)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	log.Printf("asset trading pair #%d has been create successfully\n", id)
	return id, nil
}
