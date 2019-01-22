package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/gin-gonic/gin"
)

// UpdateFetcherConfiguration update btc fetcher configuration
// and return new configuration
func (h *HTTPServer) UpdateFetcherConfiguration(c *gin.Context) {
	var (
		query common.FetcherConfiguration
	)
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"error": err.Error()},
		)
		return
	}
	if err := h.app.UpdateFetcherConfiguration(query); err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		query,
	)
}

//GetAllFetcherConfiguration returns all fetcher config
func (h *HTTPServer) GetAllFetcherConfiguration(c *gin.Context) {
	config, err := h.app.GetAllFetcherConfiguration()
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		config,
	)
}
