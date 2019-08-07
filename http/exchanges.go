package http

import (
	"fmt"
	"log"
	"math/big"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	v3common "github.com/KyberNetwork/reserve-data/v3/common"
	"github.com/gin-gonic/gin"
)

//TradeRequest request for trade on centralized exchange
type TradeRequest struct {
	Base   string  `json:"base" binding:"required"`
	Quote  string  `json:"quote" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
	Rate   float64 `json:"rate" binding:"required"`
	Type   string  `json:"type"`
}

// Trade create an order in cexs
func (s *Server) Trade(c *gin.Context) {
	var (
		tradeRequest TradeRequest
	)
	exchangeParam := c.Param("exchangeid")
	if err := c.ShouldBindJSON(&tradeRequest); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	exchange, err := common.GetExchange(exchangeParam)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	//TODO: use GetTradingPair method
	var pair v3common.TradingPairSymbols
	pairs, err := s.settingStorage.GetTradingPairs(uint64(exchange.Name()))
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	for _, p := range pairs {
		if p.BaseSymbol == tradeRequest.Base && p.QuoteSymbol == tradeRequest.Quote {
			pair = p
			break
		}
	}
	if tradeRequest.Type != "sell" && tradeRequest.Type != "buy" {
		httputil.ResponseFailure(c, httputil.WithReason(fmt.Sprintf("Trade type of %s is not supported.", tradeRequest.Type)))
		return
	}

	id, done, remaining, finished, err := s.core.Trade(
		exchange, tradeRequest.Type, pair, tradeRequest.Rate,
		tradeRequest.Amount, getTimePoint(c, false))
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithMultipleFields(gin.H{
		"id":        id,
		"done":      done,
		"remaining": remaining,
		"finished":  finished,
	}))
}

// CancelOrder cancel an order from cexs
func (s *Server) CancelOrder(c *gin.Context) {
	postForm := c.Request.Form
	exchangeParam := c.Param("exchangeid")
	id := postForm.Get("order_id")

	exchange, err := common.GetExchange(exchangeParam)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	log.Printf("Cancel order id: %s from %s\n", id, exchange.ID())
	activityID, err := common.StringToActivityID(id)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	err = s.core.CancelOrder(activityID, exchange)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

// WithdrawRequest request to withdraw
type WithdrawRequest struct {
	Asset  uint64 `json:"asset" binding:"required"`
	Amount string `json:"amount" binding:"required"`
}

// Withdraw asset to reserve from cex
func (s *Server) Withdraw(c *gin.Context) {
	var (
		r WithdrawRequest
	)
	exchangeParam := c.Param("exchangeid")
	if err := c.ShouldBindJSON(&r); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	exchange, err := common.GetExchange(exchangeParam)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	asset, err := s.settingStorage.GetAsset(r.Asset)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	amount, ok := big.NewInt(0).SetString(r.Amount, 10)
	if !ok {
		httputil.ResponseFailure(c, httputil.WithError(fmt.Errorf("cannot parse amount: %s", r.Amount)))
		return
	}
	log.Printf("Withdraw %s %d from %s\n", amount.Text(10), asset.ID, exchange.ID())
	id, err := s.core.Withdraw(exchange, asset, amount, getTimePoint(c, false))
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithField("id", id))
}

//DepositRequest deposit to centralized exchange
type DepositRequest struct {
	AssetID uint64 `json:"asset" binding:"required"`
	Amount  string `json:"amount" binding:"required"`
}

// Deposit asset into cex
func (s *Server) Deposit(c *gin.Context) {
	var (
		r DepositRequest
	)
	exchangeParam := c.Param("exchangeid")
	if err := c.ShouldBindJSON(&r); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	exchange, err := common.GetExchange(exchangeParam)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	asset, err := s.settingStorage.GetAsset(r.AssetID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	amount, ok := big.NewInt(0).SetString(r.Amount, 10)
	if !ok {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	log.Printf("Depositing %s %d to %s\n", amount.Text(10), asset.ID, exchange.ID())
	id, err := s.core.Deposit(exchange, asset, amount, getTimePoint(c, false))
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithField("id", id))
}
