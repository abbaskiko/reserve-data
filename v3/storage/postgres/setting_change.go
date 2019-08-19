package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/KyberNetwork/reserve-data/v3/common"
)

// CreateSettingChange creates an setting change in database and return id
func (s *Storage) CreateSettingChange(obj common.SettingChange) (uint64, error) {
	var id uint64
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to parse json data %+v", obj)
	}
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer rollbackUnlessCommitted(tx)

	if err = tx.Stmtx(s.stmts.newSettingChange).Get(&id, jsonData); err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	log.Printf("create setting change success with id=%d\n", id)
	return id, nil
}

type settingChangeDB struct {
	ID      uint64    `db:"id"`
	Created time.Time `db:"created"`
	Data    []byte    `db:"data"`
}

func (objDB settingChangeDB) ToCommon() (common.SettingChangeResponse, error) {
	var settingChange common.SettingChange
	err := json.Unmarshal(objDB.Data, &settingChange)
	if err != nil {
		return common.SettingChangeResponse{}, err
	}
	return common.SettingChangeResponse{
		ChangeList: settingChange.ChangeList,
		ID:         objDB.ID,
		Created:    objDB.Created,
	}, nil
}

// GetSettingChange returns a object with a given id and type
func (s *Storage) GetSettingChange(id uint64) (common.SettingChangeResponse, error) {
	var dbResult settingChangeDB
	err := s.stmts.getSettingChange.Get(&dbResult, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("setting change not found in database id=%d\n", id)
			return common.SettingChangeResponse{}, common.ErrNotFound
		}
		return common.SettingChangeResponse{}, err
	}
	res, err := dbResult.ToCommon()
	if err != nil {
		return common.SettingChangeResponse{}, err
	}
	return res, nil
}

// GetSettingChanges return list setting change.
func (s *Storage) GetSettingChanges() ([]common.SettingChangeResponse, error) {
	var dbResult []settingChangeDB
	err := s.stmts.getSettingChange.Select(&dbResult, nil)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	var result = make([]common.SettingChangeResponse, 0, 1) // although it's a slice, we expect only 1 for now.
	for _, p := range dbResult {
		rr, err := p.ToCommon()
		if err != nil {
			return nil, err
		}
		result = append(result, rr)
	}
	return result, nil
}

