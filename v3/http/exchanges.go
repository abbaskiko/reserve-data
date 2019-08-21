package http

import (
	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/http/httputil"
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
	asset, err := s.storage.GetExchange(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(asset))
}
