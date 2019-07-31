package http

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/v3/common"
)

func (s *Server) createCreateAsset(c *gin.Context) {
	var createAsset common.CreateCreateAsset

	err := c.ShouldBindJSON(&createAsset)

	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	for id, entry := range createAsset.AssetInputs {
		if err = s.checkCreateAssetEntry(entry); err != nil {
			httputil.ResponseFailure(c, httputil.WithError(err), httputil.WithField("index", id))
			return
		}
	}

	id, err := s.storage.CreateCreateAsset(createAsset)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithField("id", id))
}

func (s *Server) checkCreateAssetEntry(createEntry common.CreateAssetEntry) error {
	if createEntry.Rebalance && createEntry.RebalanceQuadratic == nil {
		return common.ErrRebalanceQuadraticMissing
	}

	if createEntry.Rebalance && createEntry.Target == nil {
		return common.ErrAssetTargetMissing
	}

	if createEntry.SetRate != common.SetRateNotSet && createEntry.PWI == nil {
		return common.ErrPWIMissing
	}

	for _, exchange := range createEntry.Exchanges {
		if common.IsZeroAddress(exchange.DepositAddress) && createEntry.Transferable {
			return errors.Wrapf(common.ErrDepositAddressMissing, "exchange %v", exchange.Symbol)
		}

		for _, tradingPair := range exchange.TradingPairs {

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

				if !createEntry.IsQuote {
					return errors.Wrapf(common.ErrQuoteAssetInvalid, "quote id: %v", tradingPair.Quote)
				}
			}
		}
	}

	return nil
}

func (s *Server) getCreateAssets(c *gin.Context) {
	result, err := s.storage.GetCreateAssets()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(result))
}

func (s *Server) getCreateAsset(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		log.Println(err)
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	result, err := s.storage.GetCreateAsset(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(result))
}

func (s *Server) confirmCreateAsset(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		log.Println(err)
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	err := s.storage.ConfirmCreateAsset(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

func (s *Server) rejectCreateAsset(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	err := s.storage.RejectCreateAsset(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}
