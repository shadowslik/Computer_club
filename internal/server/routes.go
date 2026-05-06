package server

import (
	"net/http"
	"proxy/internal/delivery/handlers"
	"proxy/internal/server/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
)

func buildRoutes(deps Deps) http.Handler {
	mux := http.NewServeMux()

	// Swagger UI
	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	// Health & metrics (created inline — just wraps the already-created Collector)
	metricsH := handlers.NewMetricsHandler(deps.Collector)
	mux.HandleFunc("GET /health", metricsH.HandleHealth)
	mux.HandleFunc("GET /metrics", metricsH.HandleMetrics)

	// IP lists
	mux.HandleFunc("/ip/whitelist", deps.WhiteHandler.ListHandler)
	mux.HandleFunc("/ip/blacklist", deps.BlackHandler.ListHandler)
	mux.HandleFunc("/ip/graylist", deps.GrayHandler.ListHandler)
	mux.HandleFunc("GET /ip/check", deps.CheckHandler.CheckIp)
	mux.HandleFunc("GET /ip/denials", deps.DenialsHandler.GetDenials)

	// Rate limiting
	mux.HandleFunc("/ratelimit", deps.RLHandler.ListHandler)
	mux.HandleFunc("/ratelimit/", deps.RLHandler.RuleHandler)
	mux.HandleFunc("GET /ratelimit/violators", deps.ViolatorsHandler.GetViolators)
	mux.HandleFunc("POST /ratelimit/violators/{ip}/unban", deps.ViolatorsHandler.UnbanViolator)

	// Upstream & stats
	mux.HandleFunc("GET /api/v1/upstream", deps.UpstreamHandler.GetStatuses)
	mux.HandleFunc("GET /top-clients", deps.StatsHandler.GetTopClients)

	// Middleware chain (outermost first)
	handler := corsMiddleware(deps.Cfg.Server.CORSAllowedOrigin)(mux)
	handler = middleware.IPFilterMiddleware(
		deps.WhiteHandler.ListUseCase,
		deps.BlackHandler.ListUseCase,
		deps.GrayHandler.ListUseCase,
		deps.Cfg.Lists.DefaultDeny,
		deps.DenialsRepo,
		deps.Log,
	)(handler)
	handler = middleware.RateLimitMiddleware(
		deps.RLUseCase,
		deps.Log,
		deps.Collector,
		deps.StatsRepo,
		deps.ViolatorsRepo,
		deps.Cfg.RateLimit.BanDuration,
	)(handler)

	return handler
}

func corsMiddleware(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
