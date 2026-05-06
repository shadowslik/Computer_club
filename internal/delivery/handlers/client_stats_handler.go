package handlers

import (
	"net/http"
	"proxy/internal/domain"
	"strconv"
)

type clientStatsUseCase interface {
	GetTop(limit int, sortBy string) []domain.ClientStat
}

type ClientStatsHandler struct {
	uc clientStatsUseCase
}

func NewClientStatsHandler(uc clientStatsUseCase) *ClientStatsHandler {
	return &ClientStatsHandler{uc: uc}
}

// GetTopClients godoc
// @Summary      Топ клиентов
// @Description  Возвращает список клиентских IP, отсортированных по выбранному критерию. Данные накапливаются через middleware с момента старта.
// @Description
// @Description  Параметр `sort`:
// @Description  - `requests` (по умолчанию) — по количеству запросов
// @Description  - `traffic` — по объёму трафика (МБ)
// @Description  - `errors` — по проценту ошибок
// @Description
// @Description  Поля ответа:
// @Description  - `ip`: IP-адрес клиента
// @Description  - `requests`: суммарное количество запросов
// @Description  - `traffic_mb`: суммарный трафик (МБ)
// @Description  - `avg_latency_ms`: средняя задержка (мс)
// @Description  - `error_rate`: процент ответов 4xx/5xx
// @Description  - `last_seen`: время последнего запроса (HH:MM:SS)
// @Tags         Мониторинг
// @Produce      json
// @Param        sort   query     string  false  "Поле сортировки: requests, traffic, errors"  Enums(requests, traffic, errors)  default(requests)
// @Param        limit  query     integer false  "Максимальное количество записей"              default(20)
// @Success      200    {array}   domain.ClientStat    "Топ клиентов"
// @Failure      400    {object}  domain.ErrorResponse "Невалидный параметр"
// @Router       /top-clients [get]
func (h *ClientStatsHandler) GetTopClients(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "requests"
	}
	limit := 20
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
		limit = l
	}
	respondWithJSON(w, http.StatusOK, h.uc.GetTop(limit, sortBy))
}
