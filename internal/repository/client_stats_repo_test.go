package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proxy/internal/domain"
)

func TestClientStatsRepo_UpdateNewEntry(t *testing.T) {
	r := NewClientStatsRepo()
	r.Update(domain.ClientStat{IP: "1.1.1.1", Requests: 5, TrafficMB: 1.5})

	top := r.GetTop(10, "requests")
	require.Len(t, top, 1)
	assert.Equal(t, "1.1.1.1", top[0].IP)
	assert.Equal(t, int64(5), top[0].Requests)
}

func TestClientStatsRepo_UpdateMergesExistingEntry(t *testing.T) {
	r := NewClientStatsRepo()
	r.Update(domain.ClientStat{IP: "2.2.2.2", Requests: 10, TrafficMB: 2.0})
	r.Update(domain.ClientStat{IP: "2.2.2.2", Requests: 5, TrafficMB: 1.0})

	top := r.GetTop(10, "requests")
	require.Len(t, top, 1)
	assert.Equal(t, int64(15), top[0].Requests)
	assert.InDelta(t, 3.0, top[0].TrafficMB, 0.001)
}

func TestClientStatsRepo_GetTop_SortByRequests(t *testing.T) {
	r := NewClientStatsRepo()
	r.Update(domain.ClientStat{IP: "a", Requests: 5})
	r.Update(domain.ClientStat{IP: "b", Requests: 50})
	r.Update(domain.ClientStat{IP: "c", Requests: 10})

	top := r.GetTop(10, "requests")
	require.Len(t, top, 3)
	assert.Equal(t, int64(50), top[0].Requests)
}

func TestClientStatsRepo_GetTop_SortByTraffic(t *testing.T) {
	r := NewClientStatsRepo()
	r.Update(domain.ClientStat{IP: "x", TrafficMB: 1.0})
	r.Update(domain.ClientStat{IP: "y", TrafficMB: 5.0})

	top := r.GetTop(10, "traffic")
	require.Len(t, top, 2)
	assert.Equal(t, float64(5.0), top[0].TrafficMB)
}

func TestClientStatsRepo_GetTop_SortByErrors(t *testing.T) {
	r := NewClientStatsRepo()
	r.Update(domain.ClientStat{IP: "p", ErrorRate: 10.0})
	r.Update(domain.ClientStat{IP: "q", ErrorRate: 50.0})

	top := r.GetTop(10, "errors")
	require.Len(t, top, 2)
	assert.Equal(t, float64(50.0), top[0].ErrorRate)
}

func TestClientStatsRepo_GetTop_LimitRespected(t *testing.T) {
	r := NewClientStatsRepo()
	for i := 0; i < 5; i++ {
		r.Update(domain.ClientStat{IP: string(rune('a' + i)), Requests: int64(i)})
	}

	top := r.GetTop(3, "requests")
	assert.Len(t, top, 3)
}

func TestClientStatsRepo_GetTop_Empty(t *testing.T) {
	r := NewClientStatsRepo()
	assert.Empty(t, r.GetTop(10, "requests"))
}
