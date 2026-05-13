package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proxy/internal/domain"
)

type mockCacheRepo struct {
	entries map[string]*domain.CacheEntry
	stats   domain.CacheStats

	deleteByRegexErr error
}

func newMockCacheRepo() *mockCacheRepo {
	return &mockCacheRepo{entries: make(map[string]*domain.CacheEntry)}
}

func (m *mockCacheRepo) Get(key string) (*domain.CacheEntry, bool) {
	e, ok := m.entries[key]
	return e, ok
}

func (m *mockCacheRepo) Set(key string, entry *domain.CacheEntry) {
	m.entries[key] = entry
}

func (m *mockCacheRepo) Delete(key string) int {
	if _, ok := m.entries[key]; ok {
		delete(m.entries, key)
		return 1
	}
	return 0
}

func (m *mockCacheRepo) DeleteByPrefix(prefix string) int {
	count := 0
	for k := range m.entries {
		if strings.HasPrefix(k, prefix) {
			delete(m.entries, k)
			count++
		}
	}
	return count
}

func (m *mockCacheRepo) DeleteByRegex(pattern string) (int, error) {
	if m.deleteByRegexErr != nil {
		return 0, m.deleteByRegexErr
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return 0, err
	}
	count := 0
	for k := range m.entries {
		if re.MatchString(k) {
			delete(m.entries, k)
			count++
		}
	}
	return count, nil
}

func (m *mockCacheRepo) DeleteByTag(tag string) int {
	count := 0
	for k, e := range m.entries {
		for _, t := range e.Tags {
			if t == tag {
				delete(m.entries, k)
				count++
				break
			}
		}
	}
	return count
}

func (m *mockCacheRepo) DeleteByCondition(cond domain.CacheCondition) int {
	return len(m.entries)
}

func (m *mockCacheRepo) Flush() int {
	n := len(m.entries)
	m.entries = make(map[string]*domain.CacheEntry)
	return n
}

func (m *mockCacheRepo) Stats() domain.CacheStats {
	return m.stats
}

func newCacheHandlerAndRepo() (*CacheHandler, *mockCacheRepo) {
	repo := newMockCacheRepo()
	return NewCacheHandler(repo), repo
}

func TestCacheHandler_GetStats_Returns200WithJSON(t *testing.T) {
	h, repo := newCacheHandlerAndRepo()
	repo.stats = domain.CacheStats{Size: 5, Hits: 10, Misses: 2}

	req := httptest.NewRequest(http.MethodGet, "/cache/stats", nil)
	rec := httptest.NewRecorder()
	h.GetStats(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var stats domain.CacheStats
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&stats))
	assert.Equal(t, 5, stats.Size)
	assert.Equal(t, int64(10), stats.Hits)
	assert.Equal(t, int64(2), stats.Misses)
}

func TestCacheHandler_Flush_Returns200WithDeletedCount(t *testing.T) {
	h, repo := newCacheHandlerAndRepo()
	repo.entries["k1"] = &domain.CacheEntry{}
	repo.entries["k2"] = &domain.CacheEntry{}

	req := httptest.NewRequest(http.MethodDelete, "/cache", nil)
	rec := httptest.NewRecorder()
	h.Flush(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp invalidateResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, 2, resp.Deleted)
}

func TestCacheHandler_DeleteByKey_ExistingKey(t *testing.T) {
	h, repo := newCacheHandlerAndRepo()
	repo.entries["foo"] = &domain.CacheEntry{}

	req := httptest.NewRequest(http.MethodDelete, "/cache/key?key=foo", nil)
	rec := httptest.NewRecorder()
	h.DeleteByKey(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp invalidateResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, 1, resp.Deleted)
}

