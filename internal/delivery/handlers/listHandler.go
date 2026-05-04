package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"proxy/internal/usecase"
)

type Handler struct {
	listUseCase *usecase.HTTPUseCase
	logger      *slog.Logger
}

func NewHandler(listUseCase *usecase.HTTPUseCase, logger *slog.Logger) *Handler {
	return &Handler{listUseCase: listUseCase,
		logger: logger}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func (handler *Handler) getIp(w http.ResponseWriter, r *http.Request) {
	res, err := handler.listUseCase.GetAll()
	if err != nil {

		respondWithError(w, http.StatusInternalServerError, "Error getting whitelist")
		return
	}

	respondWithJSON(w, http.StatusOK, res)

}

func (handler *Handler) addIp(w http.ResponseWriter, r *http.Request) {
	var ip string
	if err := json.NewDecoder(r.Body).Decode(&ip); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if ip == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	res, err := handler.listUseCase.Add(r.Context(), ip)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error adding whitelist")
		return
	}

	handler.logger.Info("Успешно добавил новый ip в список")
	respondWithJSON(w, http.StatusOK, res)
}

func (handler *Handler) deleteIp(w http.ResponseWriter, r *http.Request) {
	var ip string
	if err := json.NewDecoder(r.Body).Decode(&ip); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if ip == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	res, err := handler.listUseCase.Remove(r.Context(), ip)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error removing whitelist")
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}

func (handler *Handler) ListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handler.getIp(w, r)
	case "POST":
		handler.addIp(w, r)
	case "DELETE":
		handler.deleteIp(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "")

	}
}
