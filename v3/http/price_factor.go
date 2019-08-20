package http

import (
	"log"

	"github.com/gin-gonic/gin"

	common2 "github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/v3/common"
)

func (s *Server) setPriceFactor(c *gin.Context) {
	log.Printf("storing price factor")
	var params common.PriceFactorAtTime
	if err := c.ShouldBindJSON(&params); err != nil {
		log.Printf("cannot bind request parameter, err=%s", err.Error())
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	id, err := s.storage.CreatePriceFactor(params)
	if err != nil {
		log.Printf("cannot store price factor, err=%s", err.Error())
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithField("id", id))
}

type getPriceFactorParams struct {
	From uint64 `form:"from" binding:"required"`
	To   uint64 `form:"to" binding:"required"`
}

func convertToPriceFactorResponse(in []common.PriceFactorAtTime) []*common.AssetPriceFactorListResponse {
	var assetToPriceList = map[uint64]*common.AssetPriceFactorListResponse{}
	var res []*common.AssetPriceFactorListResponse
	for _, assetList := range in {
		for _, asset := range assetList.Data {
			var e *common.AssetPriceFactorListResponse
			var ok bool
			if e, ok = assetToPriceList[asset.AssetID]; !ok {
				e = &common.AssetPriceFactorListResponse{
					AssetID: asset.AssetID,
					Data:    nil,
				}
				assetToPriceList[asset.AssetID] = e
				res = append(res, e)
			}
			e.Data = append(e.Data, common.AssetPriceFactorResponse{
				Timestamp: assetList.Timestamp,
				AfpMid:    asset.AfpMid,
				Spread:    asset.Spread,
			})
		}
	}
	return res
}
func (s *Server) getPriceFactor(c *gin.Context) {
	log.Printf("get price factor")
	var params getPriceFactorParams
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Printf("cannot bind request parameter, err=%s", err.Error())
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	store, err := s.storage.GetPriceFactors(params.From, params.To)
	if err != nil {
		log.Printf("cannot get price factor, err=%s", err.Error())
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	data := convertToPriceFactorResponse(store)
	httputil.ResponseSuccess(c, httputil.WithMultipleFields(gin.H{
		"timestamp":  common2.GetTimepoint(),
		"returnTime": common2.GetTimepoint(),
		"data":       data,
	}))
}

func (s *Server) getSetRateStatus(c *gin.Context) {
	status, err := s.storage.GetSetRateStatus()
	if err != nil {
		log.Printf("failed to get set rate control, err:%v\n", err.Error())
		httputil.ResponseFailure(c, httputil.WithError(err))
	}
	httputil.ResponseSuccess(c, httputil.WithData(status))
}

func (s *Server) holdSetRate(c *gin.Context) {
	if err := s.storage.SetSetRateStatus(false); err != nil {
		log.Printf("failed to set set rate control, err:%v\n", err.Error())
		httputil.ResponseFailure(c, httputil.WithError(err))
	}
	httputil.ResponseSuccess(c)
}

func (s *Server) enableSetRate(c *gin.Context) {
	if err := s.storage.SetSetRateStatus(true); err != nil {
		log.Printf("failed to set set rate control, err:%v\n", err.Error())
		httputil.ResponseFailure(c, httputil.WithError(err))
	}
	httputil.ResponseSuccess(c)
}

func (s *Server) getRebalanceStatus(c *gin.Context) {
	status, err := s.storage.GetRebalanceStatus()
	if err != nil {
		log.Printf("failed to get rebalance control, err:%v\n", err.Error())
		httputil.ResponseFailure(c, httputil.WithError(err))
	}
	httputil.ResponseSuccess(c, httputil.WithData(status))
}

func (s *Server) holdRebalance(c *gin.Context) {
	if err := s.storage.SetRebalanceStatus(false); err != nil {
		log.Printf("failed to set rebalance control, err:%v\n", err.Error())
		httputil.ResponseFailure(c, httputil.WithError(err))
	}
	httputil.ResponseSuccess(c)
}

func (s *Server) enableRebalance(c *gin.Context) {
	if err := s.storage.SetRebalanceStatus(true); err != nil {
		log.Printf("failed to set rebalance control, err:%v\n", err.Error())
		httputil.ResponseFailure(c, httputil.WithError(err))
	}
	httputil.ResponseSuccess(c)
}
