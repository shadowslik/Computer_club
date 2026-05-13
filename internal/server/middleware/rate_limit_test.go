package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"proxy/internal/domain"
	"proxy/pkg/metrics"
)

type mockRateLimitUCMW struct {
	checkResult   domain.LimitResult
	checkErr      error
	connectResult domain.LimitResult
	connectErr    error
	disconnectErr error
}

func (m *mockRateLimitUCMW) Check(_ context.Context, _ string, _ int64, _ []domain.RateRule) (domain.LimitResult, error) {
	return m.checkResult, m.checkErr
}

func (m *mockRateLimitUCMW) TrackDownload(_ context.Context, _ string, _ int64) error {
	return nil
}

func (m *mockRateLimitUCMW) OnConnect(_ context.Context, _ string) (domain.LimitResult, error) {
	return m.connectResult, m.connectErr
}

func (m *mockRateLimitUCMW) OnDisconnect(_ context.Context, _ string) error {
	return m.disconnectErr
}

func (m *mockRateLimitUCMW) AddRule(_ context.Context, _ domain.RateRule) error { return nil }
func (m *mockRateLimitUCMW) RemoveRule(_ context.Context, _ string) error       { return nil }
func (m *mockRateLimitUCMW) GetRules(_ context.Context) ([]domain.RateRule, error) {
	return nil, nil
}
func (m *mockRateLimitUCMW) GetRulesForIP(_ context.Context, _ string) ([]domain.RateRule, error) {
	return nil, nil
}

type mockViolatorsRepoMW struct {
	banned   map[string]bool
	addedIPs []string
}

func newMockViolatorsRepoMW() *mockViolatorsRepoMW {
	return &mockViolatorsRepoMW{banned: make(map[string]bool)}
}

func (m *mockViolatorsRepoMW) Add(v domain.Violator) {
	m.addedIPs = append(m.addedIPs, v.IP)
}

func (m *mockViolatorsRepoMW) IsBanned(ip string) bool { return m.banned[ip] }

type mockStatsUpdaterMW struct {
	updated []domain.ClientStat
}

func (m *mockStatsUpdaterMW) Update(stat domain.ClientStat) {
	m.updated = append(m.updated, stat)
}

func buildRateLimitMW(
	uc domain.RateLimitUseCase,
	violators violatorsBanner,
	stats *mockStatsUpdaterMW,
) func(http.Handler) http.Handler {
	return RateLimitMiddleware(
		uc,
		zap.NewNop(),
		metrics.NewCollector(),
		stats,
		violators,
		5*time.Minute,
	)
}

func doRLRequest(mw func(http.Handler) http.Handler, remoteAddr string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = remoteAddr + ":9999"
	rec := httptest.NewRecorder()
	mw(okHandler()).ServeHTTP(rec, req)
	return rec
}

func TestRateLimitMiddleware_AllowedRequest_Returns200(t *testing.T) {
	uc := &mockRateLimitUCMW{
		connectResult: domain.LimitResult{Allowed: true},
		checkResult:   domain.LimitResult{Allowed: true},
	}
	violators := newMockViolatorsRepoMW()
	stats := &mockStatsUpdaterMW{}

	mw := buildRateLimitMW(uc, violators, stats)
	rec := doRLRequest(mw, "10.0.0.1")

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Len(t, stats.updated, 1)
}

func TestRateLimitMiddleware_BannedIP_Returns429(t *testing.T) {
	uc := &mockRateLimitUCMW{connectResult: domain.LimitResult{Allowed: true}}
	violators := newMockViolatorsRepoMW()
	violators.banned["10.0.0.2"] = true
	stats := &mockStatsUpdaterMW{}

	mw := buildRateLimitMW(uc, violators, stats)
	rec := doRLRequest(mw, "10.0.0.2")

	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}

func TestRateLimitMiddleware_ConnectionRejected_Returns429(t *testing.T) {
	uc := &mockRateLimitUCMW{
		connectResult: domain.LimitResult{Allowed: false, Reason: "too many conns"},
	}
	violators := newMockViolatorsRepoMW()
	stats := &mockStatsUpdaterMW{}

	mw := buildRateLimitMW(uc, violators, stats)
	rec := doRLRequest(mw, "10.0.0.3")

	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}

func TestRateLimitMiddleware_RateLimitExceeded_Returns429AndAddsViolator(t *testing.T) {
	uc := &mockRateLimitUCMW{
		connectResult: domain.LimitResult{Allowed: true},
		checkResult:   domain.LimitResult{Allowed: false, Reason: "RPS exceeded"},
	}
	violators := newMockViolatorsRepoMW()
	stats := &mockStatsUpdaterMW{}

	mw := buildRateLimitMW(uc, violators, stats)
	rec := doRLRequest(mw, "10.0.0.4")

	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
	assert.Contains(t, violators.addedIPs, "10.0.0.4")
}

func TestErrorRate_4xx_Returns100(t *testing.T) {
	assert.Equal(t, 100.0, errorRate(400))
	assert.Equal(t, 100.0, errorRate(404))
	assert.Equal(t, 100.0, errorRate(500))
}

func TestErrorRate_2xx_Returns0(t *testing.T) {
	assert.Equal(t, 0.0, errorRate(200))
	assert.Equal(t, 0.0, errorRate(301))
}
