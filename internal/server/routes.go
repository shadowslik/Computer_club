package server

import (
	"net/http"
	"proxy/internal/delivery/handlers"
	"proxy/internal/server/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
)

func buildRoutes(deps Deps) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	metricsH := handlers.NewMetricsHandler(deps.Collector)
	mux.HandleFunc("GET /health", metricsH.HandleHealth)
	mux.HandleFunc("GET /metrics", metricsH.HandleMetrics)
	mux.HandleFunc("GET /metrics/history", metricsH.HandleHistory)

	mux.HandleFunc("/ip/whitelist", deps.WhiteHandler.ListHandler)
	mux.HandleFunc("/ip/blacklist", deps.BlackHandler.ListHandler)
	mux.HandleFunc("/ip/graylist", deps.GrayHandler.ListHandler)
	mux.HandleFunc("GET /ip/check", deps.CheckHandler.CheckIp)
	mux.HandleFunc("GET /ip/denials", deps.DenialsHandler.GetDenials)

	mux.HandleFunc("/ratelimit", deps.RLHandler.ListHandler)
	mux.HandleFunc("/ratelimit/", deps.RLHandler.RuleHandler)
	mux.HandleFunc("GET /ratelimit/violators", deps.ViolatorsHandler.GetViolators)
	mux.HandleFunc("POST /ratelimit/violators/{ip}/unban", deps.ViolatorsHandler.UnbanViolator)

	mux.HandleFunc("GET /api/v1/upstream", deps.UpstreamHandler.GetStatuses)
	mux.HandleFunc("GET /top-clients", deps.StatsHandler.GetTopClients)

	if deps.ProxyHandler != nil {
		mux.Handle("/api/", deps.ProxyHandler)
	}

	mux.HandleFunc("GET /cache/stats", deps.CacheHandler.GetStats)
	mux.HandleFunc("DELETE /cache", deps.CacheHandler.Flush)
	mux.HandleFunc("DELETE /cache/key", deps.CacheHandler.DeleteByKey)
	mux.HandleFunc("DELETE /cache/prefix", deps.CacheHandler.DeleteByPrefix)
	mux.HandleFunc("DELETE /cache/regex", deps.CacheHandler.DeleteByRegex)
	mux.HandleFunc("DELETE /cache/tag", deps.CacheHandler.DeleteByTag)
	mux.HandleFunc("DELETE /cache/condition", deps.CacheHandler.DeleteByCondition)
	mux.HandleFunc("POST /cache/events", deps.CacheHandler.PostEvent)

	handler := http.Handler(mux)
	handler = middleware.CacheMiddleware(deps.CacheRepo, deps.Cfg, deps.Log)(handler)
	handler = middleware.AdminAuthMiddleware(deps.Cfg.Auth)(handler)
	handler = middleware.IPFilterMiddleware(
		deps.WhiteHandler.ListUseCase,
		deps.BlackHandler.ListUseCase,
		deps.GrayHandler.ListUseCase,
		deps.Cfg.Lists.DefaultDeny,
		deps.DenialsUC,
		deps.Log,
	)(handler)
	handler = middleware.RateLimitMiddleware(
		deps.RLUseCase,
		deps.Log,
		deps.Collector,
		deps.StatsUC,
		deps.ViolatorsUC,
		deps.Cfg.RateLimit.BanDuration,
	)(handler)
	handler = corsMiddleware(deps.Cfg.Server.CORSAllowedOrigin)(handler)

	return handler
}

func corsMiddleware(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Cache-Control, Pragma, Authorization")
			w.Header().Set("Access-Control-Max-Age", "0")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
