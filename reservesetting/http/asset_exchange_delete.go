package http

import (
	"github.com/KyberNetwork/reserve-data/reservesetting/common"
)

func (s *Server) checkDeleteAssetExchangeParams(entry common.DeleteAssetExchangeEntry) error {
	assetExchange, err := s.storage.GetAssetExchange(entry.AssetExchangeID)
	if err != nil {
		return err
	}

	if len(assetExchange.TradingPairs) != 0 {
		return common.ErrAssetExchangeDeleteViolation
	}
	return nil
}
