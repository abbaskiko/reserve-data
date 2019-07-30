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

		s.r.GET("/asset/:id", v3ProxyMW)
		return nil
	}
}
