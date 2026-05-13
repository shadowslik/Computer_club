package domain

import (
	"net/http"
	"time"
)

type CacheEntry struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Tags       []string
	CreatedAt  time.Time
	TTL        time.Duration
}

func (e *CacheEntry) IsExpired() bool {
	return e.TTL > 0 && time.Since(e.CreatedAt) > e.TTL
}

type CacheStats struct {
	Size      int   `json:"size"`
	Hits      int64 `json:"hits"`
	Misses    int64 `json:"misses"`
	Evictions int64 `json:"evictions"`
	BytesUsed int64 `json:"bytesUsed"`
}

type CacheCondition struct {
	StaleOnly bool     `json:"stale_only"`
	Tags      []string `json:"tags"`
}

type CacheCascadeRule struct {
	TriggerTag  string
	CascadeTags []string
}

type CacheEventRequest struct {
	Event  string   `json:"event"`
	Tags   []string `json:"tags"`
	Prefix string   `json:"prefix"`
	Keys   []string `json:"keys"`
}

type CacheEventResult struct {
	Deleted int `json:"deleted"`
}

type CacheRepo interface {
	Get(key string) (*CacheEntry, bool)
	Set(key string, entry *CacheEntry)
	Delete(key string) int
	DeleteByPrefix(prefix string) int
	DeleteByRegex(pattern string) (int, error)
	DeleteByTag(tag string) int
	DeleteByCondition(cond CacheCondition) int
	Flush() int
	Stats() CacheStats
}
