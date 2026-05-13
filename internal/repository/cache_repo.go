package repository

import (
	"net/http"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"proxy/internal/domain"
)

type cacheItem struct {
	entry *domain.CacheEntry
}

type CacheRepoImpl struct {
	mu           sync.RWMutex
	entries      map[string]*cacheItem
	tagIdx       map[string]map[string]struct{}
	maxSize      int
	cascadeRules []domain.CacheCascadeRule

	hits      atomic.Int64
	misses    atomic.Int64
	evictions atomic.Int64
	bytesUsed atomic.Int64
}

func NewCacheRepo(maxSize int, cleanupInterval time.Duration, cascadeRules []domain.CacheCascadeRule) *CacheRepoImpl {
	r := &CacheRepoImpl{
		entries:      make(map[string]*cacheItem),
		tagIdx:       make(map[string]map[string]struct{}),
		maxSize:      maxSize,
		cascadeRules: cascadeRules,
	}
	if cleanupInterval > 0 {
		go r.runCleanup(cleanupInterval)
	}
	return r
}

func (r *CacheRepoImpl) StartChangeDetector(interval time.Duration, client *http.Client) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			r.revalidateAll(client)
		}
	}()
}

func (r *CacheRepoImpl) revalidateAll(client *http.Client) {
	type snap struct {
		key  string
		etag string
	}
	r.mu.RLock()
	snaps := make([]snap, 0)
	for key, item := range r.entries {
		if etag := item.entry.Headers.Get("Etag"); etag != "" {
			snaps = append(snaps, snap{key, etag})
		}
	}
	r.mu.RUnlock()

	for _, s := range snaps {
		r.revalidateEntry(s.key, s.etag, client)
	}
}

func (r *CacheRepoImpl) revalidateEntry(key, etag string, client *http.Client) {
	_, rest, ok := strings.Cut(key, ":")
	if !ok {
		return
	}
	rawURL := "http://" + rest
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return
	}
	req.Header.Set("If-None-Match", etag)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		r.Delete(key)
	}
}

func (r *CacheRepoImpl) runCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		r.mu.Lock()
		r.evictExpiredLocked()
		r.mu.Unlock()
	}
}

func (r *CacheRepoImpl) deleteItemLocked(key string, item *cacheItem) {
	for _, tag := range item.entry.Tags {
		if s, ok := r.tagIdx[tag]; ok {
			delete(s, key)
			if len(s) == 0 {
				delete(r.tagIdx, tag)
			}
		}
	}
	r.bytesUsed.Add(-int64(len(item.entry.Body)))
	delete(r.entries, key)
}

func (r *CacheRepoImpl) evictExpiredLocked() {
	for key, item := range r.entries {
		if item.entry.IsExpired() {
			r.deleteItemLocked(key, item)
			r.evictions.Add(1)
		}
	}
}

func (r *CacheRepoImpl) Get(key string) (*domain.CacheEntry, bool) {
	r.mu.RLock()
	item, ok := r.entries[key]
	r.mu.RUnlock()

	if !ok {
		r.misses.Add(1)
		return nil, false
	}
	if item.entry.IsExpired() {
		r.mu.Lock()
		if it, still := r.entries[key]; still && it.entry.IsExpired() {
			r.deleteItemLocked(key, it)
			r.evictions.Add(1)
		}
		r.mu.Unlock()
		r.misses.Add(1)
		return nil, false
	}
	r.hits.Add(1)
	return item.entry, true
}

func (r *CacheRepoImpl) Set(key string, entry *domain.CacheEntry) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.maxSize > 0 && len(r.entries) >= r.maxSize {
		r.evictExpiredLocked()
		if len(r.entries) >= r.maxSize {
			return
		}
	}

	if existing, ok := r.entries[key]; ok {
		r.deleteItemLocked(key, existing)
	}

	r.entries[key] = &cacheItem{entry: entry}
	r.bytesUsed.Add(int64(len(entry.Body)))

	for _, tag := range entry.Tags {
		if r.tagIdx[tag] == nil {
			r.tagIdx[tag] = make(map[string]struct{})
		}
		r.tagIdx[tag][key] = struct{}{}
	}
}

func (r *CacheRepoImpl) Delete(key string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.entries[key]
	if !ok {
		return 0
	}
	r.deleteItemLocked(key, item)
	return 1
}

func (r *CacheRepoImpl) DeleteByPrefix(prefix string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	count := 0
	for key, item := range r.entries {
		if strings.HasPrefix(key, prefix) {
			r.deleteItemLocked(key, item)
			count++
		}
	}
	return count
}

func (r *CacheRepoImpl) DeleteByRegex(pattern string) (int, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return 0, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	count := 0
	for key, item := range r.entries {
		if re.MatchString(key) {
			r.deleteItemLocked(key, item)
			count++
		}
	}
	return count, nil
}

func (r *CacheRepoImpl) DeleteByTag(tag string) int {
	count := r.deleteByTagOnce(tag)
	for _, rule := range r.cascadeRules {
		if rule.TriggerTag == tag {
			for _, cascadeTag := range rule.CascadeTags {
				count += r.deleteByTagOnce(cascadeTag)
			}
		}
	}
	return count
}

func (r *CacheRepoImpl) deleteByTagOnce(tag string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	keys, ok := r.tagIdx[tag]
	if !ok {
		return 0
	}
	toDelete := make([]string, 0, len(keys))
	for key := range keys {
		toDelete = append(toDelete, key)
	}
	count := 0
	for _, key := range toDelete {
		if item, ok := r.entries[key]; ok {
			r.deleteItemLocked(key, item)
			count++
		}
	}
	return count
}

func (r *CacheRepoImpl) DeleteByCondition(cond domain.CacheCondition) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	count := 0
	for key, item := range r.entries {
		if r.matchesCondition(item.entry, cond) {
			r.deleteItemLocked(key, item)
			count++
		}
	}
	return count
}

func (r *CacheRepoImpl) matchesCondition(e *domain.CacheEntry, cond domain.CacheCondition) bool {
	if cond.StaleOnly && !e.IsExpired() {
		return false
	}
	if len(cond.Tags) > 0 {
		hasTag := false
		for _, reqTag := range cond.Tags {
			for _, entryTag := range e.Tags {
				if entryTag == reqTag {
					hasTag = true
					break
				}
			}
			if hasTag {
				break
			}
		}
		if !hasTag {
			return false
		}
	}
	return true
}

func (r *CacheRepoImpl) Flush() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	count := len(r.entries)
	r.entries = make(map[string]*cacheItem)
	r.tagIdx = make(map[string]map[string]struct{})
	r.bytesUsed.Store(0)
	return count
}

func (r *CacheRepoImpl) Stats() domain.CacheStats {
	r.mu.RLock()
	size := len(r.entries)
	r.mu.RUnlock()
	return domain.CacheStats{
		Size:      size,
		Hits:      r.hits.Load(),
		Misses:    r.misses.Load(),
		Evictions: r.evictions.Load(),
		BytesUsed: r.bytesUsed.Load(),
	}
}
