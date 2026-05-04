package handlers

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"proxy/internal/usecase"
	"strings"
)

type CheckIpResponse struct {
	Ip      string `json:"ip"`
	List    string `json:"list"`
	Message string `json:"message"`
}

type CheckIpHandler struct {
	whiteList *usecase.HTTPUseCase
	blackList *usecase.HTTPUseCase
	grayList  *usecase.HTTPUseCase
	logger    *slog.Logger
}

func NewCheckIpHandler(
	whiteList *usecase.HTTPUseCase,
	blackList *usecase.HTTPUseCase,
	grayList *usecase.HTTPUseCase,
	logger *slog.Logger) *CheckIpHandler {
	return &CheckIpHandler{
		whiteList: whiteList,
		blackList: blackList,
		grayList:  grayList,
		logger:    logger,
	}
}

func (h *CheckIpHandler) CheckIp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "only GET")
		return
	}

	var err error
	var ip string

	ip = r.URL.Query().Get("ip")

	if ip == "" && r.Body != nil {
		err = json.NewDecoder(r.Body).Decode(&ip)
		if err != nil {

			respondWithError(w, http.StatusBadRequest, "invalid json")
			return
		}

		ip = strings.TrimSpace(ip)
	}

	if ip == "" {

		ip, _, err = net.SplitHostPort(r.RemoteAddr)

		if err != nil {
			h.logger.Error(err.Error())
			respondWithError(w, http.StatusInternalServerError, "invalid remote address")
			return
		}

		if ip == "" {
			ip = GetRealIp(r)
		}
	}

	if ip == "" {
		respondWithError(w, http.StatusInternalServerError, "invalid remote address")
		return
	}

	if net.ParseIP(ip) == nil {
		respondWithError(w, http.StatusInternalServerError, "invalid remote address")
		return
	}

	var list string
	var message string

	list, message = h.GetListIp(r, ip)

	if list == "black" || list == "none" {
		h.logger.Info(
			"IP", ip,
			"list", list,
			"URL", r.RequestURI,
			"message", message,
		)
	}

	if list == "white" || list == "none" {
		h.logger.Info(
			"IP", ip,
			"list", list,
			"URL", r.RequestURI,
			"message", message,
		)
	}

	response := CheckIpResponse{
		Ip:      ip,
		List:    list,
		Message: message,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func GetRealIp(r *http.Request) string {

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if xrip := r.Header.Get("X-Real-Ip"); xrip != "" {
		return xrip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func (h *CheckIpHandler) GetListIp(r *http.Request, ip string) (string, string) {

	var list string
	var message string

	allowed, err := h.blackList.IsAllowed(r.Context(), ip)
	if err == nil && allowed {
		list = "black"
		message = "Доступ запрещён Ip в черном списке"
	} else {
		allowed, err = h.grayList.IsAllowed(r.Context(), ip)

		if err == nil && allowed {
			list = "gray"
			message = "Необходимо пройти капчу"
		} else {
			allowed, err = h.whiteList.IsAllowed(r.Context(), ip)
			if err == nil && allowed {
				list = "white"
				message = "Доступ разрешён"
			} else {
				list = "none"
				message = "Данный IP не найден"
			}
		}
	}

	return list, message
}

func (h *CheckIpHandler) HandlerCheck(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.CheckIp(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "")

	}
}
