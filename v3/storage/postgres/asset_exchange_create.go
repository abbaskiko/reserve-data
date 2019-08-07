package postgres

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/KyberNetwork/reserve-data/v3/common"
)

// CreateCreateAssetExchange create CreateAssetExchange
func (s *Storage) CreateCreateAssetExchange(req common.CreateCreateAssetExchange) (uint64, error) {
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
	err = tx.Stmtx(s.stmts.newCreateAssetExchange).Get(&id, jsonData)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	log.Printf("create pending asset exchange success with id = %d\n", id)
	return id, nil
}

type createAssetExchange struct {
	ID      uint64    `db:"id"`
	Created time.Time `db:"created"`
	Data    []byte    `db:"data"`
}

func (p createAssetExchange) toCommon() common.CreateAssetExchange {
	return common.CreateAssetExchange{
		ID:      p.ID,
		Created: p.Created,
		Data:    p.Data,
	}
}

// GetCreateAssetExchanges list all pending asset exchange
func (s *Storage) GetCreateAssetExchanges() ([]common.CreateAssetExchange, error) {
	var (
		recs   []createAssetExchange
		result []common.CreateAssetExchange
	)
	err := s.stmts.getCreateAssetExchanges.Select(&recs, nil)
	if err != nil {
		return nil, err
	}
	for _, p := range recs {
		result = append(result, p.toCommon())
	}
	return result, nil
}

// GetCreateAssetExchange list all pending asset exchange
func (s *Storage) GetCreateAssetExchange(id uint64) (common.CreateAssetExchange, error) {
	var (
		res createAssetExchange
	)
	err := s.stmts.getCreateAssetExchanges.Get(&res, id)
	if err != nil {
		return common.CreateAssetExchange{}, err
	}

	return res.toCommon(), nil
}

// ConfirmCreateAssetExchange confirm pending asset exchange, return err if any
func (s *Storage) ConfirmCreateAssetExchange(id uint64) error {
	var createAssetExchange common.CreateAssetExchange
	err := s.stmts.getCreateAssetExchanges.Get(&createAssetExchange, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("create asset_exchange request not found in database id=%d", id)
			return common.ErrNotFound
		}
		return err
	}
	var ccAssetExchange common.CreateCreateAssetExchange
	err = json.Unmarshal(createAssetExchange.Data, &ccAssetExchange)
	if err != nil {
		return err
	}
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)
	for _, r := range ccAssetExchange.AssetExchanges {
		_, err = s.createAssetExchange(tx, r.ExchangeID, r.AssetID, r.Symbol, r.DepositAddress, r.MinDeposit,
			r.WithdrawFee, r.TargetRecommended, r.TargetRatio)
		if err != nil {
			return err
		}
	}
	_, err = s.stmts.deleteCreateAssetExchange.Exec(id)
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

// RejectCreateAssetExchange reject pending asset exchange
func (s *Storage) RejectCreateAssetExchange(id uint64) error {
	_, err := s.stmts.deleteCreateAssetExchange.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
