package postgres

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/KyberNetwork/reserve-data/v3/common"
)

// CreateCreateTradingPair create a new pending TradingPair, this will delete all old pending as we maintain 1 pending only.
func (s *Storage) CreateCreateTradingPair(c common.CreateCreateTradingPair) (uint64, error) {
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
	err = tx.Stmtx(s.stmts.newCreateTradingPair).Get(&id, jsonData)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	log.Printf("create pending trading pair success with id=%d\n", id)
	return id, nil
}

type createTradingPair struct {
	ID      uint64    `db:"id"`
	Created time.Time `db:"created"`
	Data    []byte    `db:"data"`
}

func (p *createTradingPair) ToCommon() common.CreateTradingPair {
	return common.CreateTradingPair{
		ID:      p.ID,
		Created: p.Created,
		Data:    p.Data,
	}
}

// GetCreateTradingPairs list all CreateTradingPair exist in database
func (s *Storage) GetCreateTradingPairs() ([]common.CreateTradingPair, error) {
	var pendings []createTradingPair
	err := s.stmts.getCreateTradingPairs.Select(&pendings, nil)
	if err != nil {
		return nil, err
	}
	var result = make([]common.CreateTradingPair, 0, 1) // although it's a slice, we expect only 1 for now.
	for _, p := range pendings {
		result = append(result, p.ToCommon())
	}
	return result, nil
}

// GetCreateTradingPair get CreateTradingPair by id
func (s *Storage) GetCreateTradingPair(id uint64) (common.CreateTradingPair, error) {
	var result createTradingPair
	err := s.stmts.getCreateTradingPairs.Get(&result, nil)
	if err != nil {
		return common.CreateTradingPair{}, err
	}
	return result.ToCommon(), nil
}

func (s *Storage) RejectCreateTradingPair(id uint64) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)
	_, err = tx.Stmtx(s.stmts.deleteCreateTradingPair).Exec(id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("reject pending trading pair success with id=%d\n", id)
	return nil
}

func (s *Storage) ConfirmCreateTradingPair(id uint64) error {
	var pending createTradingPair
	err := s.stmts.getCreateTradingPairs.Get(&pending, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("create trading pair request not found in database id=%d", id)
			return common.ErrNotFound
		}
		return err
	}
	var createCreateTradingPair common.CreateCreateTradingPair
	err = json.Unmarshal(pending.Data, &createCreateTradingPair)
	if err != nil {
		return err
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)
	for _, a := range createCreateTradingPair.TradingPairs {
		_, err := s.createTradingPair(tx, a.ExchangeID, a.Base, a.Quote, a.PricePrecision, a.AmountPrecision,
			a.AmountLimitMin, a.AmountLimitMax, a.PriceLimitMin, a.PriceLimitMax, a.MinNotional)
		if err != nil {
			return err
		}
	}
	_, err = tx.Stmtx(s.stmts.deleteCreateTradingPair).Exec(id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("pending trading pair #%d has been confirm successfully\n", id)
	return nil
}
