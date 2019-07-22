package postgres

import (
	"encoding/json"
	"log"
	"time"

	"github.com/KyberNetwork/reserve-data/v3/common"
)

// CreateCreateAsset create a new pending asset, this will delete all old pending as we maintain 1 pending asset only.
func (s *Storage) CreateCreateAsset(c common.CreateCreateAsset) (uint64, error) {
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
	err = tx.Stmtx(s.stmts.newCreateAsset).Get(&id, jsonData)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	log.Printf("create pending asset success with id=%d\n", id)
	return id, nil
}

type createAsset struct {
	ID      uint64    `db:"id"`
	Created time.Time `db:"created"`
	Data    []byte    `db:"data"`
}

func (p *createAsset) ToCommon() common.CreateAsset {
	return common.CreateAsset{
		ID:      p.ID,
		Created: p.Created,
		Data:    p.Data,
	}
}

// GetCreateAssets list all CreateAsset exist in database
func (s *Storage) GetCreateAssets() ([]common.CreateAsset, error) {
	var pendings []createAsset
	err := s.stmts.getCreateAssets.Select(&pendings, nil)
	if err != nil {
		return nil, err
	}
	var result = make([]common.CreateAsset, 0, 1) // although it's a slice, we expect only 1 for now.
	for _, p := range pendings {
		result = append(result, p.ToCommon())
	}
	return result, nil
}

func (s *Storage) RejectCreateAsset(id uint64) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)
	_, err = tx.Stmtx(s.stmts.deleteCreateAsset).Exec(id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("reject pending asset success with id=%d\n", id)
	return nil
}

func (s *Storage) ConfirmCreateAsset(id uint64) error {
	var pending createAsset
	err := s.stmts.getCreateAssets.Get(&pending, id)
	if err != nil {
		return err
	}
	var createCreateAsset common.CreateCreateAsset
	err = json.Unmarshal(pending.Data, &createCreateAsset)
	if err != nil {
		return err
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)
	for _, a := range createCreateAsset.AssetInputs {
		_, err := s.createAsset(tx, a.Symbol, a.Name, a.Address, a.Decimals, a.Transferable, a.SetRate, a.Rebalance,
			a.IsQuote, a.PWI, a.RebalanceQuadratic, a.Exchanges, a.Target)
		if err != nil {
			return err
		}
	}
	_, err = tx.Stmtx(s.stmts.deleteCreateAsset).Exec(id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("pending asset #%d has been confirm successfully\n", id)
	return nil
}
