package http

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/v3/common"
)

func (s *Server) createUpdateAsset(c *gin.Context) {
	var updateAsset common.CreateUpdateAsset

	err := c.ShouldBindJSON(&updateAsset)

	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	id, err := s.storage.CreateUpdateAsset(updateAsset)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithField("id", id))
}

func (s *Server) getUpdateAssets(c *gin.Context) {
	result, err := s.storage.GetUpdateAssets()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(result))
}

func (s *Server) confirmUpdateAsset(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		log.Println(err)
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	err := s.storage.ConfirmUpdateAsset(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

func (s *Server) rejectUpdateAsset(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	err := s.storage.RejectUpdateAsset(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}
