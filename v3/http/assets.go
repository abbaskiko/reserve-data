package http

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/v3/common"
	"github.com/KyberNetwork/reserve-data/v3/storage"
)

func (s *Server) getAsset(c *gin.Context) {

	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	asset, err := s.storage.GetAsset(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(asset))
}

func (s *Server) getAssets(c *gin.Context) {
	assets, err := s.storage.GetAssets()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
	}
	httputil.ResponseSuccess(c, httputil.WithData(assets))
}

func (s *Server) createAssetExchange(c *gin.Context) {

	var r common.CreateAssetExchange
	err := c.ShouldBindJSON(&r)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		log.Println("failed to bind request", err)
		return
	}

	id, err := s.storage.CreateAssetExchange(r.ExchangeID, r.AssetID,
		r.Symbol, r.DepositAddress, r.MinDeposit, r.WithdrawFee,
		r.TargetRecommended, r.TargetRatio)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithField("id", id))
}

func (s *Server) createPendingAssetExchange(c *gin.Context) {
	var p common.PendingAssetExchange

	if err := c.ShouldBindJSON(&p); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	id, err := s.storage.CreatePendingAssetExchange(p)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithField("id", id))
}

func (s *Server) updateAssetExchange(c *gin.Context) {
	var u storage.UpdateAssetExchangeOpts

	err := c.ShouldBindJSON(&u)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	err = s.storage.UpdateAssetExchange(input.ID, u)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}
