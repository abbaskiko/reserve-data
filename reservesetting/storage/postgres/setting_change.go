package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/KyberNetwork/reserve-data/reservesetting/common"
)

const (
	settingChangeCatUnique = "setting_change_cat_key"
)

// CreateSettingChange creates an setting change in database and return id
func (s *Storage) CreateSettingChange(cat common.ChangeCatalog, obj common.SettingChange) (uint64, error) {
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

	if err = tx.Stmtx(s.stmts.newSettingChange).Get(&id, cat.String(), jsonData); err != nil {
		pErr, ok := err.(*pq.Error)
		if !ok {
			return 0, fmt.Errorf("unknown returned err=%s", err.Error())
		}

		log.Printf("failed to create new setting change err=%s", pErr.Message)
		if pErr.Code == errCodeUniqueViolation && pErr.Constraint == settingChangeCatUnique {
			return 0, common.ErrSettingChangeExists
		}
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

// GetSettingChange returns a object with a given id
func (s *Storage) GetSettingChange(id uint64) (common.SettingChangeResponse, error) {
	return s.getSettingChange(nil, id)
}

func (s *Storage) getSettingChange(tx *sqlx.Tx, id uint64) (common.SettingChangeResponse, error) {
	var dbResult settingChangeDB
	sts := s.stmts.getSettingChange
	if tx != nil {
		sts = tx.Stmtx(sts)
	}
	err := sts.Get(&dbResult, id, nil)
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

// GetSettingChanges return list setting change.
func (s *Storage) GetSettingChanges(cat common.ChangeCatalog) ([]common.SettingChangeResponse, error) {
	var dbResult []settingChangeDB
	err := s.stmts.getSettingChange.Select(&dbResult, nil, cat.String())
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

func (s *Storage) applyChange(tx *sqlx.Tx, i int, entry common.SettingChangeEntry) error {
	var err error
	switch e := entry.Data.(type) {
	case *common.ChangeAssetAddressEntry:
		err = s.changeAssetAddress(tx, e.ID, e.Address)
		if err != nil {
			msg := fmt.Sprintf("change asset address %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	case *common.CreateAssetEntry:
		_, err = s.createAsset(tx, e.Symbol, e.Name, e.Address, e.Decimals, e.Transferable, e.SetRate, e.Rebalance,
			e.IsQuote, e.PWI, e.RebalanceQuadratic, e.Exchanges, e.Target)
		if err != nil {
			msg := fmt.Sprintf("create asset %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	case *common.CreateAssetExchangeEntry:
		_, err = s.createAssetExchange(tx, e.ExchangeID, e.AssetID, e.Symbol, e.DepositAddress, e.MinDeposit,
			e.WithdrawFee, e.TargetRecommended, e.TargetRatio, e.TradingPairs)
		if err != nil {
			msg := fmt.Sprintf("create asset exchange %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}

	case *common.CreateTradingByEntry:
		_, err = s.createTradingBy(tx, e.AssetID, e.TradingPairID)
		if err != nil {
			msg := fmt.Sprintf("create trading by %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	case *common.CreateTradingPairEntry:
		pairID, err := s.createTradingPair(tx, e.ExchangeID, e.Base, e.Quote, e.PricePrecision, e.AmountPrecision, e.AmountLimitMin,
			e.AmountLimitMax, e.PriceLimitMin, e.PriceLimitMax, e.MinNotional)
		if err != nil {
			msg := fmt.Sprintf("create trading pair %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
		// TODO: should we move create trading by into createTradingPair
		_, err = s.createTradingBy(tx, e.AssetID, pairID)
		if err != nil {
			msg := fmt.Sprintf("create trading by at position %d failed, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	case *common.UpdateAssetEntry:
		err = s.updateAsset(tx, e.AssetID, *e)
		if err != nil {
			msg := fmt.Sprintf("update asset %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	case *common.UpdateAssetExchangeEntry:
		err = s.updateAssetExchange(tx, e.ID, *e)
		if err != nil {
			msg := fmt.Sprintf("update asset exchange at %d failed, err = %v\n", i, err)
			log.Println(msg)
			return err
		}
	case *common.UpdateExchangeEntry:
		err = s.updateExchange(tx, e.ExchangeID, *e)
		if err != nil {
			msg := fmt.Sprintf("update exchange %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	case *common.DeleteAssetExchangeEntry:
		err = s.deleteAssetExchange(tx, e.AssetExchangeID)
		if err != nil {
			msg := fmt.Sprintf("delete asset exchange id=%d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	// case common.ChangeTypeDeleteTradingBy:
	case *common.DeleteTradingPairEntry:
		err = s.deleteTradingPair(tx, e.TradingPairID)
		if err != nil {
			msg := fmt.Sprintf("delete trading pair %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	case *common.UpdateStableTokenParamsEntry:
		err = s.updateStableTokenParams(tx, e.Params)
		if err != nil {
			msg := fmt.Sprintf("update stable token params %d, err=%v\n", i, err)
			log.Println(msg)
			return err
		}
	default:
		return fmt.Errorf("unexpected change object %+v", e)
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
	log.Printf("setting change will be reverted due commit=false, id=%d\n", id)
	return nil
}
