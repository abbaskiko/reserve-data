package postgres

import (
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
