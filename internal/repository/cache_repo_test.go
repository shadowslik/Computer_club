package repository

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proxy/internal/domain"
)

func newEntry(body []byte, tags []string, ttl time.Duration) *domain.CacheEntry {
	return &domain.CacheEntry{
		StatusCode: http.StatusOK,
		Headers:    http.Header{"Content-Type": []string{"text/plain"}},
		Body:       body,
		Tags:       tags,
		CreatedAt:  time.Now(),
		TTL:        ttl,
	}
}

func TestCacheRepo_GetMiss(t *testing.T) {
	r := NewCacheRepo(10, 0, nil)
	entry, ok := r.Get("nonexistent")
	assert.False(t, ok)
	assert.Nil(t, entry)
	assert.Equal(t, int64(1), r.Stats().Misses)
}

func TestCacheRepo_SetAndGet_Hit(t *testing.T) {
	r := NewCacheRepo(10, 0, nil)
	e := newEntry([]byte("hello"), nil, time.Minute)
	r.Set("key1", e)

	got, ok := r.Get("key1")
	require.True(t, ok)
	assert.Equal(t, []byte("hello"), got.Body)
	assert.Equal(t, int64(1), r.Stats().Hits)
	assert.Equal(t, int64(0), r.Stats().Misses)
}

func TestCacheRepo_GetExpiredEntry_Miss(t *testing.T) {
	r := NewCacheRepo(10, 0, nil)
	e := &domain.CacheEntry{
		StatusCode: http.StatusOK,
		Headers:    http.Header{},
		Body:       []byte("expired"),
		Tags:       nil,
		CreatedAt:  time.Now().Add(-10 * time.Second),
		TTL:        5 * time.Second, // already expired
	}
	r.Set("expkey", e)

	got, ok := r.Get("expkey")
	assert.False(t, ok)
	assert.Nil(t, got)
	assert.Equal(t, int64(1), r.Stats().Misses)
	assert.Equal(t, int64(1), r.Stats().Evictions)
}

func TestCacheRepo_SetOverwritesKey(t *testing.T) {
	r := NewCacheRepo(10, 0, nil)
	r.Set("k", newEntry([]byte("first"), nil, time.Minute))
	r.Set("k", newEntry([]byte("second"), nil, time.Minute))

	got, ok := r.Get("k")
	require.True(t, ok)
	assert.Equal(t, []byte("second"), got.Body)
	assert.Equal(t, 1, r.Stats().Size)
}

func TestCacheRepo_DeleteExistingKey(t *testing.T) {
	r := NewCacheRepo(10, 0, nil)
	r.Set("del", newEntry([]byte("x"), nil, time.Minute))

	n := r.Delete("del")
	assert.Equal(t, 1, n)
	assert.Equal(t, 0, r.Stats().Size)
}

func TestCacheRepo_DeleteNonExistingKey(t *testing.T) {
	r := NewCacheRepo(10, 0, nil)
	n := r.Delete("ghost")
	assert.Equal(t, 0, n)
}

func TestCacheRepo_DeleteByPrefix(t *testing.T) {
	r := NewCacheRepo(100, 0, nil)
	r.Set("GET:/api/users", newEntry([]byte("u"), nil, time.Minute))
	r.Set("GET:/api/posts", newEntry([]byte("p"), nil, time.Minute))
	r.Set("POST:/api/users", newEntry([]byte("pu"), nil, time.Minute))

	n := r.DeleteByPrefix("GET:")
	assert.Equal(t, 2, n)
	assert.Equal(t, 1, r.Stats().Size)
}

func TestCacheRepo_DeleteByRegex_Valid(t *testing.T) {
	r := NewCacheRepo(100, 0, nil)
	r.Set("GET:/api/v1/users", newEntry([]byte("a"), nil, time.Minute))
	r.Set("GET:/api/v2/users", newEntry([]byte("b"), nil, time.Minute))
	r.Set("POST:/other", newEntry([]byte("c"), nil, time.Minute))

	n, err := r.DeleteByRegex(`GET:/api/v[0-9]+/users`)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, 1, r.Stats().Size)
}

func TestCacheRepo_DeleteByRegex_InvalidPattern(t *testing.T) {
	r := NewCacheRepo(10, 0, nil)
	r.Set("key", newEntry([]byte("v"), nil, time.Minute))

	_, err := r.DeleteByRegex(`[invalid`)
	assert.Error(t, err)
	assert.Equal(t, 1, r.Stats().Size) // nothing deleted
}

