package middleware

import (
	"golang.org/x/time/rate"
	"sync"
)

type limiterEntry struct {
	limiter  *rate.Limiter
	speedBps int64
}

type ipLimiters struct {
	mu       sync.Mutex
	download map[string]*limiterEntry
	upload   map[string]*limiterEntry
}

var limiters = &ipLimiters{
	download: make(map[string]*limiterEntry),
	upload:   make(map[string]*limiterEntry),
}

func (l *ipLimiters) getOrUpdate(m map[string]*limiterEntry, ip string, speedBps int64) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()
	entry, ok := m[ip]
	if !ok || entry.speedBps != speedBps {
		lim := rate.NewLimiter(rate.Limit(speedBps), 32*1024)
		m[ip] = &limiterEntry{limiter: lim, speedBps: speedBps}
		return lim
	}
	return entry.limiter
}
