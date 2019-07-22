package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/v3/common"
)

func (s *Server) getExchanges(c *gin.Context) {
	exhs, err := s.storage.GetExchanges()
	if err != nil {
		responseWithBackendError(c, err)
		return
	}
	responseData(c, http.StatusOK, exhs)
}
func (s *Server) updateExchange(c *gin.Context) {
	var updateExchange common.UpdateExchange
	if err := c.ShouldBindJSON(&updateExchange); err != nil {
		responseError(c, http.StatusBadRequest, "failed to bind request")
		return
	}
	var input struct {
		ID uint64 `uri:"id" binding:"gte=0"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		responseError(c, http.StatusBadRequest, "id uri is required")
		return
	}

	err := s.storage.UpdateExchange(input.ID, updateExchange)
	if err != nil {
		responseWithBackendError(c, err)
		return
	}
	responseStatus(c, http.StatusOK, "success")
}