func TestCacheRepo_DeleteByTag(t *testing.T) {
	r := NewCacheRepo(100, 0, nil)
	r.Set("k1", newEntry([]byte("a"), []string{"public", "user"}, time.Minute))
	r.Set("k2", newEntry([]byte("b"), []string{"public"}, time.Minute))
	r.Set("k3", newEntry([]byte("c"), []string{"private"}, time.Minute))

	n := r.DeleteByTag("public")
	assert.Equal(t, 2, n)
	assert.Equal(t, 1, r.Stats().Size)
}

func TestCacheRepo_DeleteByTag_CascadeRule(t *testing.T) {
	rules := []domain.CacheCascadeRule{
		{TriggerTag: "a", CascadeTags: []string{"b"}},
	}
	r := NewCacheRepo(100, 0, rules)
	r.Set("ka", newEntry([]byte("a"), []string{"a"}, time.Minute))
	r.Set("kb", newEntry([]byte("b"), []string{"b"}, time.Minute))
	r.Set("kc", newEntry([]byte("c"), []string{"c"}, time.Minute))

	n := r.DeleteByTag("a")
	assert.Equal(t, 2, n)
	assert.Equal(t, 1, r.Stats().Size)
}

func TestCacheRepo_DeleteByCondition_StaleOnly(t *testing.T) {
	r := NewCacheRepo(100, 0, nil)
	fresh := newEntry([]byte("fresh"), nil, time.Minute)
	stale := &domain.CacheEntry{
		StatusCode: http.StatusOK,
		Headers:    http.Header{},
		Body:       []byte("stale"),
		CreatedAt:  time.Now().Add(-10 * time.Second),
		TTL:        1 * time.Second,
	}
	r.Set("fresh", fresh)
	r.Set("stale", stale)

	n := r.DeleteByCondition(domain.CacheCondition{StaleOnly: true})
	assert.Equal(t, 1, n)
	assert.Equal(t, 1, r.Stats().Size)

	_, ok := r.Get("fresh")
	assert.True(t, ok)
}

func TestCacheRepo_DeleteByCondition_TagsFilter(t *testing.T) {
	r := NewCacheRepo(100, 0, nil)
	r.Set("t1", newEntry([]byte("1"), []string{"alpha"}, time.Minute))
	r.Set("t2", newEntry([]byte("2"), []string{"beta"}, time.Minute))
	r.Set("t3", newEntry([]byte("3"), []string{"gamma"}, time.Minute))

	n := r.DeleteByCondition(domain.CacheCondition{Tags: []string{"alpha", "gamma"}})
	assert.Equal(t, 2, n)
	assert.Equal(t, 1, r.Stats().Size)
}

func TestCacheRepo_Flush(t *testing.T) {
	r := NewCacheRepo(100, 0, nil)
	r.Set("a", newEntry([]byte("1"), nil, time.Minute))
	r.Set("b", newEntry([]byte("2"), nil, time.Minute))
	r.Set("c", newEntry([]byte("3"), nil, time.Minute))

	n := r.Flush()
	assert.Equal(t, 3, n)
	assert.Equal(t, 0, r.Stats().Size)
	assert.Equal(t, int64(0), r.Stats().BytesUsed)
}

func TestCacheRepo_Stats_Correctness(t *testing.T) {
	r := NewCacheRepo(100, 0, nil)
	body := []byte("hello")
	r.Set("k", newEntry(body, nil, time.Minute))

	r.Get("k")
	r.Get("missing")

	s := r.Stats()
	assert.Equal(t, 1, s.Size)
	assert.Equal(t, int64(1), s.Hits)
	assert.Equal(t, int64(1), s.Misses)
	assert.Equal(t, int64(len(body)), s.BytesUsed)
}

func TestCacheRepo_MaxSize_FullAndFresh_DropsNewSet(t *testing.T) {
	r := NewCacheRepo(2, 0, nil)
	r.Set("k1", newEntry([]byte("1"), nil, time.Minute))
	r.Set("k2", newEntry([]byte("2"), nil, time.Minute))

	r.Set("k3", newEntry([]byte("3"), nil, time.Minute))

	assert.Equal(t, 2, r.Stats().Size)
	_, ok := r.Get("k3")
	assert.False(t, ok)
}

func TestCacheRepo_TagIndexCleanedAfterDelete(t *testing.T) {
	r := NewCacheRepo(10, 0, nil)
	r.Set("k", newEntry([]byte("v"), []string{"mytag"}, time.Minute))
	r.Delete("k")

	n := r.DeleteByTag("mytag")
	assert.Equal(t, 0, n)
	assert.Equal(t, 0, r.Stats().Size)
}
