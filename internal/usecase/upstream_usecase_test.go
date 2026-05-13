package usecase

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proxy/internal/domain"
)

func TestUpstreamService_GetStatuses_ReturnsSnapshot(t *testing.T) {
	targets := []domain.UpstreamTarget{
		{ID: "t1", Name: "Alpha", URL: "http://alpha.local", Type: "http"},
		{ID: "t2", Name: "Beta", URL: "http://beta.local", Type: "http"},
	}
	svc := NewUpstreamService(targets)

	statuses := svc.GetStatuses()
	require.Len(t, statuses, 2)
	assert.Equal(t, "t1", statuses[0].ID)
	assert.Equal(t, "Alpha", statuses[0].Name)
	assert.Equal(t, "t2", statuses[1].ID)
}

func TestUpstreamService_GetStatuses_Empty(t *testing.T) {
	svc := NewUpstreamService(nil)
	statuses := svc.GetStatuses()
	assert.Empty(t, statuses)
}

func TestUpstreamService_InitUpstreams_LatencyHistoryPreallocated(t *testing.T) {
	targets := []domain.UpstreamTarget{
		{ID: "u1", Name: "Upstream1", URL: "http://example.com"},
	}
	svc := NewUpstreamService(targets)

	statuses := svc.GetStatuses()
	require.Len(t, statuses, 1)
	assert.Len(t, statuses[0].LatencyHistory, 20)
}

func TestResetIfNeeded_ResetsSecondCounter(t *testing.T) {
	c := &domain.Counter{
		RequestsSecond: 99,
		LastSecond:     time.Now().Add(-2 * time.Second), // stale — older than 1s
		LastMinute:     time.Now(),
		LastHour:       time.Now(),
		LastDay:        time.Now(),
	}
	resetIfNeeded(c)
	assert.Equal(t, 0, c.RequestsSecond)
}

func TestResetIfNeeded_ResetsMinuteCounter(t *testing.T) {
	c := &domain.Counter{
		RequestsMinute: 50,
		LastSecond:     time.Now(),
		LastMinute:     time.Now().Add(-2 * time.Minute),
		LastHour:       time.Now(),
		LastDay:        time.Now(),
	}
	resetIfNeeded(c)
	assert.Equal(t, 0, c.RequestsMinute)
}

func TestResetIfNeeded_ResetsHourCounter(t *testing.T) {
	c := &domain.Counter{
		RequestsHour: 999,
		LastSecond:   time.Now(),
		LastMinute:   time.Now(),
		LastHour:     time.Now().Add(-2 * time.Hour),
		LastDay:      time.Now(),
	}
	resetIfNeeded(c)
	assert.Equal(t, 0, c.RequestsHour)
}

func TestResetIfNeeded_ResetsDayCounters(t *testing.T) {
	c := &domain.Counter{
		RequestsDay: 5000,
		UploadDay:   1024,
		DownloadDay: 2048,
		LastSecond:  time.Now(),
		LastMinute:  time.Now(),
		LastHour:    time.Now(),
		LastDay:     time.Now().Add(-25 * time.Hour),
	}
	resetIfNeeded(c)
	assert.Equal(t, 0, c.RequestsDay)
	assert.Equal(t, int64(0), c.UploadDay)
	assert.Equal(t, int64(0), c.DownloadDay)
}

func TestUpstreamService_CheckAllUpstreams_HealthyServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	targets := []domain.UpstreamTarget{
		{ID: "u1", Name: "Test", URL: srv.URL},
	}
	svc := NewUpstreamService(targets)

	svc.checkAllUpstreams()

	statuses := svc.GetStatuses()
	require.Len(t, statuses, 1)
	assert.Equal(t, "healthy", statuses[0].Status)
}

func TestUpstreamService_CheckAllUpstreams_DownServer(t *testing.T) {
	targets := []domain.UpstreamTarget{
		{ID: "u2", Name: "Down", URL: "http://127.0.0.1:19999"},
	}
	svc := NewUpstreamService(targets)
	svc.httpClient = &http.Client{Timeout: 100 * time.Millisecond}

	svc.checkAllUpstreams()

	statuses := svc.GetStatuses()
	require.Len(t, statuses, 1)
	assert.Equal(t, "down", statuses[0].Status)
}

func TestResetIfNeeded_NoResetWhenCurrent(t *testing.T) {
	c := &domain.Counter{
		RequestsSecond: 5,
		RequestsMinute: 50,
		RequestsHour:   500,
		RequestsDay:    5000,
		LastSecond:     time.Now(),
		LastMinute:     time.Now(),
		LastHour:       time.Now(),
		LastDay:        time.Now(),
	}
	resetIfNeeded(c)
	assert.Equal(t, 5, c.RequestsSecond)
	assert.Equal(t, 50, c.RequestsMinute)
	assert.Equal(t, 500, c.RequestsHour)
	assert.Equal(t, 5000, c.RequestsDay)
}
