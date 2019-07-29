package postgres

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/KyberNetwork/reserve-data/v3/common"
	"github.com/KyberNetwork/reserve-data/v3/storage"
)

// CreateUpdateExchange create a new update exchange,
// this will delete all old pending as we maintain 1 pending.
func (s *Storage) CreateUpdateExchange(c common.CreateUpdateExchange) (uint64, error) {
	var id uint64
	jsonData, err := json.Marshal(c)
	if err != nil {
		return 0, err
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer rollbackUnlessCommitted(tx)
	err = tx.Stmtx(s.stmts.newUpdateExchange).Get(&id, jsonData)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	log.Printf("create update exchange success with id=%d\n", id)
	return id, nil
}

type updateExchange struct {
	ID      uint64    `db:"id"`
	Created time.Time `db:"created"`
	Data    []byte    `db:"data"`
}

func (p *updateExchange) ToCommon() common.UpdateExchange {
	return common.UpdateExchange{
		ID:      p.ID,
		Created: p.Created,
		Data:    p.Data,
	}
}

// GetUpdateExchanges list all UpdateExchange(waiting for apply) exists in database
func (s *Storage) GetUpdateExchanges() ([]common.UpdateExchange, error) {
	var res []updateExchange
	err := s.stmts.getUpdateExchanges.Select(&res, nil)
	if err != nil {
		return nil, err
	}
	var result = make([]common.UpdateExchange, 0, 1) // although it's a slice, we expect only 1 for now.
	for _, p := range res {
		result = append(result, p.ToCommon())
	}
	return result, nil
}

// GetUpdateExchange return an UpdateExchange object with correspond id.
func (s *Storage) GetUpdateExchange(id uint64) (common.UpdateExchange, error) {
	var res updateExchange
	err := s.stmts.getUpdateExchanges.Get(&res, id)
	if err != nil {
		return common.UpdateExchange{}, err
	}

	return res.ToCommon(), nil
}

// RejectUpdateExchange delete UpdateExchange from DB.
func (s *Storage) RejectUpdateExchange(id uint64) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)
	_, err = tx.Stmtx(s.stmts.deleteUpdateExchange).Exec(id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("reject UpdateExchange success with id=%d\n", id)
	return nil
}

// ConfirmUpdateExchange apply pending changes in UpdateExchange object.
func (s *Storage) ConfirmUpdateExchange(id uint64) error {
	var update updateExchange
	err := s.stmts.getUpdateExchanges.Get(&update, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("update exchange not found in database id=%d", id)
			return common.ErrNotFound
		}
		return err
	}
	var r common.CreateUpdateExchange
	err = json.Unmarshal(update.Data, &r)
	if err != nil {
		return err
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)
	for _, e := range r.Exchanges {
		err = s.updateExchange(tx, e.ExchangeID, storage.UpdateExchangeOpts{
			TradingFeeMaker: e.TradingFeeMaker,
			TradingFeeTaker: e.TradingFeeTaker,
			Disable:         e.Disable,
		})
		if err != nil {
			return err
		}
	}

	_, err = tx.Stmtx(s.stmts.deleteUpdateExchange).Exec(id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("update exchange #%d has been confirm successfully\n", id)
	return nil
}
