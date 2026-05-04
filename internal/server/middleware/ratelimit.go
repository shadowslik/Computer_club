package middleware

import (
	"context"
	"io"
	"net"
	"net/http"
	"proxy/internal/domain"
	"sync"

	"golang.org/x/time/rate"
)

// limiterEntry хранит лимитер вместе с текущей скоростью
// чтобы пересоздавать при изменении правил
type limiterEntry struct {
	limiter  *rate.Limiter
	speedBps int64
}

type ipLimiters struct {
	mu       sync.Mutex
	download map[string]*limiterEntry
	upload   map[string]*limiterEntry
}

var limiters = &ipLimiters{
	download: make(map[string]*limiterEntry),
	upload:   make(map[string]*limiterEntry),
}

// getOrUpdate — возвращает лимитер для IP
// если скорость изменилась — пересоздаёт лимитер
func (l *ipLimiters) getOrUpdate(m map[string]*limiterEntry, ip string, speedBps int64) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry, ok := m[ip]
	// пересоздаём если лимитера нет или скорость изменилась
	if !ok || entry.speedBps != speedBps {
		lim := rate.NewLimiter(rate.Limit(speedBps), 32*1024)
		m[ip] = &limiterEntry{limiter: lim, speedBps: speedBps}
		return lim
	}
	return entry.limiter
}

func RateLimitMiddleware(uc domain.RateLimitUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			// 2.2.3 — проверяем соединения
			connResult, err := uc.OnConnect(r.Context(), ip)
			if err != nil || !connResult.Allowed {
				writeError(w, connResult.Reason)
				return
			}
			defer uc.OnDisconnect(r.Context(), ip)

			// получаем правила именно для этого IP
			rules, err := uc.GetRulesForIP(r.Context(), ip)
			if err != nil {
				rules = []domain.RateRule{}
			}

			// извлекаем лимиты скорости из правил этого IP
			var maxUploadBps, maxDownloadBps int64 = -1, -1
			for _, rule := range rules {
				if rule.MaxUploadBps >= 0 {
					maxUploadBps = rule.MaxUploadBps
				}
				if rule.MaxDownloadBps >= 0 {
					maxDownloadBps = rule.MaxDownloadBps
				}
			}

			// 2.2.2.2 — upload: оборачиваем r.Body в rateLimitedReader
			var uploadBytes int64
			if maxUploadBps > 0 && r.Body != nil {
				lim := limiters.getOrUpdate(limiters.upload, ip, maxUploadBps)
				r.Body = &rateLimitedReader{
					reader:  r.Body,
					limiter: lim,
					ctx:     r.Context(),
				}
			}
			if r.ContentLength > 0 {
				uploadBytes = r.ContentLength
			}

			// 2.2.1 + 2.2.2.3 + 2.2.3.1 — основные счётчики
			result, err := uc.Check(r.Context(), ip, uploadBytes)
			if err != nil || !result.Allowed {
				writeError(w, result.Reason)
				return
			}

			// 2.2.2.1 — download: оборачиваем ResponseWriter
			var respWriter responseWriterWithSize
			if maxDownloadBps > 0 {
				lim := limiters.getOrUpdate(limiters.download, ip, maxDownloadBps)
				respWriter = &rateLimitedWriter{
					ResponseWriter: w,
					limiter:        lim,
					ctx:            r.Context(),
				}
			} else {
				respWriter = &trackingWriter{ResponseWriter: w}
			}

			next.ServeHTTP(respWriter, r)

			// трекаем download для счётчика дневного трафика
			uc.TrackDownload(r.Context(), ip, respWriter.Written())
		})
	}
}

// responseWriterWithSize — общий интерфейс для обоих writer-ов
type responseWriterWithSize interface {
	http.ResponseWriter
	Written() int64
}

// ── rateLimitedWriter — throttle download ────────────────────────────────────

type rateLimitedWriter struct {
	http.ResponseWriter
	limiter *rate.Limiter
	ctx     context.Context
	written int64
}

func (w *rateLimitedWriter) Write(p []byte) (int, error) {
	total := 0
	chunkSize := 32 * 1024
	for len(p) > 0 {
		chunk := p
		if len(chunk) > chunkSize {
			chunk = p[:chunkSize]
		}
		if err := w.limiter.WaitN(w.ctx, len(chunk)); err != nil {
			return total, err
		}
		n, err := w.ResponseWriter.Write(chunk)
		total += n
		w.written += int64(n)
		if err != nil {
			return total, err
		}
		p = p[n:]
	}
	return total, nil
}

func (w *rateLimitedWriter) Written() int64 { return w.written }

// ── rateLimitedReader — throttle upload ──────────────────────────────────────

type rateLimitedReader struct {
	reader  io.ReadCloser
	limiter *rate.Limiter
	ctx     context.Context
}

func (r *rateLimitedReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if n > 0 {
		if waitErr := r.limiter.WaitN(r.ctx, n); waitErr != nil {
			return n, waitErr
		}
	}
	return n, err
}

func (r *rateLimitedReader) Close() error { return r.reader.Close() }

// ── trackingWriter — без throttle, только счётчик ────────────────────────────

type trackingWriter struct {
	http.ResponseWriter
	written int64
}

func (tw *trackingWriter) Write(b []byte) (int, error) {
	n, err := tw.ResponseWriter.Write(b)
	tw.written += int64(n)
	return n, err
}

func (tw *trackingWriter) Written() int64 { return tw.written }

// ── helpers ───────────────────────────────────────────────────────────────────

func writeError(w http.ResponseWriter, reason string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte(`{"error":"` + reason + `"}`))
}
