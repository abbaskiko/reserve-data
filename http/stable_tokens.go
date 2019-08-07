package http

import (
	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/http/httputil"
)

// SetStableTokenParams set stable token params
func (s *Server) SetStableTokenParams(c *gin.Context) {
	postForm := c.Request.Form
	value := []byte(postForm.Get("value"))
	if len(value) > maxDataSize {
		httputil.ResponseFailure(c, httputil.WithReason(errDataSizeExceed.Error()))
		return
	}
	err := s.metric.SetStableTokenParams(value)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

// ConfirmStableTokenParams confirm stable token params
func (s *Server) ConfirmStableTokenParams(c *gin.Context) {
	postForm := c.Request.Form
	value := []byte(postForm.Get("value"))
	if len(value) > maxDataSize {
		httputil.ResponseFailure(c, httputil.WithReason(errDataSizeExceed.Error()))
		return
	}
	err := s.metric.ConfirmStableTokenParams(value)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

// RejectStableTokenParams reject stable token params
func (s *Server) RejectStableTokenParams(c *gin.Context) {
	err := s.metric.RemovePendingStableTokenParams()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

// GetPendingStableTokenParams return pending stable token params
func (s *Server) GetPendingStableTokenParams(c *gin.Context) {
	data, err := s.metric.GetPendingStableTokenParams()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(data))
}

// GetStableTokenParams return all stable token params
func (s *Server) GetStableTokenParams(c *gin.Context) {
	data, err := s.metric.GetStableTokenParams()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(data))
}
