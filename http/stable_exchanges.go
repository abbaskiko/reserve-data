package http

import (
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-data/common"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/http/httputil"
)

// GetGoldData return gold data feed
func (s *Server) GetGoldData(c *gin.Context) {
	zap.S().Info("Getting gold data")

	data, err := s.app.GetGoldData(getTimePoint(c, true, s.l))
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
	} else {
		httputil.ResponseSuccess(c, httputil.WithData(data))
	}
}

// GetBTCData return BTC data feed
func (s *Server) GetBTCData(c *gin.Context) {
	data, err := s.app.GetBTCData(getTimePoint(c, true, s.l))
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
	} else {
		httputil.ResponseSuccess(c, httputil.WithData(data))
	}
}

// GetUSDData return BTC data feed
func (s *Server) GetUSDData(c *gin.Context) {
	data, err := s.app.GetUSDData(getTimePoint(c, true, s.l))
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
	} else {
		httputil.ResponseSuccess(c, httputil.WithData(data))
	}
}

// UpdateFeedConfiguration update configuration for feed
func (s *Server) UpdateFeedConfiguration(c *gin.Context) {
	var input common.FeedConfigurationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	if err := s.app.UpdateFeedConfiguration(input.Data.Name, input.Data.Enabled, input.Data.BaseVolatilitySpread); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

// GetFeedConfiguration return feed configuration
func (s *Server) GetFeedConfiguration(c *gin.Context) {
	data, err := s.app.GetFeedConfiguration()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(data))
}
