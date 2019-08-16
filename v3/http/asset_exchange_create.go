package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/v3/common"
)

func (s *Server) createCreateAssetExchange(c *gin.Context) {
	var createAssetExchange common.CreateCreateAssetExchange

	err := c.ShouldBindJSON(&createAssetExchange)

	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	for index, entry := range createAssetExchange.AssetExchanges {
		if err := s.checkCreateAssetExchangeParams(entry); err != nil {
			httputil.ResponseFailure(c, httputil.WithError(err), httputil.WithField("index", index),
				httputil.WithField("asset_id", entry.AssetID), httputil.WithField("exchange_id", entry.ExchangeID))
			return
		}
		for i, tp := range entry.TradingPairs {
			if tp.Base == 0 {
				entry.TradingPairs[i].Base = entry.AssetID
			}
			if tp.Quote == 0 {
				entry.TradingPairs[i].Quote = entry.AssetID
			}
		}
		createAssetExchange.AssetExchanges[index] = entry
	}

	id, err := s.storage.CreatePendingObject(createAssetExchange, common.PendingTypeCreateAssetExchange)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithField("id", id))
}

func (s *Server) checkCreateAssetExchangeParams(createEntry common.CreateAssetExchangeEntry) error {
	asset, err := s.storage.GetAsset(createEntry.AssetID)
	if err != nil {
		return errors.Wrap(err, "asset not found")
	}

	_, err = s.storage.GetExchange(createEntry.ExchangeID)
	if err != nil {
		return errors.Wrap(err, "exchange not found")
	}

	for _, exchange := range asset.Exchanges {
		if exchange.ExchangeID == createEntry.ExchangeID {
			return common.ErrAssetExchangeAlreadyExist
		}
	}
	if asset.Transferable && common.IsZeroAddress(createEntry.DepositAddress) {
		return common.ErrDepositAddressMissing
	}
	for _, tradingPair := range createEntry.TradingPairs {
		if tradingPair.Base != 0 && tradingPair.Quote != 0 {
			return errors.Wrapf(common.ErrBadTradingPairConfiguration, "base id:%v quote id:%v", tradingPair.Base, tradingPair.Quote)
		}

		if tradingPair.Base == 0 && tradingPair.Quote == 0 {
			return errors.Wrapf(common.ErrBadTradingPairConfiguration, "base id:%v quote id:%v", tradingPair.Base, tradingPair.Quote)
		}

		if tradingPair.Base == 0 {
			quoteAsset, err := s.storage.GetAsset(tradingPair.Quote)
			if err != nil {
				return errors.Wrapf(common.ErrQuoteAssetInvalid, "quote id: %v", tradingPair.Quote)
			}
			if !quoteAsset.IsQuote {
				return errors.Wrapf(common.ErrQuoteAssetInvalid, "quote id: %v", tradingPair.Quote)
			}
		}

		if tradingPair.Quote == 0 {
			_, err := s.storage.GetAsset(tradingPair.Base)
			if err != nil {
				return errors.Wrapf(common.ErrBaseAssetInvalid, "base id: %v", tradingPair.Base)
			}

			if !asset.IsQuote {
				return errors.Wrapf(common.ErrQuoteAssetInvalid, "quote id: %v", tradingPair.Quote)
			}
		}
	}
	return nil
}
