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
