package server

import (
	"net/http"
	"remote-server-api/config"
)

// NewServer creates and configures an HTTP server with the given router and config
func NewServer(handler http.Handler, cfg config.ServerConfig) *http.Server {
	return &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
