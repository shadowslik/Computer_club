package handlers

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"proxy/pkg/utils"
	"strings"

	"go.uber.org/zap"
)

type ipCheckUseCase interface {
	IsAllowed(ctx context.Context, ip string) (bool, error)
}

type denialsRecorder interface {
	Record(ip, reason string)
}

type CheckIpResponse struct {
	IP      string `json:"ip"      example:"192.168.1.1"`
	List    string `json:"list"    example:"white"          enums:"white,black,gray,none"`
	Message string `json:"message" example:"доступ разрешён"`
}

type CheckIpHandler struct {
	whiteList ipCheckUseCase
	blackList ipCheckUseCase
	grayList  ipCheckUseCase
	log       *zap.Logger
	denials   denialsRecorder
}

func NewCheckIpHandler(
	whiteList, blackList, grayList ipCheckUseCase,
	log *zap.Logger,
	denials denialsRecorder,
) *CheckIpHandler {
	return &CheckIpHandler{
		whiteList: whiteList,
		blackList: blackList,
		grayList:  grayList,
		log:       log,
		denials:   denials,
	}
}

// CheckIp godoc
// @Summary      Проверить принадлежность IP к спискам
// @Description  Определяет, в каком списке находится IP-адрес (black/gray/white/none).
// @Description  Порядок проверки: blacklist → graylist → whitelist → none.
// @Description  IP можно передать через query-параметр, тело запроса (JSON-строка) или он будет взят из RemoteAddr.
// @Tags         IP-управление
// @Produce      json
// @Param        ip   query     string  false  "Проверяемый IP-адрес"  example(192.168.1.1)
// @Success      200  {object}  CheckIpResponse       "Результат проверки"
// @Failure      400  {object}  domain.ErrorResponse  "Невалидный IP-адрес"
// @Failure      500  {object}  domain.ErrorResponse  "Внутренняя ошибка"
// @Router       /ip/check [get]
func (h *CheckIpHandler) CheckIp(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")

	if ip == "" && r.Body != nil {
		var s string
		if err := json.NewDecoder(r.Body).Decode(&s); err == nil {
			ip = strings.TrimSpace(s)
		}
	}

	if ip == "" {
		var err error
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = utils.GetRealIP(r)
		}
	}

	if ip == "" || net.ParseIP(ip) == nil {
		respondWithError(w, http.StatusBadRequest, "невалидный IP-адрес")
		return
	}

	list, message := h.resolveList(r.Context(), ip)

	if list != "white" {
		h.log.Info("IP check",
			zap.String("ip", ip),
			zap.String("list", list),
			zap.String("url", r.RequestURI),
		)
	}

	respondWithJSON(w, http.StatusOK, CheckIpResponse{IP: ip, List: list, Message: message})
}

func (h *CheckIpHandler) resolveList(ctx context.Context, ip string) (string, string) {
	if ok, _ := h.blackList.IsAllowed(ctx, ip); ok {
		h.denials.Record(ip, "blacklist")
		return "black", "доступ запрещён — IP в чёрном списке"
	}
	if ok, _ := h.grayList.IsAllowed(ctx, ip); ok {
		h.denials.Record(ip, "graylist")
		return "gray", "необходимо пройти CAPTCHA"
	}
	if ok, _ := h.whiteList.IsAllowed(ctx, ip); ok {
		return "white", "доступ разрешён"
	}
	return "none", "IP не найден ни в одном списке"
}
