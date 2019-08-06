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
	server := &Server{
		storage: storage,
		r:       r,
	}
	g := r.Group("/v3")

	g.GET("/asset/:id", server.getAsset)
	g.GET("/asset", server.getAssets)
	g.GET("/exchange/:id", server.getExchange)
	g.GET("/exchange", server.getExchanges)
	g.GET("/trading-pair/:id", server.getTradingPair)

	// because we don't allow to create asset directly, it must go through pending operation
	// so all 'create' operation mean to operate on pending object.

	// API for CreateAsset object
	g.POST("/create-asset", server.createCreateAsset)
	g.GET("/create-asset", server.getCreateAssets)
	g.GET("/create-asset/:id", server.getCreateAsset)
	g.PUT("/create-asset/:id", server.confirmCreateAsset)
	g.DELETE("/create-asset/:id", server.rejectCreateAsset)

	// API for UpdateAsset object
	g.POST("/update-asset", server.createUpdateAsset)
	g.GET("/update-asset", server.getUpdateAssets)
	g.GET("/update-asset/:id", server.getUpdateAsset)
	g.PUT("/update-asset/:id", server.confirmUpdateAsset)
	g.DELETE("/update-asset/:id", server.rejectUpdateAsset)

	// API for ChangeAssetAddress object
	g.POST("/change-asset-address", server.createChangeAssetAddress)
	g.GET("/change-asset-address", server.getChangeAssetAddresses)
	g.GET("/change-asset-address/:id", server.getChangeAssetAddress)
	g.PUT("/change-asset-address/:id", server.confirmChangeAssetAddress)
	g.DELETE("/change-asset-address/:id", server.rejectChangeAssetAddress)

	// API for CreateAssetExchange object
	g.PUT("/create-asset-exchange/:id", server.confirmCreateAssetExchange)
	g.GET("/create-asset-exchange", server.getCreateAssetExchanges)
	g.GET("/create-asset-exchange/:id", server.getCreateAssetExchange)
	g.POST("/create-asset-exchange", server.createCreateAssetExchange)
	g.DELETE("/create-asset-exchange/:id", server.rejectCreateAssetExchange)

	g.PUT("/update-asset-exchange/:id", server.confirmUpdateAssetExchange)
	g.GET("/update-asset-exchange", server.getUpdateAssetExchanges)
	g.GET("/update-asset-exchange/:id", server.getUpdateAssetExchange)
	g.POST("/update-asset-exchange", server.createUpdateAssetExchange)
	g.DELETE("/update-asset-exchange/:id", server.rejectUpdateAssetExchange)

	// API for UpdateExchange object
	g.PUT("/update-exchange/:id", server.confirmUpdateExchange)
	g.GET("/update-exchange", server.getUpdateExchanges)
	g.GET("/update-exchange/:id", server.getUpdateExchange)
	g.POST("/update-exchange", server.createUpdateExchange)
	g.DELETE("/update-exchange/:id", server.rejectUpdateExchange)

	g.PUT("/create-trading-pair/:id", server.confirmCreateTradingPair)
	g.GET("/create-trading-pair", server.getCreateTradingPairs)
	g.GET("/create-trading-pair/:id", server.getCreateTradingPair)
	g.POST("/create-trading-pair", server.createCreateTradingPair)
	g.DELETE("/create-trading-pair/:id", server.rejectCreateTradingPair)

	g.PUT("/update-trading-pair/:id", server.confirmUpdateTradingPair)
	g.GET("/update-trading-pair", server.getUpdateTradingPairs)
	g.GET("/update-trading-pair/:id", server.getUpdateTradingPair)
	g.POST("/update-trading-pair", server.createUpdateTradingPair)
	g.DELETE("/update-trading-pair/:id", server.rejectUpdateTradingPair)

	g.POST("/create-trading-by", server.createCreateTradingBy)
	g.GET("/create-trading-by", server.getCreateTradingBy)
	g.GET("/create-trading-by/:id", server.getCreateTradingBy)
	g.PUT("/create-trading-by/:id", server.confirmCreateTradingBy)
	g.DELETE("/create-trading-by/:id", server.rejectCreateTradingBy)

	return server
}
