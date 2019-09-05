package http

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/http/httputil"
)

func (s *Server) getStableTokenParams(c *gin.Context) {
	params, err := s.storage.GetStableTokenParams()
	if err != nil {
		fmt.Printf("failed to get stable token params, err=%v\n", err.Error())
		httputil.ResponseFailure(c, httputil.WithError(err))
	}
	httputil.ResponseSuccess(c, httputil.WithData(params))
}
