package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"proxy/internal/config"
	"proxy/internal/delivery/handlers"
	"proxy/internal/usecase"
	"time"
)

type Server struct {
	httpServer *http.Server
	log        *slog.Logger
	cfg        *config.Config

	whiteHandler *handlers.Handler
	blackHandler *handlers.Handler
	grayHandler  *handlers.Handler
	checkHandler *handlers.CheckIpHandler
	rlHandler    *handlers.RateLimitHandler

	rlUC *usecase.RateLimitService
}

func NewServer(
	cfg *config.Config,
	log *slog.Logger,
	white, gray, black *handlers.Handler,
	check *handlers.CheckIpHandler,
	rlHandler *handlers.RateLimitHandler,
	rlUC *usecase.RateLimitService,
) *Server {

	server := &Server{
		log: log,
		cfg: cfg,

		whiteHandler: white,
		blackHandler: black,
		grayHandler:  gray,
		checkHandler: check,

		rlHandler: rlHandler,
		rlUC:      rlUC,
	}

	server.httpServer = &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      server.routes(),
		ReadTimeout:  0,
		WriteTimeout: 0,
		IdleTimeout:  60 * time.Second,
	}

	return server
}

func (s *Server) Run(ctx context.Context) error {
	s.log.Info("starting server", "port", s.cfg.Port)

	go func() {
		<-ctx.Done()
		s.log.Info("shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.httpServer.Shutdown(shutdownCtx)
	}()

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
