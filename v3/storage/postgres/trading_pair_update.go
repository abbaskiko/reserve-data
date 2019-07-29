package postgres

import (
	"encoding/json"
	"log"
	"time"

	"github.com/KyberNetwork/reserve-data/v3/common"
)

// CreateUpdateTradingPair create a new UpdateTradingPair, this will delete all old exist pending.
func (s *Storage) CreateUpdateTradingPair(c common.CreateUpdateTradingPair) (uint64, error) {
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
	err = tx.Stmtx(s.stmts.newUpdateTradingPair).Get(&id, jsonData)
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

type updateTradingPair struct {
	ID      uint64    `db:"id"`
	Created time.Time `db:"created"`
	Data    []byte    `db:"data"`
}

func (p *updateTradingPair) ToCommon() common.UpdateTradingPair {
	return common.UpdateTradingPair{
		ID:      p.ID,
		Created: p.Created,
		Data:    p.Data,
	}
}

// GetUpdateTradingPairs list all GetUpdateTradingPair exist in database
func (s *Storage) GetUpdateTradingPairs() ([]common.UpdateTradingPair, error) {
	var pendings []updateTradingPair
	err := s.stmts.getUpdateTradingPairs.Select(&pendings, nil)
	if err != nil {
		return nil, err
	}
	var result = make([]common.UpdateTradingPair, 0, 1) // although it's a slice, we expect only 1 for now.
	for _, p := range pendings {
		result = append(result, p.ToCommon())
	}
	return result, nil
}

// GetUpdateTradingPair get GetUpdateTradingPair by id
func (s *Storage) GetUpdateTradingPair(id uint64) (common.UpdateTradingPair, error) {
	var result updateTradingPair
	err := s.stmts.getCreateTradingPairs.Get(&result, nil)
	if err != nil {
		return common.UpdateTradingPair{}, err
	}
	return result.ToCommon(), nil
}

func (s *Storage) RejectUpdateTradingPair(id uint64) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)
	_, err = tx.Stmtx(s.stmts.deleteUpdateTradingPair).Exec(id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("reject update trading pair success with id=%d\n", id)
	return nil
}

func (s *Storage) ConfirmUpdateTradingPair(id uint64) error {
	var pending updateTradingPair
	err := s.stmts.getUpdateTradingPairs.Get(&pending, id)
	if err != nil {
		return err
	}
	var createUpdateTradingPair common.CreateUpdateTradingPair
	err = json.Unmarshal(pending.Data, &createUpdateTradingPair)
	if err != nil {
		return err
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)
	for _, a := range createUpdateTradingPair.TradingPairs {
		err = s.updateTradingPair(tx, a.ID, a)
		if err != nil {
			return err
		}
	}
	_, err = tx.Stmtx(s.stmts.deleteUpdateTradingPair).Exec(id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("update trading pair #%d has been confirm successfully\n", id)
	return nil
}
