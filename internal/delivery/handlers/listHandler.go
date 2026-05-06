package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type ipListUseCase interface {
	GetAll() (map[string]bool, error)
	IsAllowed(ctx context.Context, ip string) (bool, error)
	Add(ctx context.Context, ip string) (bool, error)
	Remove(ctx context.Context, ip string) (bool, error)
}

type Handler struct {
	ListUseCase ipListUseCase
	log         *zap.Logger
}

func NewHandler(listUseCase ipListUseCase, log *zap.Logger) *Handler {
	return &Handler{ListUseCase: listUseCase, log: log}
}

// getIp godoc
// @Summary      Получить IP-список
// @Description  Возвращает все записи из выбранного списка (IP-адреса, CIDR-подсети, диапазоны).
// @Tags         IP-управление
// @Produce      json
// @Success      200  {object}  map[string]bool  "Карта записей: ключ — IP/CIDR/диапазон, значение — true"
// @Failure      500  {object}  domain.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /ip/whitelist [get]
// @Router       /ip/blacklist [get]
// @Router       /ip/graylist  [get]
func (h *Handler) getIp(w http.ResponseWriter, r *http.Request) {
	res, err := h.ListUseCase.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "ошибка получения списка")
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}

// addIp godoc
// @Summary      Добавить запись в IP-список
// @Description  Добавляет IP-адрес, CIDR-подсеть (например 10.0.0.0/8) или диапазон (192.168.1.1-192.168.1.254) в список. Дубликаты игнорируются.
// @Tags         IP-управление
// @Accept       json
// @Produce      json
// @Param        ip  body      string  true  "IP-адрес, CIDR или диапазон"  example("192.168.1.0/24")
// @Success      200  {object}  map[string]bool       "true — запись добавлена, false — уже существует"
// @Failure      400  {object}  domain.ErrorResponse  "Невалидный IP/CIDR/диапазон или отсутствует тело запроса"
// @Failure      500  {object}  domain.ErrorResponse  "Ошибка сохранения"
// @Router       /ip/whitelist [post]
// @Router       /ip/blacklist [post]
// @Router       /ip/graylist  [post]
func (h *Handler) addIp(w http.ResponseWriter, r *http.Request) {
	var ip string
	if err := json.NewDecoder(r.Body).Decode(&ip); err != nil || ip == "" {
		respondWithError(w, http.StatusBadRequest, "невалидный запрос")
		return
	}

	res, err := h.ListUseCase.Add(r.Context(), ip)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "ошибка добавления IP")
		return
	}

	h.log.Info("IP added to list", zap.String("ip", ip))
	respondWithJSON(w, http.StatusOK, res)
}

// deleteIp godoc
// @Summary      Удалить запись из IP-списка
// @Description  Удаляет IP-адрес, CIDR-подсеть или диапазон из списка. Возвращает false, если записи не было.
// @Tags         IP-управление
// @Accept       json
// @Produce      json
// @Param        ip  body      string  true  "IP-адрес, CIDR или диапазон"  example("192.168.1.0/24")
// @Success      200  {object}  map[string]bool       "true — запись удалена, false — не найдена"
// @Failure      400  {object}  domain.ErrorResponse  "Невалидный IP или отсутствует тело запроса"
// @Failure      500  {object}  domain.ErrorResponse  "Ошибка сохранения"
// @Router       /ip/whitelist [delete]
// @Router       /ip/blacklist [delete]
// @Router       /ip/graylist  [delete]
func (h *Handler) deleteIp(w http.ResponseWriter, r *http.Request) {
	var ip string
	if err := json.NewDecoder(r.Body).Decode(&ip); err != nil || ip == "" {
		respondWithError(w, http.StatusBadRequest, "невалидный запрос")
		return
	}

	res, err := h.ListUseCase.Remove(r.Context(), ip)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "ошибка удаления IP")
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}

// ListHandler routes GET/POST/DELETE for an IP list.
func (h *Handler) ListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getIp(w, r)
	case http.MethodPost:
		h.addIp(w, r)
	case http.MethodDelete:
		h.deleteIp(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
	}
}
