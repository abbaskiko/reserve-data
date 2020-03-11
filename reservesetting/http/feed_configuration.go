package http

import (
	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/reservesetting/common"
)

func (s *Server) getFeedConfigurations(c *gin.Context) {
	feedConfigurations, err := s.storage.GetFeedConfigurations()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(feedConfigurations))
}

type feedStatusEntry struct {
	Enabled bool           `json:"enabled" binding:"required"`
	Name    string         `json:"name" binding:"required"`
	SetRate common.SetRate `json:"set_rate" binding:"required"`
}

func (s *Server) updateFeedStatus(c *gin.Context) {
	var fse feedStatusEntry
	if err := c.ShouldBindJSON(&fse); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	feed, err := s.storage.GetFeedConfiguration(fse.Name, fse.SetRate)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	if err := s.storage.UpdateFeedStatus(feed.Name, fse.SetRate, fse.Enabled); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}
