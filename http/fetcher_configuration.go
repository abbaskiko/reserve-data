package http

import (
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/gin-gonic/gin"
)

// UpdateFetcherConfiguration update btc fetcher configuration
// and return new configuration
func (h *HTTPServer) UpdateFetcherConfiguration(c *gin.Context) {
	var (
		query common.FetcherConfiguration
	)
	if err := c.ShouldBindJSON(&query); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	if err := h.app.UpdateFetcherConfiguration(query); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

//GetAllFetcherConfiguration returns all fetcher config
func (h *HTTPServer) GetAllFetcherConfiguration(c *gin.Context) {
	config, err := h.app.GetAllFetcherConfiguration()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(config))
}
