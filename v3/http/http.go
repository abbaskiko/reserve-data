package http

import (
	"log"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	v1common "github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/v3/storage"
)

// Server is the HTTP server of token V3.
type Server struct {
	storage            storage.Interface
	r                  *gin.Engine
	host               string
	supportedExchanges map[v1common.ExchangeID]v1common.LiveExchange
}

// NewServer creates new HTTP server for v3 APIs.
func NewServer(storage storage.Interface, host string, supportedExchanges map[v1common.ExchangeID]v1common.LiveExchange) *Server {
	r := gin.Default()
	server := &Server{
		storage:            storage,
		r:                  r,
		host:               host,
		supportedExchanges: supportedExchanges,
	}
	g := r.Group("/v3")

	g.GET("/asset/:id", server.getAsset)
	g.GET("/asset", server.getAssets)
	g.GET("/exchange/:id", server.getExchange)
	g.GET("/exchange", server.getExchanges)
	g.GET("/trading-pair/:id", server.getTradingPair)

	// because we don't allow to create asset directly, it must go through pending operation
	// so all 'create' operation mean to operate on pending object.

	g.POST("/setting-change", server.createSettingChange)
	g.GET("/setting-change", server.getSettingChanges)
	g.GET("/setting-change/:id", server.getSettingChange)
	g.PUT("/setting-change/:id", server.confirmSettingChange)
	g.DELETE("/setting-change/:id", server.rejectSettingChange)

	g.GET("/price-factor", server.getPriceFactor)
	g.POST("/price-factor", server.setPriceFactor)

	g.GET("/set-rate-status", server.getSetRateStatus)
	g.POST("/hold-set-rate", server.holdSetRate)
	g.POST("/enable-set-rate", server.enableSetRate)

	g.GET("/rebalance-status", server.getRebalanceStatus)
	g.POST("/hold-rebalance", server.holdRebalance)
	g.POST("/enable-rebalance", server.enableRebalance)

	return server
}

// EnableProfiler enable profiler on path "/debug/pprof"
func (s *Server) EnableProfiler() {
	pprof.Register(s.r)
}

// Run the server
func (s *Server) Run() {
	if err := s.r.Run(s.host); err != nil {
		log.Panic(err)
	}
}
