package handlers

import (
	"net/http"
	"proxy/internal/domain"
)

type denialsUseCase interface {
	GetAll() []domain.DenialStat
}

type DenialsHandler struct {
	uc denialsUseCase
}

func NewDenialsHandler(uc denialsUseCase) *DenialsHandler {
	return &DenialsHandler{uc: uc}
}

// GetDenials godoc
// @Summary      Статистика отказов доступа
// @Description  Возвращает агрегированную статистику отказов: сколько раз каждый IP был заблокирован и по какой причине.
// @Description  Источники записей: blacklist (403), graylist (429/CAPTCHA), not_in_whitelist (403 при default_deny=true).
// @Description  Записи накапливаются с момента запуска сервера (in-memory хранилище).
// @Description
// @Description  Поля ответа:
// @Description  - `ip`: IP-адрес
// @Description  - `denials`: суммарное количество отказов
// @Description  - `last_denied`: время последнего отказа (RFC3339)
// @Description  - `reason`: последняя причина ("blacklist" | "graylist" | "not_in_whitelist")
// @Tags         IP-управление
// @Produce      json
// @Success      200  {array}   domain.DenialStat  "Список записей об отказах"
// @Router       /ip/denials [get]
func (h *DenialsHandler) GetDenials(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, h.uc.GetAll())
}
