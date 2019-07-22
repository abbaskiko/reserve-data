package http

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/v3/common"
)

func (s *Server) getExchanges(c *gin.Context) {
	exhs, err := s.storage.GetExchanges()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(exhs))
}
func (s *Server) updateExchange(c *gin.Context) {
	var updateExchange common.UpdateExchange
	if err := c.ShouldBindJSON(&updateExchange); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	var input struct {
		ID uint64 `uri:"id" binding:"gte=0"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(errors.New("id uri is required")))
		return
	}

	err := s.storage.UpdateExchange(input.ID, updateExchange)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}
