package http

import (
	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/reservesetting/common"
	"github.com/KyberNetwork/reserve-data/reservesetting/storage"
)

func (s *Server) getExchanges(c *gin.Context) {
	exhs, err := s.storage.GetExchanges()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(exhs))
}

func (s *Server) getExchange(c *gin.Context) {

	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	exchange, err := s.storage.GetExchange(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(exchange))
}

type exchangeStatusEntry struct {
	Disable bool `form:"disable" binding:"required"`
}

func (s *Server) updateExchangeStatus(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	exchange, err := s.storage.GetExchange(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	var exStatus exchangeStatusEntry
	if err = c.ShouldBindJSON(&exStatus); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	if err := s.storage.UpdateExchange(exchange.ID, storage.UpdateExchangeOpts{
		Disable: common.BoolPointer(exStatus.Disable),
	}); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}
