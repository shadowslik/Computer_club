package handlers

import (
	"net/http"
	"proxy/internal/domain"
	"strings"
)

type violatorsUseCase interface {
	GetAll() []domain.Violator
	Unban(ip string)
}

type ViolatorsHandler struct {
	uc violatorsUseCase
}

func NewViolatorsHandler(uc violatorsUseCase) *ViolatorsHandler {
	return &ViolatorsHandler{uc: uc}
}

// GetViolators godoc
// @Summary      Список нарушителей rate limiting
// @Description  Возвращает список IP-адресов, заблокированных за превышение лимитов. Записи автоматически удаляются по истечении BannedUntil.
// @Description
// @Description  Поля ответа:
// @Description  - `ip`: IP-адрес нарушителя
// @Description  - `endpoint`: URL, на котором было превышение
// @Description  - `blocked_at`: время блокировки (RFC3339)
// @Description  - `banned_until`: время снятия блокировки (RFC3339)
// @Tags         Rate Limiting
// @Produce      json
// @Success      200  {array}   domain.Violator  "Список нарушителей"
// @Router       /ratelimit/violators [get]
func (h *ViolatorsHandler) GetViolators(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, h.uc.GetAll())
}

// UnbanViolator godoc
// @Summary      Разблокировать IP
// @Description  Досрочно снимает блокировку с IP-адреса, нарушившего rate limit. IP передаётся в пути URL.
// @Tags         Rate Limiting
// @Produce      json
// @Param        ip   path      string  true  "IP-адрес для разблокировки"  example(1.2.3.4)
// @Success      200  {object}  map[string]string    "{\"status\":\"unbanned\"}"
// @Failure      400  {object}  domain.ErrorResponse "IP не указан в пути"
// @Failure      405  {object}  domain.ErrorResponse "Метод не поддерживается"
// @Router       /ratelimit/violators/{ip}/unban [post]
func (h *ViolatorsHandler) UnbanViolator(w http.ResponseWriter, r *http.Request) {
	ip := strings.TrimPrefix(r.URL.Path, "/ratelimit/violators/")
	ip = strings.TrimSuffix(ip, "/unban")
	if ip == "" {
		respondWithError(w, http.StatusBadRequest, "IP обязателен")
		return
	}
	h.uc.Unban(ip)
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "unbanned"})
}
