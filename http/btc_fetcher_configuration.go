package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/gin-gonic/gin"
)

// UpdateBTCFetcherConfiguration update btc fetcher configuration
// and return new configuration
func (h *HTTPServer) UpdateBTCFetcherConfiguration(c *gin.Context) {
	var (
		query common.BTCFetcherConfigurationRequest
	)
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}
	if err := h.app.UpdateBTCFetcherConfiguration(query); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		query,
	)
}

// GetBTCFetcherConfiguration return BTC fetcher configuration
// true or false
func (h *HTTPServer) GetBTCFetcherConfiguration(c *gin.Context) {
	config, err := h.app.GetBTCFetcherConfiguration()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		config,
	)
}
