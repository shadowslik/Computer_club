// @title           Proxy API
// @version         1.0
// @description     HTTP proxy с управлением доступом по IP (whitelist/blacklist/graylist), rate limiting (RPS/RPM/RPH/RPD/bandwidth/соединения) и мониторингом в реальном времени.
// @contact.name    API Support
// @contact.email   artiomk00o@gmail.com
// @license.name    MIT
// @host            localhost:8081
// @BasePath        /
// @schemes         http

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "proxy/docs"
	"proxy/internal/config"
	"proxy/internal/delivery/handlers"
	"proxy/internal/domain"
	"proxy/internal/repository"
	"proxy/internal/server"
	"proxy/internal/usecase"
	"proxy/pkg/logger"
	"proxy/pkg/metrics"

	"go.uber.org/zap"
)

func main() {
	cfgPath := os.Getenv("CONFIG_FILE")
	if cfgPath == "" {
		cfgPath = "configs/config.yaml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	log, err := logger.New(cfg.Log.Level)
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}
	defer log.Sync()

	log.Info("config loaded", zap.String("path", cfgPath))

	// --- Repositories ---
	whiteRepo, err := repository.NewIpRepoImpl(cfg.Lists.WhitelistFile, cfg.Lists.CacheSize, cfg.Lists.CacheTTL, log)
	if err != nil {
		log.Fatal("failed to load whitelist", zap.Error(err))
	}
	blackRepo, err := repository.NewIpRepoImpl(cfg.Lists.BlacklistFile, cfg.Lists.CacheSize, cfg.Lists.CacheTTL, log)
	if err != nil {
		log.Fatal("failed to load blacklist", zap.Error(err))
	}
	grayRepo, err := repository.NewIpRepoImpl(cfg.Lists.GraylistFile, cfg.Lists.CacheSize, cfg.Lists.CacheTTL, log)
	if err != nil {
		log.Fatal("failed to load graylist", zap.Error(err))
	}
	rlRepo, err := repository.NewRateLimitStore(cfg.Lists.RateLimitFile)
	if err != nil {
		log.Fatal("failed to load ratelimit store", zap.Error(err))
	}

	denialsRepo := repository.NewDenialsRepo()
	violatorsRepo := repository.NewViolatorsRepo()
	statsRepo := repository.NewClientStatsRepo()

	// --- Use Cases ---
	whiteUC := usecase.NewHTTPUseCase(whiteRepo)
	blackUC := usecase.NewHTTPUseCase(blackRepo)
	grayUC := usecase.NewHTTPUseCase(grayRepo)
	rlUC := usecase.NewRateLimitService(rlRepo, cfg)
	denialsUC := usecase.NewDenialsUseCase(denialsRepo)
	violatorsUC := usecase.NewViolatorsUseCase(violatorsRepo)
	statsUC := usecase.NewClientStatsUseCase(statsRepo)

	targets := make([]domain.UpstreamTarget, len(cfg.Upstreams))
	for i, u := range cfg.Upstreams {
		targets[i] = domain.UpstreamTarget{ID: u.ID, Name: u.Name, URL: u.URL, Type: u.Type}
	}
	upstreamSvc := usecase.NewUpstreamService(targets)

	// --- Handlers ---
	whiteH := handlers.NewHandler(whiteUC, log)
	blackH := handlers.NewHandler(blackUC, log)
	grayH := handlers.NewHandler(grayUC, log)
	checkH := handlers.NewCheckIpHandler(whiteUC, blackUC, grayUC, log, denialsRepo)
	rlH := handlers.NewRateLimitHandler(rlUC)
	upstreamH := handlers.NewUpstreamHandler(upstreamSvc)
	statsH := handlers.NewClientStatsHandler(statsUC)
	violatorsH := handlers.NewViolatorsHandler(violatorsUC)
	denialsH := handlers.NewDenialsHandler(denialsUC)

	col := metrics.NewCollector()

	srv := server.NewServer(server.Deps{
		Cfg:              cfg,
		Log:              log,
		Collector:        col,
		WhiteHandler:     whiteH,
		BlackHandler:     blackH,
		GrayHandler:      grayH,
		CheckHandler:     checkH,
		RLHandler:        rlH,
		RLUseCase:        rlUC,
		UpstreamHandler:  upstreamH,
		StatsHandler:     statsH,
		ViolatorsHandler: violatorsH,
		DenialsHandler:   denialsH,
		DenialsRepo:      denialsRepo,
		ViolatorsRepo:    violatorsRepo,
		StatsRepo:        statsRepo,
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Info("server starting", zap.String("port", cfg.Server.Port))
	if err := srv.Run(ctx); err != nil {
		log.Fatal("server error", zap.Error(err))
	}
}
