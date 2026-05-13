package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proxy/internal/domain"
)

type mockViolatorsUseCase struct {
	violators  []domain.Violator
	unbannedIP string
}

func (m *mockViolatorsUseCase) GetAll() []domain.Violator {
	return m.violators
}

func (m *mockViolatorsUseCase) Unban(ip string) {
	m.unbannedIP = ip
}

func TestViolatorsHandler_GetViolators_Returns200(t *testing.T) {
	uc := &mockViolatorsUseCase{
		violators: []domain.Violator{
			{IP: "1.2.3.4", Endpoint: "/api", Requests: 200, Limit: 100, BlockedAt: time.Now()},
		},
	}
	h := NewViolatorsHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/ratelimit/violators", nil)
	rec := httptest.NewRecorder()
	h.GetViolators(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var list []domain.Violator
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&list))
	assert.Len(t, list, 1)
	assert.Equal(t, "1.2.3.4", list[0].IP)
}

func TestViolatorsHandler_GetViolators_EmptyList(t *testing.T) {
	uc := &mockViolatorsUseCase{}
	h := NewViolatorsHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/ratelimit/violators", nil)
	rec := httptest.NewRecorder()
	h.GetViolators(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestViolatorsHandler_UnbanViolator_ValidIP_Returns200(t *testing.T) {
	uc := &mockViolatorsUseCase{}
	h := NewViolatorsHandler(uc)

	req := httptest.NewRequest(http.MethodPost, "/ratelimit/violators/5.5.5.5/unban", nil)
	rec := httptest.NewRecorder()
	h.UnbanViolator(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "5.5.5.5", uc.unbannedIP)

	var resp map[string]string
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, "unbanned", resp["status"])
}

func TestViolatorsHandler_UnbanViolator_MissingIP_Returns400(t *testing.T) {
	uc := &mockViolatorsUseCase{}
	h := NewViolatorsHandler(uc)

	req := httptest.NewRequest(http.MethodPost, "/ratelimit/violators//unban", nil)
	rec := httptest.NewRecorder()
	h.UnbanViolator(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
