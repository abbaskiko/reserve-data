package postgres

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/KyberNetwork/reserve-data/v3/common"
	"github.com/KyberNetwork/reserve-data/v3/storage"
)

// CreateUpdateAsset create a new update asset,
// this will delete all old pending as we maintain 1 pending update asset only.
func (s *Storage) CreateUpdateAsset(c common.CreateUpdateAsset) (uint64, error) {
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
	err = tx.Stmtx(s.stmts.newUpdateAsset).Get(&id, jsonData)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	log.Printf("create update asset success with id=%d\n", id)
	return id, nil
}

type updateAsset struct {
	ID      uint64    `db:"id"`
	Created time.Time `db:"created"`
	Data    []byte    `db:"data"`
}

func (p *updateAsset) ToCommon() common.UpdateAsset {
	return common.UpdateAsset{
		ID:      p.ID,
		Created: p.Created,
		Data:    p.Data,
	}
}

// GetUpdateAssets list all UpdateAsset(waiting for apply) exists in database
func (s *Storage) GetUpdateAssets() ([]common.UpdateAsset, error) {
	var res []updateAsset
	err := s.stmts.getUpdateAssets.Select(&res, nil)
	if err != nil {
		return nil, err
	}
	var result = make([]common.UpdateAsset, 0, 1) // although it's a slice, we expect only 1 for now.
	for _, p := range res {
		result = append(result, p.ToCommon())
	}
	return result, nil
}

// RejectUpdateAsset reject by delete that UpdateAsset.
func (s *Storage) RejectUpdateAsset(id uint64) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)
	_, err = tx.Stmtx(s.stmts.deleteUpdateAsset).Exec(id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("reject UpdateAsset success with id=%d\n", id)
	return nil
}

func (s *Storage) ConfirmUpdateAsset(id uint64) error {
	var update updateAsset
	err := s.stmts.getUpdateAssets.Get(&update, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("update asset not found in database id=%d", id)
			return common.ErrNotFound
		}
		return err
	}
	var r common.CreateUpdateAsset
	err = json.Unmarshal(update.Data, &r)
	if err != nil {
		return err
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)
	for _, e := range r.Assets {
		err = s.updateAsset(tx, e.AssetID, storage.UpdateAssetOpts{
			Symbol:       e.Symbol,
			Transferable: e.Transferable,
			Address:      e.Address,
			IsQuote:      e.IsQuote,
			Rebalance:    e.Rebalance,
			SetRate:      e.SetRate,
			Decimals:     e.Decimals,
			Name:         e.Name,
		})
		if err != nil {
			return err
		}
	}

	_, err = tx.Stmtx(s.stmts.deleteUpdateAsset).Exec(id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("update asset #%d has been confirm successfully\n", id)
	return nil
}
