package http

//Option define initialize behavior for server
type Option func(*Server) error

// WithV1Endpoint set endpoint gateway for V3
func WithV1Endpoint(v1Endpoint string) Option {
	return func(s *Server) error {
		v1ProxyMW, err := newReverseProxyMW(v1Endpoint)
		if err != nil {
			return err
		}
		g := s.r.Group("/v3")

		g.GET("/prices-version", v1ProxyMW)
		g.GET("/prices", v1ProxyMW)
		g.GET("/prices/:base/:quote", v1ProxyMW)
		g.GET("/getrates", v1ProxyMW)
		g.GET("/get-all-rates", v1ProxyMW)

		g.GET("/authdata-version", v1ProxyMW)
		g.GET("/authdata", v1ProxyMW)
		g.GET("/activities", v1ProxyMW)
		g.GET("/immediate-pending-activities", v1ProxyMW)
		g.GET("/price-factor", v1ProxyMW)
		g.POST("/price-factor", v1ProxyMW)

		g.POST("/cancelorder/:exchangeid", v1ProxyMW)
		g.POST("/deposit/:exchangeid", v1ProxyMW)
		g.POST("/withdraw/:exchangeid", v1ProxyMW)
		g.POST("/trade/:exchangeid", v1ProxyMW)
		g.POST("/setrates", v1ProxyMW)
		g.GET("/tradehistory", v1ProxyMW)

		g.GET("/timeserver", v1ProxyMW)

		g.GET("/rebalancestatus", v1ProxyMW)
		g.POST("/holdrebalance", v1ProxyMW)
		g.POST("/enablerebalance", v1ProxyMW)

		g.GET("/setratestatus", v1ProxyMW)
		g.POST("/holdsetrate", v1ProxyMW)
		g.POST("/enablesetrate", v1ProxyMW)

		g.POST("/set-stable-token-params", v1ProxyMW)
		g.POST("/confirm-stable-token-params", v1ProxyMW)
		g.POST("/reject-stable-token-params", v1ProxyMW)
		g.GET("/pending-stable-token-params", v1ProxyMW)
		g.GET("/stable-token-params", v1ProxyMW)

		g.GET("/gold-feed", v1ProxyMW)
		g.GET("/btc-feed", v1ProxyMW)
		g.POST("/set-feed-configuration", v1ProxyMW)
		g.GET("/get-feed-configuration", v1ProxyMW)

		return nil
	}
}

// WithV3Endpoint set endpoint gateway for V3
func WithV3Endpoint(v3Endpoint string) Option {
	return func(s *Server) error {
		v3ProxyMW, err := newReverseProxyMW(v3Endpoint)
		if err != nil {
			return err
		}
		g := s.r.Group("/v3")

		g.GET("/asset/:id", v3ProxyMW)
		g.GET("/asset", v3ProxyMW)
		g.GET("/exchange/:id", v3ProxyMW)
		g.GET("/exchange", v3ProxyMW)

		g.GET("/setting-change", v3ProxyMW)
		g.GET("setting-change/:id", v3ProxyMW)
		g.POST("/setting-change", v3ProxyMW)
		g.PUT("/setting-change/:id", v3ProxyMW)
		g.DELETE("/setting-change/:id", v3ProxyMW)

		return nil
	}
}
