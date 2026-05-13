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
	"go.uber.org/zap"
)

type mockIpListUseCase struct {
	data            map[string]bool
	getAllErr       error
	addResult       bool
	addErr          error
	removeResult    bool
	removeErr       error
	isAllowedResult bool
	isAllowedErr    error
}

func newMockList() *mockIpListUseCase {
	return &mockIpListUseCase{data: make(map[string]bool)}
}

func (m *mockIpListUseCase) GetAll() (map[string]bool, error) {
	return m.data, m.getAllErr
}

func (m *mockIpListUseCase) IsAllowed(_ context.Context, _ string) (bool, error) {
	return m.isAllowedResult, m.isAllowedErr
}

func (m *mockIpListUseCase) Add(_ context.Context, ip string) (bool, error) {
	if m.addErr != nil {
		return false, m.addErr
	}
	m.data[ip] = true
	return m.addResult, nil
}

func (m *mockIpListUseCase) Remove(_ context.Context, ip string) (bool, error) {
	if m.removeErr != nil {
		return false, m.removeErr
	}
	delete(m.data, ip)
	return m.removeResult, nil
}

func newListHandler(uc *mockIpListUseCase) *Handler {
	return NewHandler(uc, zap.NewNop(), nil, "")
}

func jsonBody(v interface{}) *strings.Reader {
	b, _ := json.Marshal(v)
	return strings.NewReader(string(b))
}

func TestListHandler_GET_Returns200WithMap(t *testing.T) {
	uc := newMockList()
	uc.data["1.2.3.4"] = true
	h := newListHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/ip/whitelist", nil)
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result map[string]bool
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&result))
	assert.True(t, result["1.2.3.4"])
}

func TestListHandler_GET_RepoError_Returns500(t *testing.T) {
	uc := newMockList()
	uc.getAllErr = errors.New("storage down")
	h := newListHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/ip/whitelist", nil)
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListHandler_POST_ValidIP_Returns200(t *testing.T) {
	uc := newMockList()
	uc.addResult = true
	h := newListHandler(uc)

	req := httptest.NewRequest(http.MethodPost, "/ip/whitelist", jsonBody("10.0.0.1"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestListHandler_POST_EmptyBody_Returns400(t *testing.T) {
	uc := newMockList()
	h := newListHandler(uc)

	req := httptest.NewRequest(http.MethodPost, "/ip/whitelist", jsonBody(""))
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListHandler_POST_InvalidBody_Returns400(t *testing.T) {
	uc := newMockList()
	h := newListHandler(uc)

	req := httptest.NewRequest(http.MethodPost, "/ip/whitelist", strings.NewReader("{not json"))
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListHandler_POST_AddError_Returns500(t *testing.T) {
	uc := newMockList()
	uc.addErr = errors.New("db write error")
	h := newListHandler(uc)

	req := httptest.NewRequest(http.MethodPost, "/ip/whitelist", jsonBody("5.5.5.5"))
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListHandler_DELETE_ValidIP_Returns200(t *testing.T) {
	uc := newMockList()
	uc.data["6.6.6.6"] = true
	uc.removeResult = true
	h := newListHandler(uc)

	req := httptest.NewRequest(http.MethodDelete, "/ip/whitelist", jsonBody("6.6.6.6"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestListHandler_DELETE_EmptyBody_Returns400(t *testing.T) {
	uc := newMockList()
	h := newListHandler(uc)

	req := httptest.NewRequest(http.MethodDelete, "/ip/whitelist", jsonBody(""))
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListHandler_DELETE_RemoveError_Returns500(t *testing.T) {
	uc := newMockList()
	uc.removeErr = errors.New("remove failed")
	h := newListHandler(uc)

	req := httptest.NewRequest(http.MethodDelete, "/ip/whitelist", jsonBody("7.7.7.7"))
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListHandler_UnknownMethod_Returns405(t *testing.T) {
	uc := newMockList()
	h := newListHandler(uc)

	req := httptest.NewRequest(http.MethodPatch, "/ip/whitelist", nil)
	rec := httptest.NewRecorder()
	h.ListHandler(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}
