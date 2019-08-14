package http

//Option define initialize behavior for server
type Option func(*Server) error

// WithV3Endpoint set endpoint gateway for V3
func WithV3Endpoint(v3Endpoint string) Option {
	return func(s *Server) error {
		v3ProxyMW, err := newReverseProxyMW(v3Endpoint)
		if err != nil {
			return err
		}
		g := s.r.Group("/v3")

		g.GET("/prices-version", v3ProxyMW)
		g.GET("/prices", v3ProxyMW)
		g.GET("/prices/:base/:quote", v3ProxyMW)
		g.GET("/getrates", v3ProxyMW)
		g.GET("/get-all-rates", v3ProxyMW)

		g.GET("/authdata-version", v3ProxyMW)
		g.GET("/authdata", v3ProxyMW)
		g.GET("/activities", v3ProxyMW)
		g.GET("/immediate-pending-activities", v3ProxyMW)
		g.GET("/price-factor", v3ProxyMW)
		g.POST("/price-factor", v3ProxyMW)

		g.POST("/cancelorder/:exchangeid", v3ProxyMW)
		g.POST("/deposit/:exchangeid", v3ProxyMW)
		g.POST("/withdraw/:exchangeid", v3ProxyMW)
		g.POST("/trade/:exchangeid", v3ProxyMW)
		g.POST("/setrates", v3ProxyMW)
		g.GET("/tradehistory", v3ProxyMW)

		g.GET("/timeserver", v3ProxyMW)

		g.GET("/rebalancestatus", v3ProxyMW)
		g.POST("/holdrebalance", v3ProxyMW)
		g.POST("/enablerebalance", v3ProxyMW)

		g.GET("/setratestatus", v3ProxyMW)
		g.POST("/holdsetrate", v3ProxyMW)
		g.POST("/enablesetrate", v3ProxyMW)

		g.POST("/set-stable-token-params", v3ProxyMW)
		g.POST("/confirm-stable-token-params", v3ProxyMW)
		g.POST("/reject-stable-token-params", v3ProxyMW)
		g.GET("/pending-stable-token-params", v3ProxyMW)
		g.GET("/stable-token-params", v3ProxyMW)

		g.GET("/gold-feed", v3ProxyMW)
		g.GET("/btc-feed", v3ProxyMW)
		g.POST("/set-feed-configuration", v3ProxyMW)
		g.GET("/get-feed-configuration", v3ProxyMW)

		g.GET("/asset/:id", v3ProxyMW)
		g.GET("/asset", v3ProxyMW)
		g.GET("/exchange/:id", v3ProxyMW)
		g.GET("/exchange", v3ProxyMW)

		g.GET("/create-asset", v3ProxyMW)
		g.GET("/create-asset/:id", v3ProxyMW)
		g.POST("/create-asset", v3ProxyMW)
		g.PUT("/create-asset/:id", v3ProxyMW)
		g.DELETE("/create-asset/:id", v3ProxyMW)

		g.GET("/update-asset", v3ProxyMW)
		g.GET("/update-asset/:id", v3ProxyMW)
		g.POST("/update-asset", v3ProxyMW)
		g.PUT("/update-asset/:id", v3ProxyMW)
		g.DELETE("/update-asset/:id", v3ProxyMW)

		g.POST("/change-asset-address", v3ProxyMW)
		g.GET("/change-asset-address", v3ProxyMW)
		g.GET("/change-asset-address/:id", v3ProxyMW)
		g.PUT("/change-asset-address/:id", v3ProxyMW)
		g.DELETE("/change-asset-address/:id", v3ProxyMW)

		g.GET("/create-asset-exchange", v3ProxyMW)
		g.GET("/create-asset-exchange/:id", v3ProxyMW)
		g.POST("/create-asset-exchange", v3ProxyMW)
		g.PUT("/create-asset-exchange/:id", v3ProxyMW)
		g.DELETE("/create-asset-exchange/:id", v3ProxyMW)

		g.GET("/update-asset-exchange", v3ProxyMW)
		g.GET("/update-asset-exchange/:id", v3ProxyMW)
		g.POST("/update-asset-exchange", v3ProxyMW)
		g.PUT("/update-asset-exchange/:id", v3ProxyMW)
		g.DELETE("/update-asset-exchange/:id", v3ProxyMW)

		g.GET("/update-exchange", v3ProxyMW)
		g.GET("/update-exchange/:id", v3ProxyMW)
		g.POST("/update-exchange", v3ProxyMW)
		g.PUT("/update-exchange/:id", v3ProxyMW)
		g.DELETE("/update-exchange/:id", v3ProxyMW)

		g.GET("/create-trading-pair", v3ProxyMW)
		g.GET("/create-trading-pair/:id", v3ProxyMW)
		g.POST("/create-trading-pair", v3ProxyMW)
		g.PUT("/create-trading-pair/:id", v3ProxyMW)
		g.DELETE("/create-trading-pair/:id", v3ProxyMW)

		g.GET("/update-trading-pair", v3ProxyMW)
		g.GET("/update-trading-pair/:id", v3ProxyMW)
		g.POST("/update-trading-pair", v3ProxyMW)
		g.PUT("/update-trading-pair/:id", v3ProxyMW)
		g.DELETE("/update-trading-pair/:id", v3ProxyMW)

		g.GET("/setting-change", v3ProxyMW)
		g.GET("setting-change/:id", v3ProxyMW)
		g.POST("/setting-change", v3ProxyMW)
		g.PUT("/setting-change/:id", v3ProxyMW)
		g.DELETE("/setting-change/:id", v3ProxyMW)

		return nil
	}
}
