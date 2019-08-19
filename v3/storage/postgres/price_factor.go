package postgres

import (
	"database/sql"

	"github.com/pkg/errors"

	common2 "github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/v3/common"
)

// PriceFactorAtTime store price factor into DB
func (s *Storage) CreatePriceFactor(pf common.PriceFactorAtTime) (uint64, error) {
	var lastID uint64
	err := s.stmts.newPriceFactor.Get(&lastID, pf.Timestamp, pf.Data)
	if err != nil {
		return 0, err
	}

	return lastID, err
}

type allPriceFactor struct {
	ID        uint64                      `db:"id"`
	TimePoint uint64                      `db:"timepoint"`
	Data      common.AssetPriceFactorList `db:"data"`
}

func (a allPriceFactor) toCommon() common.PriceFactorAtTime {
	return common.PriceFactorAtTime{
		ID:        a.ID,
		Timestamp: a.TimePoint,
		Data:      a.Data,
	}
}
func (s *Storage) GetPriceFactors(from, to uint64) ([]common.PriceFactorAtTime, error) {
	var dbResult []allPriceFactor
	err := s.stmts.getPriceFactor.Select(&dbResult, from, to)
	if err != nil {
		return nil, err
	}
	var res = make([]common.PriceFactorAtTime, 0, len(dbResult))
	for _, e := range dbResult {
		res = append(res, e.toCommon())
	}
	return res, nil
}

type setRateDB struct {
	ID        uint64 `db:"id"`
	Timepoint uint64 `db:"timepoint"`
	Status    bool   `db:"status"`
}

func (s *Storage) GetSetRateStatus() (bool, error) {
	var setRateResult setRateDB
	err := s.stmts.getSetRate.Get(&setRateResult)
	switch err {
	case sql.ErrNoRows:
		err := s.SetSetRateStatus(false)
		if err != nil {
			return false, errors.Wrapf(err, "fail to set set-rate control 1st time")
		}
		return false, nil
	case nil:
		return setRateResult.Status, nil
	default:
		return false, err
	}
}

func (s *Storage) SetSetRateStatus(status bool) error {
	_, err := s.stmts.newSetRate.Exec(common2.GetTimepoint(), status)
	return err
}

type rebalanceDB struct {
	ID        uint64 `db:"id"`
	Timepoint uint64 `db:"timepoint"`
	Status    bool   `db:"status"`
}

func (s *Storage) GetRebalanceStatus() (bool, error) {
	var rebalanceResult rebalanceDB
	err := s.stmts.getRebalance.Get(&rebalanceResult)
	switch err {
	case sql.ErrNoRows:
		err := s.SetRebalanceStatus(false)
		if err != nil {
			return false, errors.Wrapf(err, "fail to set rebalance 1st time")
		}
		return false, nil
	case nil:
		return rebalanceResult.Status, nil
	default:
		return false, err
	}
}

func (s *Storage) SetRebalanceStatus(status bool) error {
	_, err := s.stmts.newRebalance.Exec(common2.GetTimepoint(), status)
	return err
}
