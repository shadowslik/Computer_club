package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proxy/pkg/metrics"
)

func TestMetricsHandler_HandleMetrics_Returns200(t *testing.T) {
	col := metrics.NewCollector()
	col.RecordRequest(10, false)
	h := NewMetricsHandler(col)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	h.HandleMetrics(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var snap metrics.Snapshot
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&snap))
	assert.Equal(t, int64(1), snap.TotalRequests)
}

func TestMetricsHandler_HandleHealth_Returns200(t *testing.T) {
	col := metrics.NewCollector()
	h := NewMetricsHandler(col)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.HandleHealth(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, "ok", resp["status"])
}
