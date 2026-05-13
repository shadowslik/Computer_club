package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockIPChecker struct {
	allowed map[string]bool
	err     error
}

func newMockIPChecker(ips ...string) *mockIPChecker {
	m := &mockIPChecker{allowed: make(map[string]bool)}
	for _, ip := range ips {
		m.allowed[ip] = true
	}
	return m
}

func (m *mockIPChecker) IsAllowed(_ context.Context, ip string) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.allowed[ip], nil
}

type mockDenialsRepoMW struct {
	records []struct{ IP, Reason string }
}

func (m *mockDenialsRepoMW) Record(ip, reason string) {
	m.records = append(m.records, struct{ IP, Reason string }{ip, reason})
}

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}

func buildIPFilterMW(
	whiteUC, blackUC, grayUC ipChecker,
	defaultDeny bool,
	denials denialsRecorder,
) func(http.Handler) http.Handler {
	return IPFilterMiddleware(whiteUC, blackUC, grayUC, defaultDeny, denials, zap.NewNop())
}

func doIPRequest(mw func(http.Handler) http.Handler, remoteAddr string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = remoteAddr + ":1234"
	rec := httptest.NewRecorder()
	mw(okHandler()).ServeHTTP(rec, req)
	return rec
}

func TestIPFilterMiddleware_BlacklistedIP_Returns403(t *testing.T) {
	blackUC := newMockIPChecker("10.0.0.1")
	whiteUC := newMockIPChecker()
	grayUC := newMockIPChecker()
	denials := &mockDenialsRepoMW{}

	mw := buildIPFilterMW(whiteUC, blackUC, grayUC, false, denials)
	rec := doIPRequest(mw, "10.0.0.1")

	assert.Equal(t, http.StatusForbidden, rec.Code)
	require.Len(t, denials.records, 1)
	assert.Equal(t, "blacklist", denials.records[0].Reason)
}

func TestIPFilterMiddleware_GraylistedIP_Returns429WithCaptchaHeader(t *testing.T) {
	blackUC := newMockIPChecker()
	grayUC := newMockIPChecker("10.0.0.2")
	whiteUC := newMockIPChecker()
	denials := &mockDenialsRepoMW{}

	mw := buildIPFilterMW(whiteUC, blackUC, grayUC, false, denials)
	rec := doIPRequest(mw, "10.0.0.2")

	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
	assert.Equal(t, "true", rec.Header().Get("X-CAPTCHA-Required"))
	require.Len(t, denials.records, 1)
	assert.Equal(t, "graylist", denials.records[0].Reason)
}

func TestIPFilterMiddleware_WhitelistedIP_Returns200(t *testing.T) {
	whiteUC := newMockIPChecker("10.0.0.3")
	blackUC := newMockIPChecker()
	grayUC := newMockIPChecker()
	denials := &mockDenialsRepoMW{}

	mw := buildIPFilterMW(whiteUC, blackUC, grayUC, false, denials)
	rec := doIPRequest(mw, "10.0.0.3")

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, denials.records)
}

func TestIPFilterMiddleware_UnknownIP_DefaultDeny_Returns403(t *testing.T) {
	whiteUC := newMockIPChecker()
	blackUC := newMockIPChecker()
	grayUC := newMockIPChecker()
	denials := &mockDenialsRepoMW{}

	mw := buildIPFilterMW(whiteUC, blackUC, grayUC, true, denials)
	rec := doIPRequest(mw, "10.0.0.4")

	assert.Equal(t, http.StatusForbidden, rec.Code)
	require.Len(t, denials.records, 1)
	assert.Equal(t, "not_in_whitelist", denials.records[0].Reason)
}

func TestIPFilterMiddleware_UnknownIP_DefaultAllow_Returns200(t *testing.T) {
	whiteUC := newMockIPChecker()
	blackUC := newMockIPChecker()
	grayUC := newMockIPChecker()
	denials := &mockDenialsRepoMW{}

	mw := buildIPFilterMW(whiteUC, blackUC, grayUC, false, denials)
	rec := doIPRequest(mw, "10.0.0.5")

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, denials.records)
}

func TestIPFilterMiddleware_BlacklistRecordsDenial(t *testing.T) {
	blackUC := newMockIPChecker("192.168.1.100")
	whiteUC := newMockIPChecker()
	grayUC := newMockIPChecker()
	denials := &mockDenialsRepoMW{}

	mw := buildIPFilterMW(whiteUC, blackUC, grayUC, false, denials)
	doIPRequest(mw, "192.168.1.100")

	require.Len(t, denials.records, 1)
	assert.Equal(t, "192.168.1.100", denials.records[0].IP)
	assert.Equal(t, "blacklist", denials.records[0].Reason)
}

func TestIPFilterMiddleware_GraylistRecordsDenial(t *testing.T) {
	grayUC := newMockIPChecker("172.16.0.1")
	blackUC := newMockIPChecker()
	whiteUC := newMockIPChecker()
	denials := &mockDenialsRepoMW{}

	mw := buildIPFilterMW(whiteUC, blackUC, grayUC, false, denials)
	doIPRequest(mw, "172.16.0.1")

	require.Len(t, denials.records, 1)
	assert.Equal(t, "172.16.0.1", denials.records[0].IP)
	assert.Equal(t, "graylist", denials.records[0].Reason)
}
