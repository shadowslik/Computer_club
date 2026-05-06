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

// ListHandler godoc
// @Summary      Управление правилами rate limiting
// @Description  GET — возвращает все правила. POST — создаёт новое правило для IP или подсети.
// @Description
// @Description  Поля правила (все лимиты -1 = отключён):
// @Description  - `target`: "ip" или "subnet"
// @Description  - `value`: конкретный IP или CIDR (например 10.0.0.0/8)
// @Description  - `max_rps` / `max_rpm` / `max_rph` / `max_rpd`: лимиты запросов
// @Description  - `max_upload_bps` / `max_download_bps`: лимиты пропускной способности (байт/с)
// @Description  - `max_traffic_day`: дневной лимит трафика (байт)
// @Description  - `max_concurrent`: макс. одновременных соединений
// @Description  - `max_new_per_second`: макс. новых соединений в секунду
// @Tags         Rate Limiting
// @Accept       json
// @Produce      json
// @Param        rule  body      domain.RateRule  false  "Правило rate limiting (только для POST)"
// @Success      200   {array}   domain.RateRule       "Список всех правил"
// @Success      201   {object}  map[string]string     "Правило создано: {\"status\":\"ok\"}"
// @Failure      400   {object}  domain.ErrorResponse  "Невалидный JSON или отсутствует поле value/target"
// @Failure      405   {object}  domain.ErrorResponse  "Метод не поддерживается"
// @Failure      500   {object}  domain.ErrorResponse  "Ошибка сохранения"
// @Router       /ratelimit [get]
// @Router       /ratelimit [post]
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

// RuleHandler godoc
// @Summary      Удалить правило rate limiting
// @Description  Удаляет правило по его UUID. ID правила берётся из пути: /ratelimit/{id}.
// @Tags         Rate Limiting
// @Produce      json
// @Param        id  path      string  true  "UUID правила"  example(550e8400-e29b-41d4-a716-446655440000)
// @Success      200  {object}  map[string]string    "Правило удалено: {\"status\":\"deleted\"}"
// @Failure      400  {object}  domain.ErrorResponse "ID правила не указан"
// @Failure      405  {object}  domain.ErrorResponse "Метод не поддерживается"
// @Failure      500  {object}  domain.ErrorResponse "Ошибка удаления"
// @Router       /ratelimit/{id} [delete]
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
