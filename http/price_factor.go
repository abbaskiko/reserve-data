package http

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/http/httputil"
)

type getPriceFactorParams struct {
	From uint64 `form:"from" binding:"required"`
	To   uint64 `form:"to" binding:"required"`
}

// getPriceFactor return all metrics
func (s *Server) getPriceFactor(c *gin.Context) {
	log.Printf("get price factor")
	var params getPriceFactorParams
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Printf("cannot bind request parameter, err=%s", err.Error())
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	assets, err := s.settingStorage.GetAssets()
	if err != nil {
		log.Printf("failed to list assets, err=%s", err.Error())
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	data, err := s.metric.GetMetric(assets, params.From, params.To)
	if err != nil {
		log.Printf("failed to get metric, err=%s", err.Error())
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	response := common.MetricResponse{
		Timestamp: common.GetTimepoint(),
	}
	var assetsMetric = make([]common.AssetMetric, 0, len(data))
	for assetID, entry := range data {
		assetsMetric = append(assetsMetric, common.AssetMetric{
			AssetID: assetID,
			Data:    entry,
		})
	}
	response.ReturnTime = common.GetTimepoint()
	response.Data = assetsMetric
	httputil.ResponseSuccess(c, httputil.WithMultipleFields(gin.H{
		"timestamp":  response.Timestamp,
		"returnTime": response.ReturnTime,
		"data":       response.Data,
	}))
}

type setPriceFactorParam struct {
	Timestamp uint64 `json:"timestamp"`
	Data      []struct {
		AssetID uint64  `json:"id"`
		AfpMid  float64 `json:"afp_mid"`
		Spread  float64 `json:"spread"`
	} `json:"data"`
}

// setPriceFactor store metrics into db
func (s *Server) setPriceFactor(c *gin.Context) {
	log.Printf("storing price factor")
	var params setPriceFactorParam
	if err := c.ShouldBindJSON(&params); err != nil {
		log.Printf("cannot bind request parameter, err=%s", err.Error())
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	tokenMetric := map[uint64]common.TokenMetric{}
	for _, e := range params.Data {
		tokenMetric[e.AssetID] = common.TokenMetric{
			AfpMid: e.AfpMid,
			Spread: e.Spread,
		}
	}
	metricEntry := common.MetricEntry{}
	metricEntry.Timestamp = params.Timestamp
	metricEntry.Data = map[uint64]common.TokenMetric{}

	err := s.metric.StoreMetric(&metricEntry, common.GetTimepoint())
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
	} else {
		httputil.ResponseSuccess(c)
	}
}
