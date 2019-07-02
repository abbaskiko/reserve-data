package http

import (
	"github.com/KyberNetwork/reserve-data/v3/storage"
)

// Server is the HTTP server of token V3.
type Server struct {
	storage storage.Interface
}

// NewServer creates new HTTP server for v3 APIs.
func NewServer(storage storage.Interface) *Server {
	return &Server{storage: storage}
}
