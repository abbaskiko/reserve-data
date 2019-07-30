package postgres

import (
	"database/sql"
	"log"

	"github.com/pkg/errors"

	"github.com/KyberNetwork/reserve-data/v3/common"
)

func (s *Storage) GetAssetExchange(id uint64) (common.AssetExchange, error) {
	var (
		result assetExchangeDB
	)
	err := s.stmts.getAssetExchange.Get(&result, nil, id)
	switch err {
	case sql.ErrNoRows:
		log.Printf("asset exchange not found id=%d", id)
		return common.AssetExchange{}, common.ErrNotFound
	case nil:
		return result.ToCommon(), nil
	default:
		return common.AssetExchange{}, errors.Errorf("failed to get asset exchange from database id=%d err=%s", id, err.Error())
	}
}
