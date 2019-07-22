package http

import (
	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/v3/storage"
)

// Server is the HTTP server of token V3.
type Server struct {
	storage storage.Interface
	r       *gin.Engine
}

// NewServer creates new HTTP server for v3 APIs.
func NewServer(storage storage.Interface, r *gin.Engine) *Server {
	if r == nil {
		r = gin.Default()
	}
	server := &Server{storage: storage, r: r}
	g := r.Group("/v3")

	g.GET("/asset/:id", server.getAsset)
	g.GET("/asset", server.getAssets)

	g.POST("/pending-asset", server.createPendingAsset)
	g.GET("/pending-asset", server.listPendingAsset)
	g.PUT("/pending-asset/:id", server.confirmPendingAsset)
	g.DELETE("/pending-asset/:id", server.rejectPendingAsset)

	g.POST("/asset-exchange", server.createAssetExchange)
	g.PUT("/asset-exchange/:id", server.updateAssetExchange)

	g.PUT("/exchange/:id", server.updateExchange)
	g.GET("/exchange", server.getExchanges)
	return server
}
