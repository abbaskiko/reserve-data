package postgres

import (
	"encoding/json"
	"log"
	"time"

	"github.com/KyberNetwork/reserve-data/v3/common"
)

// CreateChangeAssetAddress create a new change asset address
func (s *Storage) CreateChangeAssetAddress(c common.ChangeAssetAddress) (uint64, error) {
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
	err = tx.Stmtx(s.stmts.newChangeAssetAddress).Get(&id, jsonData)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	log.Printf("create change asset address success with id=%d\n", id)
	return id, nil
}

type changeAssetAddress struct {
	ID      uint64    `db:"id"`
	Created time.Time `db:"created"`
	Data    []byte    `db:"data"`
}

func (caa changeAssetAddress) toCommon() common.ChangeAssetAddressPending {
	return common.ChangeAssetAddressPending{
		ID:      caa.ID,
		Created: caa.Created,
		Data:    caa.Data,
	}
}

// GetChangeAssetAddress get a pending change asset address by id
func (s *Storage) GetChangeAssetAddress(id uint64) (common.ChangeAssetAddressPending, error) {
	var (
		res changeAssetAddress
	)
	err := s.stmts.getChangeAssetAddresses.Get(&res, id)
	if err != nil {
		return common.ChangeAssetAddressPending{}, err
	}

	return res.toCommon(), nil
}

// GetChangeAssetAddresses get all new pending change asset address
func (s *Storage) GetChangeAssetAddresses() ([]common.ChangeAssetAddressPending, error) {
	var (
		pending []changeAssetAddress
		res     []common.ChangeAssetAddressPending
	)
	err := s.stmts.getChangeAssetAddresses.Get(&res, nil)
	if err != nil {
		return res, err
	}
	for _, p := range pending {
		res = append(res, p.toCommon())
	}
	return res, nil
}

// RejectChangeAssetAddress reject by delete the pending change asset address.
func (s *Storage) RejectChangeAssetAddress(id uint64) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)
	_, err = tx.Stmtx(s.stmts.deleteChangeAssetAddress).Exec(id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("reject ChangeAssetAddress success with id=%d\n", id)
	return nil
}

// ConfirmChangeAssetAddress confirm the pending change asset address.
func (s *Storage) ConfirmChangeAssetAddress(id uint64) error {
	recordedData, err := s.GetChangeAssetAddress(id)
	if err != nil {
		log.Printf("update change_asset_addresses request not found in database id=%d", id)
		return common.ErrNotFound
	}
	var changeAssetAddress common.ChangeAssetAddress
	if err = json.Unmarshal(recordedData.Data, &changeAssetAddress); err != nil {
		return err
	}
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)

	for _, a := range changeAssetAddress.Assets {
		_, err = tx.Stmtx(s.stmts.changeAssetAddress).Exec(a.ID, a.Address)
		if err != nil {
			return err
		}
	}
	_, err = tx.Stmtx(s.stmts.deleteChangeAssetAddress).Exec(id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("confirm ChangeAssetAddress success with id=%d\n", id)
	return nil
}
