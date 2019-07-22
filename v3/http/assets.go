package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/v3/common"
	"github.com/KyberNetwork/reserve-data/v3/storage"
)

func (s *Server) getAsset(c *gin.Context) {

	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		responseError(c, http.StatusBadRequest, err.Error())
		return
	}
	asset, err := s.storage.GetAsset(input.ID)
	if err != nil {
		responseWithBackendError(c, err)
		return
	}
	responseData(c, http.StatusOK, asset)
}

func (s *Server) getAssets(c *gin.Context) {
	assets, err := s.storage.GetAssets()
	if err != nil {
		responseWithBackendError(c, err)
	}
	responseData(c, http.StatusOK, assets)
}

func (s *Server) createPendingAsset(c *gin.Context) {
	var createPendingAsset common.CreatePendingAsset

	err := c.ShouldBindJSON(&createPendingAsset)

	if err != nil {
		responseError(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := s.storage.CreatePendingAsset(createPendingAsset)
	if err != nil {
		responseWithBackendError(c, err)
		return
	}
	responseData(c, http.StatusCreated, gin.H{"id": id})
}

func (s *Server) listPendingAsset(c *gin.Context) {
	result, err := s.storage.ListPendingAsset()
	if err != nil {
		responseWithBackendError(c, err)
		return
	}
	responseData(c, http.StatusOK, result)
}

func (s *Server) confirmPendingAsset(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		log.Println(err)
		responseError(c, http.StatusBadRequest, err.Error())
		return
	}
	err := s.storage.ConfirmPendingAsset(input.ID)
	if err != nil {
		responseWithBackendError(c, err)
		return
	}
	responseStatus(c, http.StatusOK, "success")
}

func (s *Server) rejectPendingAsset(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		responseError(c, http.StatusBadRequest, err.Error())
		return
	}
	err := s.storage.RejectPendingAsset(input.ID)
	if err != nil {
		responseWithBackendError(c, err)
		return
	}
	responseStatus(c, http.StatusOK, "success")
}

func (s *Server) createAssetExchange(c *gin.Context) {

	var r common.CreateAssetExchange
	err := c.ShouldBindJSON(&r)
	if err != nil {
		responseError(c, http.StatusBadRequest, "failed to bind request")
		log.Println("failed to bind request", err)
		return
	}

	id, err := s.storage.CreateAssetExchange(r.ExchangeID, r.AssetID,
		r.Symbol, r.DepositAddress, r.MinDeposit, r.WithdrawFee,
		r.TargetRecommended, r.TargetRatio)
	if err != nil {
		responseWithBackendError(c, err)
		return
	}
	responseData(c, http.StatusCreated, gin.H{"id": id})
}

func (s *Server) updateAssetExchange(c *gin.Context) {
	var u storage.UpdateAssetExchangeOpts

	err := c.ShouldBindJSON(&u)
	if err != nil {
		responseError(c, http.StatusBadRequest, "failed to bind request")
		return
	}
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		responseError(c, http.StatusBadRequest, err.Error())
		return
	}

	err = s.storage.UpdateAssetExchange(input.ID, u)
	if err != nil {
		responseWithBackendError(c, err)
		return
	}
	responseStatus(c, http.StatusOK, "success")
}
