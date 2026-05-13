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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	cascadeRules := make([]domain.CacheCascadeRule, len(cfg.Cache.CascadeRules))
	for i, r := range cfg.Cache.CascadeRules {
		cascadeRules[i] = domain.CacheCascadeRule{TriggerTag: r.TriggerTag, CascadeTags: r.CascadeTags}
	}
	cacheRepo := repository.NewCacheRepo(cfg.Cache.MaxSize, time.Minute, cascadeRules)
	if cfg.Cache.Enabled {
		cacheRepo.StartChangeDetector(5*time.Minute, &http.Client{Timeout: 3 * time.Second})
	}

	whiteUC := usecase.NewListUseCase(whiteRepo)
	blackUC := usecase.NewListUseCase(blackRepo)
	grayUC := usecase.NewListUseCase(grayRepo)
	rlUC := usecase.NewRateLimitService(rlRepo, cfg)
	denialsUC := usecase.NewDenialsUseCase(denialsRepo)
	violatorsUC := usecase.NewViolatorsUseCase(violatorsRepo)
	statsUC := usecase.NewClientStatsUseCase(statsRepo)

	targets := make([]domain.UpstreamTarget, len(cfg.Upstreams))
	for i, u := range cfg.Upstreams {
		targets[i] = domain.UpstreamTarget{ID: u.ID, Name: u.Name, URL: u.URL, Type: u.Type}
	}
	upstreamSvc := usecase.NewUpstreamService(targets)

	whiteH := handlers.NewHandler(whiteUC, log, cacheRepo, "whitelist")
	blackH := handlers.NewHandler(blackUC, log, cacheRepo, "blacklist")
	grayH := handlers.NewHandler(grayUC, log, cacheRepo, "graylist")
	checkH := handlers.NewCheckIpHandler(whiteUC, blackUC, grayUC, log, denialsUC)
	rlH := handlers.NewRateLimitHandler(rlUC, cacheRepo)
	upstreamH := handlers.NewUpstreamHandler(upstreamSvc)
	statsH := handlers.NewClientStatsHandler(statsUC)
	violatorsH := handlers.NewViolatorsHandler(violatorsUC)
	denialsH := handlers.NewDenialsHandler(denialsUC)
	cacheH := handlers.NewCacheHandler(cacheRepo)

	col := metrics.NewCollector()

	var proxyH http.Handler
	if cfg.Server.ProxyTarget != "" {
		proxyH, err = handlers.NewProxyHandler(cfg.Server.ProxyTarget)
		if err != nil {
			log.Fatal("failed to create proxy handler", zap.Error(err))
		}
		log.Info("proxy target configured", zap.String("target", cfg.Server.ProxyTarget))
	}

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
		DenialsUC:        denialsUC,
		ViolatorsUC:      violatorsUC,
		StatsUC:          statsUC,
		CacheRepo:        cacheRepo,
		CacheHandler:     cacheH,
		ProxyHandler:     proxyH,
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := srv.Run(ctx); err != nil {
		log.Fatal("server error", zap.Error(err))
	}
}
