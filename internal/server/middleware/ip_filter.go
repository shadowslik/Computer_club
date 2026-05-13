package middleware

import (
	"context"
	"net"
	"net/http"
	"proxy/pkg/utils"
	"time"

	"go.uber.org/zap"
)

type ipChecker interface {
	IsAllowed(ctx context.Context, ip string) (bool, error)
}

type denialsRecorder interface {
	Record(ip, reason string)
}

func IPFilterMiddleware(
	whiteUC, blackUC, grayUC ipChecker,
	defaultDeny bool,
	denials denialsRecorder,
	log *zap.Logger,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}
			if ip == "" {
				ip = utils.GetRealIP(r)
			}

			ctx := r.Context()

			if blacklisted, _ := blackUC.IsAllowed(ctx, ip); blacklisted {
				log.Warn("access denied: blacklist",
					zap.String("ip", ip),
					zap.String("path", r.URL.Path),
					zap.String("reason", "blacklist"),
					zap.Time("time", time.Now()),
				)
				denials.Record(ip, "blacklist")
				writeJSONError(w, http.StatusForbidden, "IP находится в чёрном списке")
				return
			}

			if graylisted, _ := grayUC.IsAllowed(ctx, ip); graylisted {
				log.Warn("access requires CAPTCHA: graylist",
					zap.String("ip", ip),
					zap.String("path", r.URL.Path),
					zap.String("reason", "graylist"),
					zap.Time("time", time.Now()),
				)
				denials.Record(ip, "graylist")
				w.Header().Set("X-CAPTCHA-Required", "true")
				writeJSONError(w, http.StatusTooManyRequests, "требуется CAPTCHA")
				return
			}

			if whitelisted, _ := whiteUC.IsAllowed(ctx, ip); whitelisted {
				log.Info("access allowed: whitelist",
					zap.String("ip", ip),
					zap.String("path", r.URL.Path),
					zap.Time("time", time.Now()),
				)
				next.ServeHTTP(w, r)
				return
			}

			if defaultDeny {
				log.Warn("access denied: not in whitelist",
					zap.String("ip", ip),
					zap.String("path", r.URL.Path),
					zap.String("reason", "not_in_whitelist"),
					zap.Time("time", time.Now()),
				)
				denials.Record(ip, "not_in_whitelist")
				writeJSONError(w, http.StatusForbidden, "доступ запрещён")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
