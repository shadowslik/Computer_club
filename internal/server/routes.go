package server

import (
	"log"
	"net/http"
	"proxy/internal/server/middleware"
)

func (s *Server) routes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", s.handleHealth)

	mux.HandleFunc("GET /hello", s.handleHello)

	mux.HandleFunc("/ip/whiteList", s.whiteHandler.ListHandler)

	mux.HandleFunc("/ip/blackList", s.blackHandler.ListHandler)

	mux.HandleFunc("/ip/grayList", s.grayHandler.ListHandler)

	mux.HandleFunc("/ip/check", s.checkHandler.HandlerCheck)

	mux.HandleFunc("/ratelimit", s.rlHandler.ListHandler)

	mux.HandleFunc("/ratelimit/", s.rlHandler.RuleHandler)

	return middleware.RateLimitMiddleware(s.rlUC)(loggingMiddleware(mux))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.log.Debug("health check")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleHello(w http.ResponseWriter, r *http.Request) {
	s.log.Info("hello request", "addr", r.RemoteAddr)
	w.Write([]byte("Hello, World!"))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		log.Printf("%s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
