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

type mockDenialsUseCase struct {
	stats []domain.DenialStat
}

func (m *mockDenialsUseCase) GetAll() []domain.DenialStat {
	return m.stats
}

func TestDenialsHandler_GetDenials_Returns200(t *testing.T) {
	uc := &mockDenialsUseCase{
		stats: []domain.DenialStat{
			{IP: "10.0.0.1", Denials: 5, LastDenied: time.Now(), Reason: "blacklist"},
			{IP: "10.0.0.2", Denials: 2, LastDenied: time.Now(), Reason: "graylist"},
		},
	}
	h := NewDenialsHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/ip/denials", nil)
	rec := httptest.NewRecorder()
	h.GetDenials(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var list []domain.DenialStat
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&list))
	assert.Len(t, list, 2)
}

func TestDenialsHandler_GetDenials_EmptyList(t *testing.T) {
	uc := &mockDenialsUseCase{}
	h := NewDenialsHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/ip/denials", nil)
	rec := httptest.NewRecorder()
	h.GetDenials(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var list []domain.DenialStat
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&list))
	assert.Empty(t, list)
}
