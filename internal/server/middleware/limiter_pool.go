package middleware

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type limiterEntry struct {
	limiter  *rate.Limiter
	speedBps int64
	lastUsed time.Time
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

func init() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			limiters.cleanup(10 * time.Minute)
		}
	}()
}

func (l *ipLimiters) getOrUpdate(m map[string]*limiterEntry, ip string, speedBps int64) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()
	entry, ok := m[ip]
	if !ok || entry.speedBps != speedBps {
		lim := rate.NewLimiter(rate.Limit(speedBps), 32*1024)
		m[ip] = &limiterEntry{limiter: lim, speedBps: speedBps, lastUsed: time.Now()}
		return lim
	}
	entry.lastUsed = time.Now()
	return entry.limiter
}

func (l *ipLimiters) cleanup(maxAge time.Duration) {
	cutoff := time.Now().Add(-maxAge)
	l.mu.Lock()
	defer l.mu.Unlock()
	for ip, e := range l.download {
		if e.lastUsed.Before(cutoff) {
			delete(l.download, ip)
		}
	}
	for ip, e := range l.upload {
		if e.lastUsed.Before(cutoff) {
			delete(l.upload, ip)
		}
	}
}
