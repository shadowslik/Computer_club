package handlers

import (
	"encoding/json"
	"net/http"
	"proxy/pkg/metrics"
)

type MetricsHandler struct {
	collector *metrics.Collector
}

func NewMetricsHandler(c *metrics.Collector) *MetricsHandler {
	return &MetricsHandler{collector: c}
}

// HandleMetrics godoc
// @Summary      Метрики сервера
// @Description  Возвращает агрегированные метрики с момента запуска: количество запросов, ошибок, трафик, RPS, латентность, активные соединения.
// @Tags         Мониторинг
// @Produce      json
// @Success      200  {object}  metrics.Snapshot  "Снимок метрик"
// @Router       /metrics [get]
func (h *MetricsHandler) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	snap := h.collector.Snapshot()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snap)
}

// HandleHealth godoc
// @Summary      Health check
// @Description  Возвращает статус сервера, uptime, текущий RPS и количество активных соединений. Используется для liveness/readiness проб.
// @Tags         Мониторинг
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "status, uptime_s, rps, conns"
// @Router       /health [get]
func (h *MetricsHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	snap := h.collector.Snapshot()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":   "ok",
		"uptime_s": snap.UptimeSec,
		"rps":      snap.RPS,
		"conns":    snap.ActiveConns,
	})
}
