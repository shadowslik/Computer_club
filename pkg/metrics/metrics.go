package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

type Collector struct {
	startedAt     time.Time
	totalRequests atomic.Int64
	totalErrors   atomic.Int64
	totalBytesIn  atomic.Int64
	totalBytesOut atomic.Int64
	activeConns   atomic.Int64

	mu    sync.RWMutex
	rps   float64
	avgMs float64
}

func NewCollector() *Collector {
	return &Collector{startedAt: time.Now()}
}

func (c *Collector) RecordRequest(latencyMs int64, isError bool) {
	c.totalRequests.Add(1)
	if isError {
		c.totalErrors.Add(1)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	total := float64(c.totalRequests.Load())
	c.rps = total / time.Since(c.startedAt).Seconds()
	c.avgMs = (c.avgMs*(total-1) + float64(latencyMs)) / total
}

func (c *Collector) RecordTraffic(bytesIn, bytesOut int64) {
	c.totalBytesIn.Add(bytesIn)
	c.totalBytesOut.Add(bytesOut)
}

func (c *Collector) IncConns() { c.activeConns.Add(1) }
func (c *Collector) DecConns() { c.activeConns.Add(-1) }

// Snapshot is a point-in-time view of server metrics.
type Snapshot struct {
	UptimeSec     int64   `json:"uptime_sec"      example:"3600"`
	TotalRequests int64   `json:"total_requests"  example:"100000"`
	TotalErrors   int64   `json:"total_errors"    example:"42"`
	TotalBytesIn  int64   `json:"total_bytes_in"  example:"10485760"`
	TotalBytesOut int64   `json:"total_bytes_out" example:"52428800"`
	ActiveConns   int64   `json:"active_conns"    example:"17"`
	RPS           float64 `json:"rps"             example:"27.8"`
	AvgLatencyMs  float64 `json:"avg_latency_ms"  example:"12.4"`
}

func (c *Collector) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Snapshot{
		UptimeSec:     int64(time.Since(c.startedAt).Seconds()),
		TotalRequests: c.totalRequests.Load(),
		TotalErrors:   c.totalErrors.Load(),
		TotalBytesIn:  c.totalBytesIn.Load(),
		TotalBytesOut: c.totalBytesOut.Load(),
		ActiveConns:   c.activeConns.Load(),
		RPS:           c.rps,
		AvgLatencyMs:  c.avgMs,
	}
}
