package http

import (
	"encoding/json"
	"log"

	"github.com/KyberNetwork/reserve-data/common"

	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetGoldData(c *gin.Context) {
	log.Printf("Getting gold data")

	data, err := s.app.GetGoldData(getTimePoint(c, true))
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
	} else {
		httputil.ResponseSuccess(c, httputil.WithData(data))
	}
}

func (s *Server) GetBTCData(c *gin.Context) {
	data, err := s.app.GetBTCData(getTimePoint(c, true))
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

	inputJSON, _ := json.Marshal(input)
	log.Printf("input: %s", inputJSON)

	if err := s.app.UpdateFeedConfiguration(input.Data.Name, input.Data.Enabled); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

func (s *Server) GetFeedConfiguration(c *gin.Context) {
	data, err := s.app.GetFeedConfiguration()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(data))
}
