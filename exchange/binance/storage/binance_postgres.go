package storage

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/common/postgres"
	"github.com/KyberNetwork/reserve-data/exchange"
)

const (
	schema = `
		CREATE TABLE IF NOT EXISTS "binance_trade_history"(
		    id 				SERIAL PRIMARY KEY,
		    pair_id			BIGINT,
		    trade_id		TEXT NOT NULL,
		    price 			FLOAT NOT NULL,
		    qty 			FLOAT NOT NULL, 
		    type			TEXT NOT NULL,
		    time			BIGINT
		);
	`
)

// postgresStorage implements binance storage in postgres
type postgresStorage struct {
	db    *sqlx.DB
	stmts preparedStmt
}
type preparedStmt struct {
	storeHistoryStmt     *sqlx.NamedStmt
	getHistoryStmt       *sqlx.Stmt
	getLastIDHistoryStmt *sqlx.Stmt
}

// NewPostgresStorage creates new obj exchange.BinanceStorage with db engine = postgres
func NewPostgresStorage(db *sqlx.DB) (exchange.BinanceStorage, error) {
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to intialize database schema err=%s", err.Error())
	}
	storeHistoryStmt, err := db.PrepareNamed(`INSERT INTO "binance_trade_history"
		(pair_id, trade_id, price, qty, type, time)
		VALUES(:pair_id, :trade_id, :price, :qty, :type, :time)`)
	if err != nil {
		return nil, err
	}
	getHistoryStmt, err := db.Preparex(`SELECT pair_id, trade_id, price, qty, type, time 
		FROM "binance_trade_history"
		WHERE time >= $1 AND time <= $2`)
	if err != nil {
		return nil, err
	}

	getLastIDHistoryStmt, err := db.Preparex(`SELECT pair_id, trade_id, price, qty, type, time FROM "binance_trade_history"
											WHERE pair_id = $1 
											ORDER BY time DESC, trade_id DESC;`)
	if err != nil {
		return nil, err
	}
	storage := &postgresStorage{
		db: db,
		stmts: preparedStmt{
			storeHistoryStmt:     storeHistoryStmt,
			getHistoryStmt:       getHistoryStmt,
			getLastIDHistoryStmt: getLastIDHistoryStmt,
		},
	}
	return storage, nil
}

type exchangeTradeHistoryDB struct {
	PairID  uint64  `db:"pair_id"`
	TradeID string  `db:"trade_id"`
	Price   float64 `db:"price"`
	Qty     float64 `db:"qty"`
	Type    string  `db:"type"`
	Time    uint64  `db:"time"`
}

// StoreTradeHistory implements exchange.BinanceStorage and store trade history
func (s *postgresStorage) StoreTradeHistory(data common.ExchangeTradeHistory) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer postgres.RollbackUnlessCommitted(tx)

	for pairID, tradeHistory := range data {
		for _, history := range tradeHistory {
			_, err = tx.NamedStmt(s.stmts.storeHistoryStmt).Exec(exchangeTradeHistoryDB{
				PairID:  pairID,
				TradeID: history.ID,
				Price:   history.Price,
				Qty:     history.Qty,
				Type:    history.Type,
				Time:    history.Timestamp,
			})
			if err != nil {
				return err
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// GetTradeHistory implements exchange.BinanceStorage and get trade history within a time period
func (s *postgresStorage) GetTradeHistory(fromTime, toTime uint64) (common.ExchangeTradeHistory, error) {
	var result = make(common.ExchangeTradeHistory)
	var records []exchangeTradeHistoryDB
	err := s.stmts.getHistoryStmt.Select(&records, fromTime, toTime)
	if err != nil {
		return result, err
	}
	for _, r := range records {
		tradeHistory := result[r.PairID]
		tradeHistory = append(tradeHistory, common.TradeHistory{
			ID:        r.TradeID,
			Price:     r.Price,
			Qty:       r.Qty,
			Type:      r.Type,
			Timestamp: r.Time,
		})
		result[r.PairID] = tradeHistory
	}
	return result, nil
}

// GetLastIDTradeHistory implements exchange.BinanceStorage and get the last ID with a correspond pairID
func (s *postgresStorage) GetLastIDTradeHistory(pairID uint64) (string, error) {
	var record exchangeTradeHistoryDB
	err := s.stmts.getLastIDHistoryStmt.Get(&record, pairID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.Errorf("no history with pair_id=%v", pairID)
		}
		return "", err
	}
	return record.TradeID, nil
}
