package handlers

import (
	"net/http"
	"proxy/internal/domain"
)

type upstreamService interface {
	GetStatuses() []domain.UpstreamStatus
}

type UpstreamHandler struct {
	svc upstreamService
}

func NewUpstreamHandler(svc upstreamService) *UpstreamHandler {
	return &UpstreamHandler{svc: svc}
}

// GetStatuses godoc
// @Summary      Статус upstream-сервисов
// @Description  Возвращает актуальный статус всех настроенных upstream-сервисов (Java Backend, Redis, PostgreSQL и др.).
// @Description  Статус обновляется каждые 3 секунды через health-check (HTTP GET или TCP dial).
// @Description
// @Description  Поля ответа:
// @Description  - `status`: "healthy" | "degraded" | "down"
// @Description  - `latency`: последняя задержка (мс)
// @Description  - `latencyHistory`: последние 20 замеров (мс)
// @Description  - `errorRate`: процент ошибочных ответов за всё время
// @Description  - `requestsTotal`: общее число health-check запросов
// @Tags         Мониторинг
// @Produce      json
// @Success      200  {array}   domain.UpstreamStatus  "Список статусов upstream-сервисов"
// @Router       /api/v1/upstream [get]
func (h *UpstreamHandler) GetStatuses(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, h.svc.GetStatuses())
}
