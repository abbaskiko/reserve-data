package http

import (
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/reservesetting/common"
	"github.com/gin-gonic/gin"
)

func (s *Server) setPreferGasSource(c *gin.Context) {
	var input common.PreferGasSource
	if err := c.BindJSON(input); err != nil {
		httputil.ResponseFailure(c, httputil.WithReason(err.Error()))
		return
	}
	if err := s.storage.SetPreferGasSource(input); err != nil {
		httputil.ResponseFailure(c, httputil.WithReason(err.Error()))
		return
	}
	httputil.ResponseSuccess(c)
}

func (s *Server) getPreferGasSource(c *gin.Context) {
	preferredGasSource, err := s.storage.GetPreferGasSource()
	if err != nil {
		s.l.Errorw("failed to get prefered gas source", "err", err)
		httputil.ResponseFailure(c, httputil.WithReason(err.Error()))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(preferredGasSource))
}
