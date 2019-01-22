package http

import (
	"net/http"
	"strings"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/gin-gonic/gin"
)

// UpdateFetcherConfiguration update btc fetcher configuration
// and return new configuration
func (h *HTTPServer) UpdateFetcherConfiguration(c *gin.Context) {
	var (
		query common.FetcherConfigurationRequest
	)
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}
	if err := h.app.UpdateFetcherConfiguration(query); err != nil {
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

// GetFetcherConfiguration return BTC fetcher configuration
// true or false
func (h *HTTPServer) GetFetcherConfiguration(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "token name is require"},
		)
		return
	}
	config, err := h.app.GetFetcherConfiguration(strings.ToLower(token))
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		map[string]bool{
			strings.ToLower(token): config,
		},
	)
}

//GetAllFetcherConfiguration returns all fetcher config
func (h *HTTPServer) GetAllFetcherConfiguration(c *gin.Context) {
	config, err := h.app.GetAllFetcherConfiguration()
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
