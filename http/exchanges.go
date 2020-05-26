package http

import (
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/gin-gonic/gin"
)

// CancelOrderByOrderID cancel an open order on exchanges
func (s *Server) CancelOrderByOrderID(c *gin.Context) {
	postForm, ok := s.Authenticated(c, []string{"order_id", "exchange_id", "symbol"}, []Permission{RebalancePermission})
	if !ok {
		return
	}

	exchangeParam := postForm.Get("exchange_id")
	id := postForm.Get("order_id")
	symbol := postForm.Get("symbol")

	exchange, err := common.GetExchange(exchangeParam)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	s.l.Infof("Cancel order id: %s from %s\n", id, exchange.ID())
	err = s.core.CancelOrderByOrderID(id, symbol, exchange)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

// CancelAllOrders cancel all open orders of a symbol
func (s *Server) CancelAllOrders(c *gin.Context) {
	postForm, ok := s.Authenticated(c, []string{"exchange_id", "symbol"}, []Permission{RebalancePermission})
	if !ok {
		return
	}

	exchangeParam := postForm.Get("exchange_id")
	symbol := postForm.Get("symbol")

	exchange, err := common.GetExchange(exchangeParam)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	err = exchange.CancelAllOrders(symbol)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}
