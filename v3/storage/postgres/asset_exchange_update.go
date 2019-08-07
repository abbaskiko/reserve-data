package postgres

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/KyberNetwork/reserve-data/v3/common"
)

// CreateUpdateAssetExchange create pending asset exchange
func (s *Storage) CreateUpdateAssetExchange(req common.CreateUpdateAssetExchange) (uint64, error) {
	var (
		id uint64
	)
	jsonData, err := json.Marshal(req)
	if err != nil {
		return 0, err
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer rollbackUnlessCommitted(tx)
	err = tx.Stmtx(s.stmts.newUpdateAssetExchange).Get(&id, jsonData)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	log.Printf("create pending asset exchange success with id = %d\n", id)
	return id, nil
}

type updateAssetExchange struct {
	ID      uint64    `db:"id"`
	Created time.Time `db:"created"`
	Data    []byte    `db:"data"`
}

func (p updateAssetExchange) toCommon() common.UpdateAssetExchange {
	return common.UpdateAssetExchange{
		ID:      p.ID,
		Created: p.Created,
		Data:    p.Data,
	}
}

// GetUpdateAssetExchanges list all pending asset exchange
func (s *Storage) GetUpdateAssetExchanges() ([]common.UpdateAssetExchange, error) {
	var (
		pendings []updateAssetExchange
		result   []common.UpdateAssetExchange
	)
	err := s.stmts.getUpdateAssetExchanges.Select(&pendings, nil)
	if err != nil {
		return nil, err
	}
	for _, p := range pendings {
		result = append(result, p.toCommon())
	}
	return result, nil
}

// GetUpdateAssetExchange list all pending asset exchange
func (s *Storage) GetUpdateAssetExchange(id uint64) (common.UpdateAssetExchange, error) {
	var (
		res updateAssetExchange
	)
	err := s.stmts.getUpdateAssetExchanges.Get(&res, id)
	if err != nil {
		return common.UpdateAssetExchange{}, err
	}

	return res.toCommon(), nil
}

// ConfirmUpdateAssetExchange confirm pending asset exchange, return err if any
func (s *Storage) ConfirmUpdateAssetExchange(id uint64) error {
	var updateAssetExchange common.UpdateAssetExchange
	err := s.stmts.getUpdateAssetExchanges.Get(&updateAssetExchange, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("update asset_exchange request not found in database id=%d", id)
			return common.ErrNotFound
		}
		return err
	}
	var ccAssetExchange common.CreateUpdateAssetExchange
	err = json.Unmarshal(updateAssetExchange.Data, &ccAssetExchange)
	if err != nil {
		return err
	}
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)
	for _, r := range ccAssetExchange.AssetExchanges {
		err = s.updateAssetExchange(tx, r.ID, r)
		if err != nil {
			return err
		}
	}
	_, err = s.stmts.deleteUpdateAssetExchange.Exec(id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("pending asset exchange #%d has been confirm successfully\n", id)
	return nil
}

// RejectUpdateAssetExchange reject pending asset exchange
func (s *Storage) RejectUpdateAssetExchange(id uint64) error {
	_, err := s.stmts.deleteUpdateAssetExchange.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