// RejectSettingChange delete setting change with a given id
func (s *Storage) RejectSettingChange(id uint64) error {
	var returnedID uint64
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer rollbackUnlessCommitted(tx)
	err = tx.Stmtx(s.stmts.deleteSettingChange).Get(&returnedID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return common.ErrNotFound
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("reject setting change success with id=%d\n", id)
	return nil
}
func (s *Storage) getSettingChange(tx *sqlx.Tx, id uint64) (common.SettingChangeResponse, error) {
	var dbResult settingChangeDB
	err := tx.Stmtx(s.stmts.getSettingChange).Get(&dbResult, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("setting change not found in database id=%d\n", id)
			return common.SettingChangeResponse{}, common.ErrNotFound
		}
		return common.SettingChangeResponse{}, err
	}
	res, err := dbResult.ToCommon()
	if err != nil {
		log.Printf("failed to convert to common setting change, err=%v\n", err)
		return common.SettingChangeResponse{}, err
	}
	return res, nil
}
func (s *Storage) applyChange(tx *sqlx.Tx, i int, entry common.SettingChangeEntry) error {
	var err error
	switch entry.Type {
	case common.ChangeTypeChangeAssetAddr:
		u, ok := entry.Data.(*common.ChangeAssetAddressEntry)
		if !ok {
			msg := fmt.Sprintf("bad cast at %d to %s\n", i, common.ChangeTypeChangeAssetAddr)
			log.Println(msg)
			return errors.Wrap(err, msg)
		}
		err = s.changeAssetAddress(tx, u.ID, u.Address)
		if err != nil {
			msg := fmt.Sprintf("change asset address %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	case common.ChangeTypeCreateAsset:
		a, ok := entry.Data.(*common.CreateAssetEntry)
		if !ok {
			msg := fmt.Sprintf("bad cast at %d to %s\n", i, common.ChangeTypeCreateAsset)
			log.Println(msg)
			return errors.Wrap(err, msg)
		}
		_, err = s.createAsset(tx, a.Symbol, a.Name, a.Address, a.Decimals, a.Transferable, a.SetRate, a.Rebalance,
			a.IsQuote, a.PWI, a.RebalanceQuadratic, a.Exchanges, a.Target)
		if err != nil {
			msg := fmt.Sprintf("create asset %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	case common.ChangeTypeCreateAssetExchange:
		a, ok := entry.Data.(*common.CreateAssetExchangeEntry)
		if !ok {
			msg := fmt.Sprintf("bad cast at %d to %s\n", i, common.ChangeTypeCreateAssetExchange)
			log.Println(msg)
			return errors.Wrap(err, msg)
		}
		_, err = s.createAssetExchange(tx, a.ExchangeID, a.AssetID, a.Symbol, a.DepositAddress, a.MinDeposit,
			a.WithdrawFee, a.TargetRecommended, a.TargetRatio)
		if err != nil {
			msg := fmt.Sprintf("create asset exchange %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	case common.ChangeTypeCreateTradingBy:
		a, ok := entry.Data.(*common.CreateTradingByEntry)
		if !ok {
			msg := fmt.Sprintf("bad cast at %d to %s\n", i, common.ChangeTypeCreateTradingBy)
			log.Println(msg)
			return errors.Wrap(err, msg)
		}
		_, err = s.createTradingBy(tx, a.AssetID, a.TradingPairID)
		if err != nil {
			msg := fmt.Sprintf("create trading by %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	case common.ChangeTypeCreateTradingPair:
		a, ok := entry.Data.(*common.CreateTradingPairEntry)
		if !ok {
			msg := fmt.Sprintf("bad cast at %d to %s\n", i, common.ChangeTypeCreateTradingPair)
			log.Println(msg)
			return errors.Wrap(err, msg)
		}
		_, err = s.createTradingPair(tx, a.ExchangeID, a.Base, a.Quote, a.PricePrecision, a.AmountPrecision, a.AmountLimitMin,
			a.AmountLimitMax, a.PriceLimitMin, a.PriceLimitMax, a.MinNotional)
		if err != nil {
			msg := fmt.Sprintf("create trading pair %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	case common.ChangeTypeUpdateAsset:
		a, ok := entry.Data.(*common.UpdateAssetEntry)
		if !ok {
			msg := fmt.Sprintf("bad cast at %d to %s\n", i, common.ChangeTypeUpdateAsset)
			log.Println(msg)
			return errors.Wrap(err, msg)
		}
		err = s.updateAsset(tx, a.AssetID, *a)
		if err != nil {
			msg := fmt.Sprintf("update asset %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	case common.ChangeTypeUpdateAssetExchange:
		a, ok := entry.Data.(*common.UpdateAssetExchangeEntry)
		if !ok {
			msg := fmt.Sprintf("bad cast at %d to %s\n", i, common.ChangeTypeUpdateAssetExchange)
			log.Println(msg)
			return errors.Wrap(err, msg)
		}
		err = s.updateAssetExchange(tx, a.ID, *a)
		if err != nil {
			msg := fmt.Sprintf("bad cast at %d to %s\n", i, common.ChangeTypeUpdateAssetExchange)
			log.Println(msg)
			return err
		}
	case common.ChangeTypeUpdateExchange:
		a, ok := entry.Data.(*common.UpdateExchangeEntry)
		if !ok {
			msg := fmt.Sprintf("bad cast at %d to %s\n", i, common.ChangeTypeUpdateExchange)
			log.Println(msg)
			return errors.Wrap(err, msg)
		}
		err = s.updateExchange(tx, a.ExchangeID, *a)
		if err != nil {
			msg := fmt.Sprintf("update exchange %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	case common.ChangeTypeDeleteAssetExchange:
		a, ok := entry.Data.(*common.DeleteAssetExchangeEntry)
		if !ok {
			msg := fmt.Sprintf("bad cast at %d to %s\n", i, common.ChangeTypeDeleteAssetExchange)
			log.Println(msg)
			return errors.Wrap(err, msg)
		}
		err = s.deleteAssetExchange(tx, a.AssetExchangeID)
		if err != nil {
			msg := fmt.Sprintf("delete asset exchange id=%d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	case common.ChangeTypeDeleteTradingBy:
	case common.ChangeTypeDeleteTradingPair:
		a, ok := entry.Data.(*common.DeleteTradingPairEntry)
		if !ok {
			msg := fmt.Sprintf("bad cast at %d to %s\n", i, common.ChangeTypeUpdateExchange)
			log.Println(msg)
			return errors.Wrap(err, msg)
		}
		err = s.deleteTradingPair(tx, a.TradingPairID)
		if err != nil {
			msg := fmt.Sprintf("delete trading pair %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	case common.ChangeTypeUnknown:
		return fmt.Errorf("change type not set at %d", i)
	}
	return nil
}

// ConfirmSettingChange apply setting change with a given id
func (s *Storage) ConfirmSettingChange(id uint64, commit bool) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "create transaction error")
	}
	defer rollbackUnlessCommitted(tx)
	changeObj, err := s.getSettingChange(tx, id)
	if err != nil {
		return errors.Wrap(err, "get setting change error")
	}

	for i, change := range changeObj.ChangeList {
		if err = s.applyChange(tx, i, change); err != nil {
			return err
		}
	}
	_, err = tx.Stmtx(s.stmts.deleteSettingChange).Exec(id)
	if err != nil {
		return err
	}
	if commit {
		if err := tx.Commit(); err != nil {
			log.Printf("setting change id=%d has been failed to confirm, err=%v\n", id, err)
			return err
		}
		log.Printf("setting change has been confirmed successfully, id=%d\n", id)
		return nil
	}
	return nil
}
