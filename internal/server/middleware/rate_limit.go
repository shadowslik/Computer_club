package middleware

import (
	"net"
	"net/http"
	"proxy/internal/domain"
	"proxy/pkg/metrics"
	"time"

	"go.uber.org/zap"
)

type clientStatsUpdater interface {
	Update(stat domain.ClientStat)
}

func RateLimitMiddleware(
	uc domain.RateLimitUseCase,
	log *zap.Logger,
	col *metrics.Collector,
	statsUpdater clientStatsUpdater,
	violatorsRepo domain.ViolatorsRepo,
	banDuration time.Duration,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			col.IncConns()
			defer col.DecConns()

			connResult, err := uc.OnConnect(r.Context(), ip)
			if err != nil || !connResult.Allowed {
				log.Warn("connection rejected",
					zap.String("ip", ip),
					zap.String("reason", connResult.Reason),
					zap.String("path", r.URL.Path),
				)
				writeJSONError(w, http.StatusTooManyRequests, connResult.Reason)
				col.RecordRequest(time.Since(start).Milliseconds(), true)
				return
			}
			defer uc.OnDisconnect(r.Context(), ip)

			rules, err := uc.GetRulesForIP(r.Context(), ip)
			if err != nil {
				rules = []domain.RateRule{}
			}

			var maxUploadBps, maxDownloadBps int64 = -1, -1
			for _, rule := range rules {
				if rule.MaxUploadBps >= 0 {
					maxUploadBps = rule.MaxUploadBps
				}
				if rule.MaxDownloadBps >= 0 {
					maxDownloadBps = rule.MaxDownloadBps
				}
			}

			var uploadBytes int64
			if maxUploadBps > 0 && r.Body != nil {
				lim := limiters.getOrUpdate(limiters.upload, ip, maxUploadBps)
				r.Body = &rateLimitedReader{reader: r.Body, limiter: lim, ctx: r.Context()}
			}
			if r.ContentLength > 0 {
				uploadBytes = r.ContentLength
			}

			result, err := uc.Check(r.Context(), ip, uploadBytes, rules)
			if err != nil || !result.Allowed {
				log.Warn("request rate limited",
					zap.String("ip", ip),
					zap.String("reason", result.Reason),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
				)
				writeJSONError(w, http.StatusTooManyRequests, result.Reason)
				col.RecordRequest(time.Since(start).Milliseconds(), true)

				violatorsRepo.Add(domain.Violator{
					IP:          ip,
					Endpoint:    r.URL.Path,
					BlockedAt:   time.Now(),
					BannedUntil: time.Now().Add(banDuration),
				})
				return
			}

			var respWriter responseWriterWithSize
			sr := &statusRecorder{ResponseWriter: w, status: 200}
			if maxDownloadBps > 0 {
				lim := limiters.getOrUpdate(limiters.download, ip, maxDownloadBps)
				respWriter = &rateLimitedWriter{ResponseWriter: sr, limiter: lim, ctx: r.Context()}
			} else {
				respWriter = &trackingWriter{ResponseWriter: sr}
			}

			next.ServeHTTP(respWriter, r)

			downloadBytes := respWriter.Written()
			latencyMs := time.Since(start).Milliseconds()
			uc.TrackDownload(r.Context(), ip, downloadBytes)
			col.RecordTraffic(uploadBytes, downloadBytes)

			lvl := zap.InfoLevel
			if sr.status >= 500 {
				lvl = zap.ErrorLevel
			} else if sr.status >= 400 {
				lvl = zap.WarnLevel
			}
			log.Log(lvl, "http request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("ip", ip),
				zap.Int("status", sr.status),
				zap.Int64("latency_ms", latencyMs),
				zap.Int64("bytes_in", uploadBytes),
				zap.Int64("bytes_out", downloadBytes),
			)

			col.RecordRequest(latencyMs, sr.status >= 500)

			statsUpdater.Update(domain.ClientStat{
				IP:           ip,
				Requests:     1,
				TrafficMB:    float64(uploadBytes+downloadBytes) / (1024 * 1024),
				AvgLatencyMs: float64(latencyMs),
				ErrorRate:    errorRate(sr.status),
				LastSeen:     time.Now().Format("15:04:05"),
			})
		})
	}
}

func errorRate(status int) float64 {
	if status >= 400 {
		return 100.0
	}
	return 0.0
}
