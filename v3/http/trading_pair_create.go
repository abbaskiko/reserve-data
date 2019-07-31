package http

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/v3/common"
)

func (s *Server) createCreateTradingPair(c *gin.Context) {
	var createTradingPair common.CreateCreateTradingPair

	err := c.ShouldBindJSON(&createTradingPair)

	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	for index, entry := range createTradingPair.TradingPairs {
		if err = s.checkCreateTradingPairEntry(entry); err != nil {
			httputil.ResponseFailure(c, httputil.WithError(err), httputil.WithField("index", index),
				httputil.WithField("quote", entry.Quote), httputil.WithField("base", entry.Base))
			return
		}
	}

	id, err := s.storage.CreateCreateTradingPair(createTradingPair)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithField("id", id))
}

func (s *Server) checkCreateTradingPairEntry(createEntry common.CreateTradingPairEntry) error {
	base, err := s.storage.GetAsset(createEntry.Base)
	if err != nil {
		return errors.Wrapf(common.ErrBaseAssetInvalid, "base id: %v", createEntry.Base)
	}
	quote, err := s.storage.GetAsset(createEntry.Quote)
	if err != nil {
		return errors.Wrapf(common.ErrBaseAssetInvalid, "quote id: %v", createEntry.Quote)
	}

	if !quote.IsQuote {
		return errors.Wrap(common.ErrQuoteAssetInvalid, "quote asset should have is_quote=true")
	}

	if !isAssetConfigWithExchangeID(base, createEntry.ExchangeID) {
		return errors.Wrap(common.ErrBaseAssetInvalid, "exchange id not found")
	}
	if !isAssetConfigWithExchangeID(quote, createEntry.ExchangeID) {
		return errors.Wrap(common.ErrQuoteAssetInvalid, "exchange id not found")
	}
	return nil
}

func isAssetConfigWithExchangeID(asset common.Asset, exchangeID uint64) bool {
	for _, exchange := range asset.Exchanges {
		if exchange.ExchangeID == exchangeID {
			return true
		}
	}
	return false
}

func (s *Server) getCreateTradingPairs(c *gin.Context) {
	result, err := s.storage.GetCreateTradingPairs()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(result))
}

func (s *Server) getCreateTradingPair(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	result, err := s.storage.GetCreateTradingPair(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(result))
}

func (s *Server) confirmCreateTradingPair(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		log.Println(err)
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	err := s.storage.ConfirmCreateTradingPair(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

func (s *Server) rejectCreateTradingPair(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	err := s.storage.RejectCreateTradingPair(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}
