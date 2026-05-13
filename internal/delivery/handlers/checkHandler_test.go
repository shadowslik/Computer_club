package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockIPCheckUseCase struct {
	allowed map[string]bool
}

func newMockIPCheckUC(ips ...string) *mockIPCheckUseCase {
	m := &mockIPCheckUseCase{allowed: make(map[string]bool)}
	for _, ip := range ips {
		m.allowed[ip] = true
	}
	return m
}

func (m *mockIPCheckUseCase) IsAllowed(_ context.Context, ip string) (bool, error) {
	return m.allowed[ip], nil
}

type mockDenialsRecorder struct {
	records []struct{ ip, reason string }
}

func (m *mockDenialsRecorder) Record(ip, reason string) {
	m.records = append(m.records, struct{ ip, reason string }{ip, reason})
}

func newCheckHandler(whiteIPs, blackIPs, grayIPs []string, denials *mockDenialsRecorder) *CheckIpHandler {
	white := newMockIPCheckUC(whiteIPs...)
	black := newMockIPCheckUC(blackIPs...)
	gray := newMockIPCheckUC(grayIPs...)
	return NewCheckIpHandler(white, black, gray, zap.NewNop(), denials)
}

func TestCheckIpHandler_QueryParam_WhiteIP(t *testing.T) {
	denials := &mockDenialsRecorder{}
	h := newCheckHandler([]string{"1.2.3.4"}, nil, nil, denials)

	req := httptest.NewRequest(http.MethodGet, "/ip/check?ip=1.2.3.4", nil)
	rec := httptest.NewRecorder()
	h.CheckIp(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp CheckIpResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, "white", resp.List)
	assert.Equal(t, "1.2.3.4", resp.IP)
}

func TestCheckIpHandler_QueryParam_BlackIP(t *testing.T) {
	denials := &mockDenialsRecorder{}
	h := newCheckHandler(nil, []string{"5.5.5.5"}, nil, denials)

	req := httptest.NewRequest(http.MethodGet, "/ip/check?ip=5.5.5.5", nil)
	rec := httptest.NewRecorder()
	h.CheckIp(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp CheckIpResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, "black", resp.List)
	require.Len(t, denials.records, 1)
	assert.Equal(t, "blacklist", denials.records[0].reason)
}

func TestCheckIpHandler_QueryParam_GrayIP(t *testing.T) {
	denials := &mockDenialsRecorder{}
	h := newCheckHandler(nil, nil, []string{"6.6.6.6"}, denials)

	req := httptest.NewRequest(http.MethodGet, "/ip/check?ip=6.6.6.6", nil)
	rec := httptest.NewRecorder()
	h.CheckIp(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp CheckIpResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, "gray", resp.List)
	require.Len(t, denials.records, 1)
	assert.Equal(t, "graylist", denials.records[0].reason)
}

func TestCheckIpHandler_QueryParam_NoneIP(t *testing.T) {
	denials := &mockDenialsRecorder{}
	h := newCheckHandler(nil, nil, nil, denials)

	req := httptest.NewRequest(http.MethodGet, "/ip/check?ip=7.7.7.7", nil)
	rec := httptest.NewRecorder()
	h.CheckIp(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp CheckIpResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, "none", resp.List)
}

func TestCheckIpHandler_InvalidIP_Returns400(t *testing.T) {
	denials := &mockDenialsRecorder{}
	h := newCheckHandler(nil, nil, nil, denials)

	req := httptest.NewRequest(http.MethodGet, "/ip/check?ip=not-an-ip", nil)
	rec := httptest.NewRecorder()
	h.CheckIp(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCheckIpHandler_IPFromBody(t *testing.T) {
	denials := &mockDenialsRecorder{}
	h := newCheckHandler([]string{"8.8.8.8"}, nil, nil, denials)

	req := httptest.NewRequest(http.MethodGet, "/ip/check",
		strings.NewReader(`"8.8.8.8"`))
	rec := httptest.NewRecorder()
	h.CheckIp(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp CheckIpResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, "white", resp.List)
}

func TestCheckIpHandler_IPFromRemoteAddr(t *testing.T) {
	denials := &mockDenialsRecorder{}
	h := newCheckHandler(nil, nil, nil, denials)

	req := httptest.NewRequest(http.MethodGet, "/ip/check", nil)
	req.RemoteAddr = "9.9.9.9:5678"
	rec := httptest.NewRecorder()
	h.CheckIp(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp CheckIpResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, "9.9.9.9", resp.IP)
}
