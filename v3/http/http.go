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

	// because we don't allow to create asset directly, it must go through pending operation
	// so all 'create' operation mean to operate on pending object.
	g.POST("/create-asset", server.createCreateAsset)
	g.GET("/create-asset", server.listCreateAsset)
	g.PUT("/create-asset/:id", server.confirmCreateAsset)
	g.DELETE("/create-asset/:id", server.rejectCreateAsset)

	g.POST("/asset-exchange", server.createAssetExchange)
	g.PUT("/asset-exchange/:id", server.updateAssetExchange)

	g.PUT("/exchange/:id", server.updateExchange)
	g.GET("/exchange", server.getExchanges)
	return server
}
