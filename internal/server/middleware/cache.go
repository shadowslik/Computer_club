package middleware

import (
	"bytes"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"proxy/internal/config"
	"proxy/internal/domain"

	"go.uber.org/zap"
)

var hopByHopHeaders = map[string]bool{
	"Connection":          true,
	"Keep-Alive":          true,
	"Proxy-Authenticate":  true,
	"Proxy-Authorization": true,
	"Te":                  true,
	"Trailers":            true,
	"Transfer-Encoding":   true,
	"Upgrade":             true,
}

var corsHeaders = map[string]bool{
	"Access-Control-Allow-Origin":      true,
	"Access-Control-Allow-Methods":     true,
	"Access-Control-Allow-Headers":     true,
	"Access-Control-Allow-Credentials": true,
	"Access-Control-Expose-Headers":    true,
	"Access-Control-Max-Age":           true,
}

type capturingWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func newCapturingWriter(w http.ResponseWriter) *capturingWriter {
	return &capturingWriter{ResponseWriter: w, status: http.StatusOK}
}

func (cw *capturingWriter) WriteHeader(code int) { cw.status = code }

func (cw *capturingWriter) Write(b []byte) (int, error) { return cw.body.Write(b) }

var mutationPrefixes = []struct{ path, cachePrefix string }{
	{"/api/computer_sessions", "GET:"},
	{"/api/clients", "GET:"},
	{"/api/computers", "GET:"},
	{"/ip/whitelist", "GET:"},
	{"/ip/blacklist", "GET:"},
	{"/ip/graylist", "GET:"},
	{"/ratelimit", "GET:"},
}

func CacheMiddleware(repo domain.CacheRepo, cfg *config.Config, log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.Cache.Enabled &&
				(r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete) {
				cw := newCapturingWriter(w)
				next.ServeHTTP(cw, r)
				w.WriteHeader(cw.status)
				w.Write(cw.body.Bytes())
				if cw.status >= 200 && cw.status < 300 {
					for _, mp := range mutationPrefixes {
						if strings.HasPrefix(r.URL.Path, mp.path) {
							n := repo.DeleteByPrefix(mp.cachePrefix + r.Host + mp.path)
							if n > 0 {
								log.Debug("cache invalidated", zap.String("prefix", mp.path), zap.Int("count", n))
							}
							break
						}
					}
				}
				return
			}

			if !cfg.Cache.Enabled || (r.Method != http.MethodGet && r.Method != http.MethodHead) {
				next.ServeHTTP(w, r)
				return
			}

			cc := r.Header.Get("Cache-Control")
			if strings.Contains(cc, "no-cache") || strings.Contains(cc, "no-store") ||
				r.Header.Get("Pragma") == "no-cache" {
				w.Header().Set("X-Cache", "BYPASS")
				next.ServeHTTP(w, r)
				return
			}

			rule := matchRule(r, cfg.Cache.Rules)
			if rule != nil && rule.Enabled != nil && !*rule.Enabled {
				w.Header().Set("X-Cache", "BYPASS")
				next.ServeHTTP(w, r)
				return
			}

			key := buildCacheKey(r)

			if entry, ok := repo.Get(key); ok {
				serveFromCache(w, r, entry)
				log.Debug("cache hit", zap.String("key", key))
				return
			}

			cw := newCapturingWriter(w)
			next.ServeHTTP(cw, r)

			capturedHeaders := filterHeaders(w.Header())

			cacheable, ttl := decideCacheable(cw.status, rule, &cfg.Cache, capturedHeaders)

			bodyLen := int64(cw.body.Len())
			if cacheable && cfg.Cache.MaxBodySize > 0 && bodyLen > cfg.Cache.MaxBodySize {
				cacheable = false
			}
			if cacheable && cfg.Cache.MinBodySize > 0 && bodyLen < cfg.Cache.MinBodySize {
				cacheable = false
			}

			if cacheable {
				tags := ruleTags(rule)
				repo.Set(key, &domain.CacheEntry{
					StatusCode: cw.status,
					Headers:    capturedHeaders,
					Body:       cw.body.Bytes(),
					Tags:       tags,
					CreatedAt:  time.Now(),
					TTL:        ttl,
				})
				w.Header().Set("X-Cache", "MISS")
				log.Debug("cache stored", zap.String("key", key), zap.Duration("ttl", ttl))
			} else {
				w.Header().Set("X-Cache", "BYPASS")
			}

			w.WriteHeader(cw.status)
			w.Write(cw.body.Bytes())
		})
	}
}

func serveFromCache(w http.ResponseWriter, r *http.Request, entry *domain.CacheEntry) {
	for k, vals := range entry.Headers {
		for _, v := range vals {
			w.Header().Add(k, v)
		}
	}
	age := int64(time.Since(entry.CreatedAt).Seconds())
	w.Header().Set("Age", strconv.FormatInt(age, 10))
	w.Header().Set("X-Cache", "HIT")
	w.WriteHeader(entry.StatusCode)
	if r.Method != http.MethodHead {
		w.Write(entry.Body)
	}
}

func decideCacheable(status int, rule *config.CacheRuleConfig, cfg *config.CacheConfig, hdrs http.Header) (bool, time.Duration) {
	upCC := hdrs.Get("Cache-Control")
	if strings.Contains(upCC, "no-store") || strings.Contains(upCC, "private") {
		return false, 0
	}

	var ok bool
	var ttl time.Duration

	switch {
	case status >= 200 && status < 300:
		ok, ttl = cfg.Cache2xx, cfg.DefaultTTL
	case status >= 300 && status < 400:
		ok, ttl = cfg.Cache3xx, cfg.TTL3xx
	case status >= 400 && status < 500:
		ok, ttl = cfg.Cache4xx, cfg.TTL4xx
	case status >= 500:
		ok, ttl = cfg.Cache5xx, cfg.TTL5xx
	}
	if !ok {
		return false, 0
	}
	if rule != nil && rule.TTL > 0 {
		ttl = rule.TTL
	}
	return ttl > 0, ttl
}

func matchRule(r *http.Request, rules []config.CacheRuleConfig) *config.CacheRuleConfig {
	target := r.Host + r.URL.Path
	var best *config.CacheRuleConfig
	bestLen := -1
	for i := range rules {
		rule := &rules[i]
		if strings.HasPrefix(target, rule.Match) && len(rule.Match) > bestLen {
			best = rule
			bestLen = len(rule.Match)
		}
	}
	return best
}

func buildCacheKey(r *http.Request) string {
	q := r.URL.Query()
	qKeys := make([]string, 0, len(q))
	for k := range q {
		qKeys = append(qKeys, k)
	}
	sort.Strings(qKeys)

	var sb strings.Builder
	sb.WriteString(r.Method)
	sb.WriteByte(':')
	sb.WriteString(r.Host)
	sb.WriteString(r.URL.Path)
	if len(qKeys) > 0 {
		sb.WriteByte('?')
		for i, k := range qKeys {
			if i > 0 {
				sb.WriteByte('&')
			}
			sb.WriteString(k)
			sb.WriteByte('=')
			sb.WriteString(q.Get(k))
		}
	}
	return sb.String()
}

func filterHeaders(h http.Header) http.Header {
	out := make(http.Header, len(h))
	for k, v := range h {
		if !hopByHopHeaders[k] && !corsHeaders[k] {
			out[k] = v
		}
	}
	return out
}

func ruleTags(rule *config.CacheRuleConfig) []string {
	if rule == nil || len(rule.Tags) == 0 {
		return nil
	}
	tags := make([]string, len(rule.Tags))
	copy(tags, rule.Tags)
	return tags
}
