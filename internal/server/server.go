package server

import (
	"context"
	"fmt"
	"net/http"
	"proxy/internal/config"
	"proxy/internal/delivery/handlers"
	"proxy/internal/domain"
	"proxy/pkg/metrics"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	log        *zap.Logger
}

type Deps struct {
	Cfg              *config.Config
	Log              *zap.Logger
	Collector        *metrics.Collector
	WhiteHandler     *handlers.Handler
	BlackHandler     *handlers.Handler
	GrayHandler      *handlers.Handler
	CheckHandler     *handlers.CheckIpHandler
	RLHandler        *handlers.RateLimitHandler
	RLUseCase        domain.RateLimitUseCase
	UpstreamHandler  *handlers.UpstreamHandler
	StatsHandler     *handlers.ClientStatsHandler
	ViolatorsHandler *handlers.ViolatorsHandler
	DenialsHandler   *handlers.DenialsHandler
	DenialsRepo      domain.DenialsRepo
	ViolatorsRepo    domain.ViolatorsRepo
	StatsRepo        domain.ClientStatsRepo
}

func NewServer(deps Deps) *Server {
	s := &Server{log: deps.Log}
	s.httpServer = &http.Server{
		Addr:         ":" + deps.Cfg.Server.Port,
		Handler:      buildRoutes(deps),
		ReadTimeout:  deps.Cfg.Server.ReadTimeout,
		WriteTimeout: deps.Cfg.Server.WriteTimeout,
		IdleTimeout:  deps.Cfg.Server.IdleTimeout,
	}
	return s
}

func (s *Server) Run(ctx context.Context) error {
	s.log.Info("starting server", zap.String("addr", s.httpServer.Addr))
	go func() {
		<-ctx.Done()
		s.log.Info("shutting down server")
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.httpServer.Shutdown(shutCtx)
	}()
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}
