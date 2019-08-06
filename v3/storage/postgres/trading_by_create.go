package postgres

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/KyberNetwork/reserve-data/v3/common"
)

// CreateCreateTradingBy create a new pending TradingBy, this will delete all old pending as we maintain 1 pending only.
func (s *Storage) CreateCreateTradingBy(c common.CreateCreateTradingBy) (uint64, error) {
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
	err = tx.Stmtx(s.stmts.newCreateTradingBy).Get(&id, jsonData)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	log.Printf("create pending trading by success with id=%d\n", id)
	return id, nil
}

type createTradingBy struct {
	ID      uint64    `db:"id"`
	Created time.Time `db:"created"`
	Data    []byte    `db:"data"`
}

func (p *createTradingBy) ToCommon() common.CreateTradingBy {
	return common.CreateTradingBy{
		ID:      p.ID,
		Created: p.Created,
		Data:    p.Data,
	}
}

// GetCreateTradingBy list all CreateTradingBy exist in database
func (s *Storage) GetCreateTradingBys() ([]common.CreateTradingBy, error) {
	var pendings []createTradingBy
	err := s.stmts.getCreateTradingBy.Select(&pendings, nil)
	if err != nil {
		return nil, err
	}
	var result = make([]common.CreateTradingBy, 0, 1) // although it's a slice, we expect only 1 for now.
	for _, p := range pendings {
		result = append(result, p.ToCommon())
	}
	return result, nil
}

// GetCreateTradingBy get CreateTradingBy by id
func (s *Storage) GetCreateTradingBy(id uint64) (common.CreateTradingBy, error) {
	var result createTradingBy
	err := s.stmts.getCreateTradingBy.Get(&result, id)
	if err != nil {
		return common.CreateTradingBy{}, err
	}
	return result.ToCommon(), nil
}

// RejectCreateTradingBy to delete pending create trading by request
func (s *Storage) RejectCreateTradingBy(id uint64) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)
	_, err = tx.Stmtx(s.stmts.deleteCreateTradingBy).Exec(id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("reject pending trading by success with id=%d\n", id)
	return nil
}

// ConfirmCreateTradingBy to execute the pending trading by request
func (s *Storage) ConfirmCreateTradingBy(id uint64) error {
	var pending createTradingBy
	err := s.stmts.getCreateTradingBy.Get(&pending, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("create trading by request not found in database id=%d", id)
			return common.ErrNotFound
		}
		return err
	}
	var createCreateTradingBy common.CreateCreateTradingBy
	err = json.Unmarshal(pending.Data, &createCreateTradingBy)
	if err != nil {
		return err
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)
	for _, tradingByEntry := range createCreateTradingBy.TradingBys {
		_, err := s.createTradingBy(tx, tradingByEntry.AssetID, tradingByEntry.TradingPairID)
		if err != nil {
			return err
		}
	}
	_, err = tx.Stmtx(s.stmts.deleteCreateTradingBy).Exec(id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("pending trading by #%d has been confirm successfully\n", id)
	return nil
}
