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
	g.GET("/exchange/:id", server.getExchange)
	g.GET("/exchange", server.getExchanges)

	// because we don't allow to create asset directly, it must go through pending operation
	// so all 'create' operation mean to operate on pending object.

	// API for CreateAsset object
	g.POST("/create-asset", server.createCreateAsset)
	g.GET("/create-asset", server.getCreateAssets)
	g.PUT("/create-asset/:id", server.confirmCreateAsset)
	g.DELETE("/create-asset/:id", server.rejectCreateAsset)

	// API for UpdateAsset object
	g.POST("/update-asset", server.createUpdateAsset)
	g.GET("/update-asset", server.getUpdateAssets)
	g.PUT("/update-asset/:id", server.confirmUpdateAsset)
	g.DELETE("/update-asset/:id", server.rejectUpdateAsset)

	// API for CreateAssetExchange object
	g.POST("/asset-exchange", server.createAssetExchange)

	// API for UpdateAssetExchange object
	g.PUT("/asset-exchange/:id", server.updateAssetExchange)

	// API for UpdateExchange object
	g.PUT("/update-exchange/:id", server.confirmUpdateExchange)
	g.GET("/update-exchange", server.getUpdateExchanges)
	g.GET("/update-exchange/:id", server.getUpdateExchange)
	g.POST("/update-exchange", server.createUpdateExchange)
	g.DELETE("/update-exchange/:id", server.rejectUpdateExchange)
	g.POST("/pending-asset-exchange", server.createPendingAssetExchange)

	g.PUT("/exchange/:id", server.updateExchange)
	g.GET("/exchange", server.getExchanges)
	return server
}
