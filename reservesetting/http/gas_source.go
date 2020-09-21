package http

import (
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/reservesetting/common"
	"github.com/gin-gonic/gin"
)

func (s *Server) setPreferGasSource(c *gin.Context) {
	name := c.Request.FormValue("name")
	if name == "" {
		httputil.ResponseFailure(c, httputil.WithReason("name is required"))
	}
	if err := s.storage.SetPreferGasSource(common.PreferGasSource{Name: name}); err != nil {
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
