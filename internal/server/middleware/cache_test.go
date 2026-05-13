package middleware

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"proxy/internal/config"
	"proxy/internal/domain"
)

type mockCacheRepoMW struct {
	entries map[string]*domain.CacheEntry
}

func newMockCacheRepoMW() *mockCacheRepoMW {
	return &mockCacheRepoMW{entries: make(map[string]*domain.CacheEntry)}
}

func (m *mockCacheRepoMW) Get(key string) (*domain.CacheEntry, bool) {
	e, ok := m.entries[key]
	return e, ok
}

func (m *mockCacheRepoMW) Set(key string, entry *domain.CacheEntry) {
	m.entries[key] = entry
}

func (m *mockCacheRepoMW) Delete(key string) int {
	if _, ok := m.entries[key]; ok {
		delete(m.entries, key)
		return 1
	}
	return 0
}

func (m *mockCacheRepoMW) DeleteByPrefix(prefix string) int {
	count := 0
	for k := range m.entries {
		if strings.HasPrefix(k, prefix) {
			delete(m.entries, k)
			count++
		}
	}
	return count
}

func (m *mockCacheRepoMW) DeleteByRegex(pattern string) (int, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return 0, err
	}
	count := 0
	for k := range m.entries {
		if re.MatchString(k) {
			delete(m.entries, k)
			count++
		}
	}
	return count, nil
}

func (m *mockCacheRepoMW) DeleteByTag(tag string) int {
	count := 0
	for k, e := range m.entries {
		for _, t := range e.Tags {
			if t == tag {
				delete(m.entries, k)
				count++
				break
			}
		}
	}
	return count
}

func (m *mockCacheRepoMW) DeleteByCondition(_ domain.CacheCondition) int { return 0 }

func (m *mockCacheRepoMW) Flush() int {
	n := len(m.entries)
	m.entries = make(map[string]*domain.CacheEntry)
	return n
}

func (m *mockCacheRepoMW) Stats() domain.CacheStats {
	return domain.CacheStats{Size: len(m.entries)}
}

func minimalCacheCfg() *config.Config {
	return &config.Config{
		Cache: config.CacheConfig{
			Enabled:    true,
			Cache2xx:   true,
			DefaultTTL: time.Minute,
			MaxSize:    100,
		},
	}
}

func helloHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	})
}

func doGet(mw func(http.Handler) http.Handler, target string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, target, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	mw(helloHandler()).ServeHTTP(rec, req)
	return rec
}

func TestCacheMiddleware_FirstRequest_XCacheMiss(t *testing.T) {
	repo := newMockCacheRepoMW()
	cfg := minimalCacheCfg()
	mw := CacheMiddleware(repo, cfg, zap.NewNop())

	rec := doGet(mw, "http://example.com/page", nil)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "MISS", rec.Header().Get("X-Cache"))
	assert.Equal(t, "hello", rec.Body.String())
}

func TestCacheMiddleware_SecondRequest_XCacheHit(t *testing.T) {
	repo := newMockCacheRepoMW()
	cfg := minimalCacheCfg()
	mw := CacheMiddleware(repo, cfg, zap.NewNop())

	doGet(mw, "http://example.com/page", nil)

	rec := doGet(mw, "http://example.com/page", nil)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "HIT", rec.Header().Get("X-Cache"))
	assert.Equal(t, "hello", rec.Body.String())
}

func TestCacheMiddleware_POSTRequest_BypassesCache(t *testing.T) {
	repo := newMockCacheRepoMW()
	cfg := minimalCacheCfg()
	mw := CacheMiddleware(repo, cfg, zap.NewNop())

	req := httptest.NewRequest(http.MethodPost, "http://example.com/page", nil)
	rec := httptest.NewRecorder()
	mw(helloHandler()).ServeHTTP(rec, req)

	assert.NotEqual(t, "HIT", rec.Header().Get("X-Cache"))
	assert.Equal(t, "hello", rec.Body.String())
}

func TestCacheMiddleware_NoCacheHeader_Bypass(t *testing.T) {
	repo := newMockCacheRepoMW()
	cfg := minimalCacheCfg()
	mw := CacheMiddleware(repo, cfg, zap.NewNop())

	rec := doGet(mw, "http://example.com/page", map[string]string{
		"Cache-Control": "no-cache",
	})
	assert.Equal(t, "BYPASS", rec.Header().Get("X-Cache"))
}

func TestCacheMiddleware_CacheDisabled_NoXCacheHeader(t *testing.T) {
	repo := newMockCacheRepoMW()
	cfg := minimalCacheCfg()
	cfg.Cache.Enabled = false
	mw := CacheMiddleware(repo, cfg, zap.NewNop())

	rec := doGet(mw, "http://example.com/page", nil)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Header().Get("X-Cache"))
}

func TestCacheMiddleware_RuleDisabled_Bypass(t *testing.T) {
	repo := newMockCacheRepoMW()
	cfg := minimalCacheCfg()
	boolFalse := false
	cfg.Cache.Rules = []config.CacheRuleConfig{
		{Match: "example.com/disabled", Enabled: &boolFalse},
	}
	mw := CacheMiddleware(repo, cfg, zap.NewNop())

	rec := doGet(mw, "http://example.com/disabled", nil)
	assert.Equal(t, "BYPASS", rec.Header().Get("X-Cache"))
}

func TestCacheMiddleware_AgeHeader_OnHit(t *testing.T) {
	repo := newMockCacheRepoMW()
	cfg := minimalCacheCfg()
	mw := CacheMiddleware(repo, cfg, zap.NewNop())

	doGet(mw, "http://example.com/age-test", nil)

	rec := doGet(mw, "http://example.com/age-test", nil)
	assert.Equal(t, "HIT", rec.Header().Get("X-Cache"))
	assert.NotEmpty(t, rec.Header().Get("Age"))
}

func TestCacheMiddleware_PragmaNoCache_Bypass(t *testing.T) {
	repo := newMockCacheRepoMW()
	cfg := minimalCacheCfg()
	mw := CacheMiddleware(repo, cfg, zap.NewNop())

	rec := doGet(mw, "http://example.com/pragma", map[string]string{
		"Pragma": "no-cache",
	})
	assert.Equal(t, "BYPASS", rec.Header().Get("X-Cache"))
}
