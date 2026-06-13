package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/GhanshyamJha05/fifa-cli/internal/config"
	"github.com/GhanshyamJha05/fifa-cli/internal/handlers"
	"github.com/GhanshyamJha05/fifa-cli/internal/service"
)

// Server is the HTTP API server.
type Server struct {
	httpServer *http.Server
	svc        *service.Service
	logger     *slog.Logger
}

// New creates an HTTP server.
func New(cfg *config.Config, svc *service.Service, logger *slog.Logger) *Server {
	addr := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
	timeout := 30 * time.Second
	if cfg.CacheTTL > timeout {
		timeout = cfg.CacheTTL
	}

	handler := handlers.NewRouter(svc, logger, cfg.CORSOrigins)
	return &Server{
		svc:    svc,
		logger: logger,
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      timeout,
			IdleTimeout:       60 * time.Second,
		},
	}
}

// Start begins listening (blocking).
func (s *Server) Start() error {
	s.logger.Info("starting REST API", "addr", s.httpServer.Addr)
	s.svc.RefreshCache(context.Background())
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down REST API")
	return s.httpServer.Shutdown(ctx)
}
