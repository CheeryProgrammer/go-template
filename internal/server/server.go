package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/YOUR_ORG/myapp/internal/config"
	"github.com/YOUR_ORG/myapp/internal/handler"
	"github.com/YOUR_ORG/myapp/internal/store"
)

// Server wraps the HTTP server and manages its lifecycle.
type Server struct {
	cfg  config.Config
	http *http.Server
	stop func()
}

// New creates a Server from cfg, optionally connecting to PostgreSQL if
// cfg.DatabaseURL is set.
func New(cfg config.Config) (*Server, error) {
	s := &Server{cfg: cfg, stop: func() {}}

	var st store.Store
	if cfg.DatabaseURL != "" {
		pg, err := store.NewPostgres(cfg.DatabaseURL)
		if err != nil {
			return nil, fmt.Errorf("connect to database: %w", err)
		}
		st = pg
		s.stop = pg.Close
		slog.Info("connected to database")
	} else {
		slog.Warn("DATABASE_URL not set, starting without database")
	}

	s.http = &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      handler.New(st),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return s, nil
}

// Run starts the HTTP server and blocks until ctx is cancelled or a fatal error
// occurs. It performs a graceful shutdown before returning.
func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		slog.Info("server listening", "addr", s.http.Addr)
		if err := s.http.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		defer s.stop()
		if err := s.http.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
		return nil
	}
}
