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

type denialsWriter interface {
	Record(ip, reason string)
}

type violatorsBanner interface {
	IsBanned(ip string) bool
	Add(v domain.Violator)
}

type statsWriter interface {
	Update(stat domain.ClientStat)
}

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
	DenialsUC        denialsWriter
	ViolatorsUC      violatorsBanner
	StatsUC          statsWriter
	CacheRepo        domain.CacheRepo
	CacheHandler     *handlers.CacheHandler
	ProxyHandler     http.Handler
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
