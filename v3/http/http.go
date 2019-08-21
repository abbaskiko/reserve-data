package http

import (
	"log"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	v1common "github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/v3/common"
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

	g.POST("/setting-change-main", server.createSettingChangeWithType(common.ChangeCatalogMain))
	g.GET("/setting-change-main", server.getSettingChangeWithType(common.ChangeCatalogMain))
	g.GET("/setting-change-main/:id", server.getSettingChange)
	g.PUT("/setting-change-main/:id", server.confirmSettingChange)
	g.DELETE("/setting-change-main/:id", server.rejectSettingChange)

	g.POST("/setting-change-target", server.createSettingChangeWithType(common.ChangeCatalogSetTarget))
	g.GET("/setting-change-target", server.getSettingChangeWithType(common.ChangeCatalogSetTarget))
	g.GET("/setting-change-target/:id", server.getSettingChange)
	g.PUT("/setting-change-target/:id", server.confirmSettingChange)
	g.DELETE("/setting-change-target/:id", server.rejectSettingChange)

	g.POST("/setting-change-pwis", server.createSettingChangeWithType(common.ChangeCatalogSetPWIS))
	g.GET("/setting-change-pwis", server.getSettingChangeWithType(common.ChangeCatalogSetPWIS))
	g.GET("/setting-change-pwis/:id", server.getSettingChange)
	g.PUT("/setting-change-pwis/:id", server.confirmSettingChange)
	g.DELETE("/setting-change-pwis/:id", server.rejectSettingChange)

	g.POST("/setting-change-stable", server.createSettingChangeWithType(common.ChangeCatalogStableToken))
	g.GET("/setting-change-stable", server.getSettingChangeWithType(common.ChangeCatalogStableToken))
	g.GET("/setting-change-stable/:id", server.getSettingChange)
	g.PUT("/setting-change-stable/:id", server.confirmSettingChange)
	g.DELETE("/setting-change-stable/:id", server.rejectSettingChange)

	g.POST("/setting-change-rbquadratic", server.createSettingChangeWithType(common.ChangeCatalogRebalanceQuadratic))
	g.GET("/setting-change-rbquadratic", server.getSettingChangeWithType(common.ChangeCatalogRebalanceQuadratic))
	g.GET("/setting-change-rbquadratic/:id", server.getSettingChange)
	g.PUT("/setting-change-rbquadratic/:id", server.confirmSettingChange)
	g.DELETE("/setting-change-rbquadratic/:id", server.rejectSettingChange)

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
