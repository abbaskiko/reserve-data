package http

import (
	"log"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/v3/common"
)

func (s *Server) createChangeAssetAddress(c *gin.Context) {
	var changeAssetAddress common.ChangeAssetAddress
	if err := c.ShouldBindJSON(&changeAssetAddress); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	for _, changeAssetAddressEntry := range changeAssetAddress.Assets {
		if err := s.checkChangeAssetAddressParams(changeAssetAddressEntry); err != nil {
			log.Println("error", err)
			httputil.ResponseFailure(c, httputil.WithError(err))
			return
		}
	}

	log.Println("still run here")

	id, err := s.storage.CreateChangeAssetAddress(changeAssetAddress)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithField("id", id))
}

func (s *Server) checkChangeAssetAddressParams(changeAssetAddressEntry common.ChangeAssetAddressEntry) error {
	if _, err := s.storage.GetAsset(changeAssetAddressEntry.ID); err != nil {
		return err
	}
	if !ethereum.IsHexAddress(changeAssetAddressEntry.Address) {
		log.Printf("%s is not a valid ethereum address", changeAssetAddressEntry.Address)
		return common.ErrInvalidAddress
	}
	return nil
}

func (s *Server) getChangeAssetAddress(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		log.Println(err)
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	result, err := s.storage.GetUpdateAsset(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(result))
}

func (s *Server) getChangeAssetAddresses(c *gin.Context) {
	result, err := s.storage.GetChangeAssetAddresses()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(result))
}

func (s *Server) confirmChangeAssetAddress(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		log.Println(err)
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	err := s.storage.ConfirmChangeAssetAddress(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

func (s *Server) rejectChangeAssetAddress(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		log.Println(err)
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	err := s.storage.RejectChangeAssetAddress(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}
