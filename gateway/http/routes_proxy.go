package http

//Option define initialize behavior for server
type Option func(*Server) error

// WithCoreEndpoint set endpoint gateway for V3
func WithCoreEndpoint(coreEndpoint string) Option {
	return func(s *Server) error {
		coreProxyMW, err := newReverseProxyMW(coreEndpoint)
		if err != nil {
			return err
		}
		g := s.r.Group("/v3")

		g.GET("/prices-version", coreProxyMW)
		g.GET("/prices", coreProxyMW)
		g.GET("/prices/:base/:quote", coreProxyMW)
		g.GET("/getrates", coreProxyMW)
		g.GET("/get-all-rates", coreProxyMW)

		g.GET("/authdata-version", coreProxyMW)
		g.GET("/authdata", coreProxyMW)
		g.GET("/activities", coreProxyMW)
		g.GET("/immediate-pending-activities", coreProxyMW)
		g.GET("/price-factor", coreProxyMW)
		g.POST("/price-factor", coreProxyMW)

		g.POST("/cancelorder/:exchangeid", coreProxyMW)
		g.POST("/deposit/:exchangeid", coreProxyMW)
		g.POST("/withdraw/:exchangeid", coreProxyMW)
		g.POST("/trade/:exchangeid", coreProxyMW)
		g.POST("/setrates", coreProxyMW)
		g.GET("/tradehistory", coreProxyMW)

		g.GET("/timeserver", coreProxyMW)

		g.GET("/rebalancestatus", coreProxyMW)
		g.POST("/holdrebalance", coreProxyMW)
		g.POST("/enablerebalance", coreProxyMW)

		g.GET("/setratestatus", coreProxyMW)
		g.POST("/holdsetrate", coreProxyMW)
		g.POST("/enablesetrate", coreProxyMW)

		g.POST("/set-stable-token-params", coreProxyMW)
		g.POST("/confirm-stable-token-params", coreProxyMW)
		g.POST("/reject-stable-token-params", coreProxyMW)
		g.GET("/pending-stable-token-params", coreProxyMW)
		g.GET("/stable-token-params", coreProxyMW)

		g.GET("/gold-feed", coreProxyMW)
		g.GET("/btc-feed", coreProxyMW)
		g.POST("/set-feed-configuration", coreProxyMW)
		g.GET("/get-feed-configuration", coreProxyMW)

		return nil
	}
}

//WithSettingEndpoint set endpoint gateway for V3
func WithSettingEndpoint(settingEndpoint string) Option {
	return func(s *Server) error {
		settingProxyMW, err := newReverseProxyMW(settingEndpoint)
		if err != nil {
			return err
		}
		g := s.r.Group("/v3")

		g.GET("/asset/:id", settingProxyMW)
		g.GET("/asset", settingProxyMW)
		g.GET("/exchange/:id", settingProxyMW)
		g.GET("/exchange", settingProxyMW)

		g.GET("/setting-change", settingProxyMW)
		g.GET("setting-change/:id", settingProxyMW)
		g.POST("/setting-change", settingProxyMW)
		g.PUT("/setting-change/:id", settingProxyMW)
		g.DELETE("/setting-change/:id", settingProxyMW)

		return nil
	}
}
