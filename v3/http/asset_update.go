package http

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

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

	for _, updateEntry := range updateAsset.Assets {
		if err = s.checkUpdateAssetParams(updateEntry); err != nil {
			httputil.ResponseFailure(c, httputil.WithError(err))
			return
		}
	}

	id, err := s.storage.CreateUpdateAsset(updateAsset)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithField("id", id))
}

func (s *Server) checkUpdateAssetParams(updateEntry common.UpdateAssetEntry) error {
	access, err := s.storage.GetAsset(updateEntry.AssetID)
	if err != nil {
		return errors.Wrapf(err, "failed to get asset id: %v from db ", updateEntry.AssetID)
	}

	if updateEntry.Rebalance != nil {
		if *updateEntry.Rebalance {
			if access.RebalanceQuadratic == nil && updateEntry.RebalanceQuadratic == nil {
				return common.ErrRebalanceQuadraticMissing
			}

			if access.Target == nil && updateEntry.Target == nil {
				return common.ErrAssetTargetMissing
			}
		}
	}

	if updateEntry.SetRate != nil {
		if *updateEntry.SetRate != common.SetRateNotSet && access.PWI == nil && updateEntry.PWI == nil {
			return common.ErrPWIMissing
		}
	}

	if updateEntry.Transferable != nil {
		if *updateEntry.Transferable {
			for _, exchange := range access.Exchanges {
				if common.IsZeroAddress(exchange.DepositAddress) {
					return common.ErrDepositAddressMissing
				}
			}
		}
	}
	return nil
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
