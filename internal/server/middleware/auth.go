package middleware

import (
	"net/http"
	"strings"

	"proxy/internal/config"
)

func AdminAuthMiddleware(cfg config.AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			if !isAdminPath(r.URL.Path, cfg.AdminPaths) {
				next.ServeHTTP(w, r)
				return
			}

			token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if cfg.AdminKey == "" || token != cfg.AdminKey {
				writeJSONError(w, http.StatusForbidden, "доступ запрещён: требуется admin ключ")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isAdminPath(path string, adminPaths []string) bool {
	for _, p := range adminPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
