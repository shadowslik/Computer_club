package handlers

import (
	"encoding/json"
	"net/http"
	"proxy/internal/domain"
	"strings"
)

type RateLimitHandler struct {
	uc domain.RateLimitUseCase
}

func NewRateLimitHandler(uc domain.RateLimitUseCase) *RateLimitHandler {
	return &RateLimitHandler{uc: uc}
}

// ListHandler — GET /ratelimit → список всех правил
//
//	POST /ratelimit → добавить правило
func (h *RateLimitHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getRules(w, r)
	case http.MethodPost:
		h.addRule(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
	}
}

// RuleHandler — DELETE /ratelimit/{id} → удалить правило по ID
func (h *RateLimitHandler) RuleHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		h.removeRule(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
	}
}

func (h *RateLimitHandler) getRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.uc.GetRules(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "ошибка получения правил")
		return
	}
	respondWithJSON(w, http.StatusOK, rules)
}

func (h *RateLimitHandler) addRule(w http.ResponseWriter, r *http.Request) {
	var rule domain.RateRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		respondWithError(w, http.StatusBadRequest, "невалидный JSON")
		return
	}

	// базовая валидация
	if rule.Value == "" {
		respondWithError(w, http.StatusBadRequest, "поле value обязательно")
		return
	}
	if rule.Target != domain.TargetIP && rule.Target != domain.TargetSubnet {
		respondWithError(w, http.StatusBadRequest, "target должен быть ip или subnet")
		return
	}

	if err := h.uc.AddRule(r.Context(), rule); err != nil {
		respondWithError(w, http.StatusInternalServerError, "ошибка добавления правила")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

func (h *RateLimitHandler) removeRule(w http.ResponseWriter, r *http.Request) {
	// достаём ID из пути: /ratelimit/some-uuid
	id := strings.TrimPrefix(r.URL.Path, "/ratelimit/")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "ID правила обязателен")
		return
	}

	if err := h.uc.RemoveRule(r.Context(), id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "ошибка удаления правила")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
