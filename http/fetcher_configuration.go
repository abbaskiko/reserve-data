package http

import (
	"strconv"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/gin-gonic/gin"
)

// UpdateFetcherConfiguration update btc fetcher configuration
// and return new configuration
func (s *Server) UpdateFetcherConfiguration(c *gin.Context) {
	var (
		query common.FetcherConfiguration
	)
	postForm, ok := s.Authenticated(c, []string{}, []Permission{ConfirmConfPermission})
	if !ok {
		return
	}
	value := postForm.Get("btc")
	btcConfig, err := strconv.ParseBool(value)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
	}
	query.BTC = btcConfig
	if err := s.app.UpdateFetcherConfiguration(query); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

//GetAllFetcherConfiguration returns all fetcher config
func (s *Server) GetAllFetcherConfiguration(c *gin.Context) {
	_, ok := s.Authenticated(c, []string{}, []Permission{ReadOnlyPermission, ConfigurePermission, ConfirmConfPermission})
	if !ok {
		return
	}
	config, err := s.app.GetAllFetcherConfiguration()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(config))
}
