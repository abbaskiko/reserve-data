package http

import (
	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/http/httputil"
)

func (s *Server) getStableTokenParams(c *gin.Context) {
	params, err := s.storage.GetStableTokenParams()
	if err != nil {
		s.l.Warnw("failed to get stable token params", "err", err)
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(params))
}
