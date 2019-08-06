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

		s.r.GET("/prices-version", v3ProxyMW)
		s.r.GET("/prices", v3ProxyMW)
		s.r.GET("/prices/:base/:quote", v3ProxyMW)
		s.r.GET("/getrates", v3ProxyMW)
		s.r.GET("/get-all-rates", v3ProxyMW)

		s.r.GET("/authdata-version", v3ProxyMW)
		s.r.GET("/authdata", v3ProxyMW)
		s.r.GET("/activities", v3ProxyMW)
		s.r.GET("/immediate-pending-activities", v3ProxyMW)
		s.r.GET("/metrics", v3ProxyMW)
		s.r.POST("/metrics", v3ProxyMW)

		s.r.POST("/cancelorder/:exchangeid", v3ProxyMW)
		s.r.POST("/deposit/:exchangeid", v3ProxyMW)
		s.r.POST("/withdraw/:exchangeid", v3ProxyMW)
		s.r.POST("/trade/:exchangeid", v3ProxyMW)
		s.r.POST("/setrates", v3ProxyMW)
		s.r.GET("/tradehistory", v3ProxyMW)

		s.r.GET("/timeserver", v3ProxyMW)

		s.r.GET("/rebalancestatus", v3ProxyMW)
		s.r.POST("/holdrebalance", v3ProxyMW)
		s.r.POST("/enablerebalance", v3ProxyMW)

		s.r.GET("/setratestatus", v3ProxyMW)
		s.r.POST("/holdsetrate", v3ProxyMW)
		s.r.POST("/enablesetrate", v3ProxyMW)

		s.r.POST("/set-stable-token-params", v3ProxyMW)
		s.r.POST("/confirm-stable-token-params", v3ProxyMW)
		s.r.POST("/reject-stable-token-params", v3ProxyMW)
		s.r.GET("/pending-stable-token-params", v3ProxyMW)
		s.r.GET("/stable-token-params", v3ProxyMW)

		s.r.GET("/gold-feed", v3ProxyMW)
		s.r.GET("/btc-feed", v3ProxyMW)
		s.r.POST("/set-feed-configuration", v3ProxyMW)
		s.r.GET("/get-feed-configuration", v3ProxyMW)

		g := s.r.Group("/v3")

		g.GET("/asset/:id", v3ProxyMW)
		g.GET("/asset", v3ProxyMW)
		g.GET("/exchange/:id", v3ProxyMW)
		g.GET("/exchange")

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

		return nil
	}
}