func TestCacheHandler_DeleteByKey_MissingParam_Returns400(t *testing.T) {
	h, _ := newCacheHandlerAndRepo()

	req := httptest.NewRequest(http.MethodDelete, "/cache/key", nil)
	rec := httptest.NewRecorder()
	h.DeleteByKey(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCacheHandler_DeleteByPrefix_Returns200(t *testing.T) {
	h, repo := newCacheHandlerAndRepo()
	repo.entries["GET:/api/v1"] = &domain.CacheEntry{}
	repo.entries["GET:/api/v2"] = &domain.CacheEntry{}

	req := httptest.NewRequest(http.MethodDelete, "/cache/prefix?prefix=GET:", nil)
	rec := httptest.NewRecorder()
	h.DeleteByPrefix(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp invalidateResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, 2, resp.Deleted)
}

func TestCacheHandler_DeleteByPrefix_MissingParam_Returns400(t *testing.T) {
	h, _ := newCacheHandlerAndRepo()

	req := httptest.NewRequest(http.MethodDelete, "/cache/prefix", nil)
	rec := httptest.NewRecorder()
	h.DeleteByPrefix(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCacheHandler_DeleteByRegex_ValidPattern_Returns200(t *testing.T) {
	h, _ := newCacheHandlerAndRepo()

	req := httptest.NewRequest(http.MethodDelete, "/cache/regex?pattern=[0-9]+", nil)
	rec := httptest.NewRecorder()
	h.DeleteByRegex(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCacheHandler_DeleteByRegex_InvalidPattern_Returns400(t *testing.T) {
	h, _ := newCacheHandlerAndRepo()

	req := httptest.NewRequest(http.MethodDelete, "/cache/regex?pattern=%5Binvalid", nil)
	rec := httptest.NewRecorder()
	h.DeleteByRegex(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCacheHandler_DeleteByRegex_MissingParam_Returns400(t *testing.T) {
	h, _ := newCacheHandlerAndRepo()

	req := httptest.NewRequest(http.MethodDelete, "/cache/regex", nil)
	rec := httptest.NewRecorder()
	h.DeleteByRegex(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCacheHandler_DeleteByTag_Returns200(t *testing.T) {
	h, repo := newCacheHandlerAndRepo()
	repo.entries["k"] = &domain.CacheEntry{Tags: []string{"public"}}

	req := httptest.NewRequest(http.MethodDelete, "/cache/tag?tag=public", nil)
	rec := httptest.NewRecorder()
	h.DeleteByTag(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp invalidateResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, 1, resp.Deleted)
}

func TestCacheHandler_DeleteByTag_MissingParam_Returns400(t *testing.T) {
	h, _ := newCacheHandlerAndRepo()

	req := httptest.NewRequest(http.MethodDelete, "/cache/tag", nil)
	rec := httptest.NewRecorder()
	h.DeleteByTag(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCacheHandler_DeleteByCondition_ValidBody_Returns200(t *testing.T) {
	h, repo := newCacheHandlerAndRepo()
	repo.entries["stale"] = &domain.CacheEntry{}

	body := `{"stale_only":true}`
	req := httptest.NewRequest(http.MethodDelete, "/cache/condition", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.DeleteByCondition(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCacheHandler_DeleteByCondition_InvalidBody_Returns400(t *testing.T) {
	h, _ := newCacheHandlerAndRepo()

	req := httptest.NewRequest(http.MethodDelete, "/cache/condition", strings.NewReader("{invalid json"))
	rec := httptest.NewRecorder()
	h.DeleteByCondition(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCacheHandler_PostEvent_ValidBody_Returns200(t *testing.T) {
	h, repo := newCacheHandlerAndRepo()
	repo.entries["user-key"] = &domain.CacheEntry{Tags: []string{"users"}}

	body := `{"event":"user_updated","tags":["users"],"prefix":"","keys":[]}`
	req := httptest.NewRequest(http.MethodPost, "/cache/events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.PostEvent(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result domain.CacheEventResult
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&result))
	assert.Equal(t, 1, result.Deleted)
}

func TestCacheHandler_PostEvent_InvalidBody_Returns400(t *testing.T) {
	h, _ := newCacheHandlerAndRepo()

	req := httptest.NewRequest(http.MethodPost, "/cache/events", strings.NewReader("not json"))
	rec := httptest.NewRecorder()
	h.PostEvent(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
