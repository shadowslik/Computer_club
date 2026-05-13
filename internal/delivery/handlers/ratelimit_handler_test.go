package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proxy/internal/domain"
)

type mockRateLimitUseCase struct {
	rules         []domain.RateRule
	getRulesErr   error
	addRuleErr    error
	removeRuleErr error
}

func (m *mockRateLimitUseCase) Check(_ context.Context, _ string, _ int64, _ []domain.RateRule) (domain.LimitResult, error) {
	return domain.LimitResult{Allowed: true}, nil
}

func (m *mockRateLimitUseCase) TrackDownload(_ context.Context, _ string, _ int64) error {
	return nil
}

func (m *mockRateLimitUseCase) OnConnect(_ context.Context, _ string) (domain.LimitResult, error) {
	return domain.LimitResult{Allowed: true}, nil
}

func (m *mockRateLimitUseCase) OnDisconnect(_ context.Context, _ string) error {
	return nil
}

func (m *mockRateLimitUseCase) AddRule(_ context.Context, rule domain.RateRule) error {
	if m.addRuleErr != nil {
		return m.addRuleErr
	}
	m.rules = append(m.rules, rule)
	return nil
}

func (m *mockRateLimitUseCase) RemoveRule(_ context.Context, id string) error {
	if m.removeRuleErr != nil {
		return m.removeRuleErr
	}
	filtered := m.rules[:0]
	for _, r := range m.rules {
		if r.ID != id {
			filtered = append(filtered, r)
		}
	}
	m.rules = filtered
	return nil
}

func (m *mockRateLimitUseCase) GetRules(_ context.Context) ([]domain.RateRule, error) {
	return m.rules, m.getRulesErr
}

func (m *mockRateLimitUseCase) GetRulesForIP(_ context.Context, _ string) ([]domain.RateRule, error) {
	return m.rules, nil
}

func TestRateLimitHandler_GET_Returns200WithRules(t *testing.T) {
	uc := &mockRateLimitUseCase{
		rules: []domain.RateRule{{ID: "r1", MaxRPS: 10}},
	}
	h := NewRateLimitHandler(uc, nil)

	req := httptest.NewRequest(http.MethodGet, "/ratelimit", nil)
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var rules []domain.RateRule
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&rules))
	assert.Len(t, rules, 1)
	assert.Equal(t, "r1", rules[0].ID)
}

func TestRateLimitHandler_GET_RepoError_Returns500(t *testing.T) {
	uc := &mockRateLimitUseCase{getRulesErr: errors.New("db error")}
	h := NewRateLimitHandler(uc, nil)

	req := httptest.NewRequest(http.MethodGet, "/ratelimit", nil)
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestRateLimitHandler_POST_ValidRule_Returns201(t *testing.T) {
	uc := &mockRateLimitUseCase{}
	h := NewRateLimitHandler(uc, nil)

	rule := domain.RateRule{Target: domain.TargetIP, Value: "1.2.3.4", MaxRPS: 50}
	body, _ := json.Marshal(rule)

	req := httptest.NewRequest(http.MethodPost, "/ratelimit", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestRateLimitHandler_POST_InvalidJSON_Returns400(t *testing.T) {
	uc := &mockRateLimitUseCase{}
	h := NewRateLimitHandler(uc, nil)

	req := httptest.NewRequest(http.MethodPost, "/ratelimit", strings.NewReader("{bad"))
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRateLimitHandler_POST_MissingValue_Returns400(t *testing.T) {
	uc := &mockRateLimitUseCase{}
	h := NewRateLimitHandler(uc, nil)

	rule := domain.RateRule{Target: domain.TargetIP}
	body, _ := json.Marshal(rule)

	req := httptest.NewRequest(http.MethodPost, "/ratelimit", strings.NewReader(string(body)))
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRateLimitHandler_POST_InvalidTarget_Returns400(t *testing.T) {
	uc := &mockRateLimitUseCase{}
	h := NewRateLimitHandler(uc, nil)

	body := `{"target":"unknown","value":"1.2.3.4"}`
	req := httptest.NewRequest(http.MethodPost, "/ratelimit", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRateLimitHandler_POST_AddRuleError_Returns500(t *testing.T) {
	uc := &mockRateLimitUseCase{addRuleErr: errors.New("write error")}
	h := NewRateLimitHandler(uc, nil)

	rule := domain.RateRule{Target: domain.TargetIP, Value: "1.2.3.4"}
	body, _ := json.Marshal(rule)

	req := httptest.NewRequest(http.MethodPost, "/ratelimit", strings.NewReader(string(body)))
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestRateLimitHandler_ListHandler_UnknownMethod_Returns405(t *testing.T) {
	uc := &mockRateLimitUseCase{}
	h := NewRateLimitHandler(uc, nil)

	req := httptest.NewRequest(http.MethodPut, "/ratelimit", nil)
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestRateLimitHandler_DELETE_ValidID_Returns200(t *testing.T) {
	uc := &mockRateLimitUseCase{rules: []domain.RateRule{{ID: "rule-123"}}}
	h := NewRateLimitHandler(uc, nil)

	req := httptest.NewRequest(http.MethodDelete, "/ratelimit/rule-123", nil)
	rec := httptest.NewRecorder()
	h.RuleHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]string
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, "deleted", resp["status"])
}

func TestRateLimitHandler_DELETE_RemoveError_Returns500(t *testing.T) {
	uc := &mockRateLimitUseCase{removeRuleErr: errors.New("delete error")}
	h := NewRateLimitHandler(uc, nil)

	req := httptest.NewRequest(http.MethodDelete, "/ratelimit/rule-abc", nil)
	rec := httptest.NewRecorder()
	h.RuleHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestRateLimitHandler_RuleHandler_UnknownMethod_Returns405(t *testing.T) {
	uc := &mockRateLimitUseCase{}
	h := NewRateLimitHandler(uc, nil)

	req := httptest.NewRequest(http.MethodGet, "/ratelimit/some-id", nil)
	rec := httptest.NewRecorder()
	h.RuleHandler(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}
