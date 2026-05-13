package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proxy/internal/domain"
)

type mockClientStatsUseCase struct {
	stats []domain.ClientStat
}

func (m *mockClientStatsUseCase) GetTop(limit int, sortBy string) []domain.ClientStat {
	if limit > 0 && limit < len(m.stats) {
		return m.stats[:limit]
	}
	return m.stats
}

func TestClientStatsHandler_GetTopClients_DefaultParams(t *testing.T) {
	uc := &mockClientStatsUseCase{
		stats: []domain.ClientStat{
			{IP: "1.1.1.1", Requests: 100},
			{IP: "2.2.2.2", Requests: 50},
		},
	}
	h := NewClientStatsHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/top-clients", nil)
	rec := httptest.NewRecorder()
	h.GetTopClients(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var list []domain.ClientStat
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&list))
	assert.Len(t, list, 2)
}

func TestClientStatsHandler_GetTopClients_WithLimit(t *testing.T) {
	uc := &mockClientStatsUseCase{
		stats: []domain.ClientStat{
			{IP: "1.1.1.1", Requests: 100},
			{IP: "2.2.2.2", Requests: 50},
			{IP: "3.3.3.3", Requests: 10},
		},
	}
	h := NewClientStatsHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/top-clients?limit=2", nil)
	rec := httptest.NewRecorder()
	h.GetTopClients(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var list []domain.ClientStat
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&list))
	assert.Len(t, list, 2)
}

func TestClientStatsHandler_GetTopClients_SortByTraffic(t *testing.T) {
	uc := &mockClientStatsUseCase{stats: []domain.ClientStat{{IP: "4.4.4.4"}}}
	h := NewClientStatsHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/top-clients?sort=traffic", nil)
	rec := httptest.NewRecorder()
	h.GetTopClients(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestClientStatsHandler_GetTopClients_InvalidLimit_UsesDefault(t *testing.T) {
	uc := &mockClientStatsUseCase{}
	h := NewClientStatsHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/top-clients?limit=notanumber", nil)
	rec := httptest.NewRecorder()
	h.GetTopClients(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
