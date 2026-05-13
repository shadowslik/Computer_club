package handlers

import (
	"encoding/json"
	"net/http"

	"proxy/internal/domain"
)

type CacheHandler struct {
	repo domain.CacheRepo
}

func NewCacheHandler(repo domain.CacheRepo) *CacheHandler {
	return &CacheHandler{repo: repo}
}

type invalidateResponse struct {
	Deleted int `json:"deleted"`
}

// GetStats godoc
// @Summary      Статистика кэша
// @Description  Возвращает текущее состояние кэша: размер, hit/miss, число вытеснений и занятую память.
// @Tags         Кэш
// @Produce      json
// @Success      200  {object}  domain.CacheStats  "Статистика кэша"
// @Router       /cache/stats [get]
func (h *CacheHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, h.repo.Stats())
}

// Flush godoc
// @Summary      Полная очистка кэша
// @Description  Удаляет все записи из кэша.
// @Tags         Кэш
// @Produce      json
// @Success      200  {object}  invalidateResponse  "Число удалённых записей"
// @Router       /cache [delete]
func (h *CacheHandler) Flush(w http.ResponseWriter, r *http.Request) {
	n := h.repo.Flush()
	respondWithJSON(w, http.StatusOK, invalidateResponse{Deleted: n})
}

// DeleteByKey godoc
// @Summary      Инвалидация по точному ключу
// @Description  Удаляет одну запись по точному ключу кэша (метод:хост:путь?параметры).
// @Tags         Кэш
// @Produce      json
// @Param        key  query  string  true  "Ключ кэша"
// @Success      200  {object}  invalidateResponse
// @Failure      400  {object}  map[string]string
// @Router       /cache/key [delete]
func (h *CacheHandler) DeleteByKey(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		respondWithError(w, http.StatusBadRequest, "query param 'key' is required")
		return
	}
	respondWithJSON(w, http.StatusOK, invalidateResponse{Deleted: h.repo.Delete(key)})
}

// DeleteByPrefix godoc
// @Summary      Инвалидация по префиксу
// @Description  Удаляет все записи, ключ которых начинается с указанного префикса.
// @Tags         Кэш
// @Produce      json
// @Param        prefix  query  string  true  "Префикс ключа"
// @Success      200  {object}  invalidateResponse
// @Failure      400  {object}  map[string]string
// @Router       /cache/prefix [delete]
func (h *CacheHandler) DeleteByPrefix(w http.ResponseWriter, r *http.Request) {
	prefix := r.URL.Query().Get("prefix")
	if prefix == "" {
		respondWithError(w, http.StatusBadRequest, "query param 'prefix' is required")
		return
	}
	respondWithJSON(w, http.StatusOK, invalidateResponse{Deleted: h.repo.DeleteByPrefix(prefix)})
}

// DeleteByRegex godoc
// @Summary      Инвалидация по регулярному выражению
// @Description  Удаляет все записи, ключ которых совпадает с регулярным выражением.
// @Tags         Кэш
// @Produce      json
// @Param        pattern  query  string  true  "Регулярное выражение"
// @Success      200  {object}  invalidateResponse
// @Failure      400  {object}  map[string]string
// @Router       /cache/regex [delete]
func (h *CacheHandler) DeleteByRegex(w http.ResponseWriter, r *http.Request) {
	pattern := r.URL.Query().Get("pattern")
	if pattern == "" {
		respondWithError(w, http.StatusBadRequest, "query param 'pattern' is required")
		return
	}
	n, err := h.repo.DeleteByRegex(pattern)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid regex: "+err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, invalidateResponse{Deleted: n})
}

// DeleteByTag godoc
// @Summary      Инвалидация по тегу
// @Description  Удаляет все записи, помеченные указанным тегом.
// @Tags         Кэш
// @Produce      json
// @Param        tag  query  string  true  "Тег"
// @Success      200  {object}  invalidateResponse
// @Failure      400  {object}  map[string]string
// @Router       /cache/tag [delete]
func (h *CacheHandler) DeleteByTag(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	if tag == "" {
		respondWithError(w, http.StatusBadRequest, "query param 'tag' is required")
		return
	}
	respondWithJSON(w, http.StatusOK, invalidateResponse{Deleted: h.repo.DeleteByTag(tag)})
}

// PostEvent godoc
// @Summary      Инвалидация по событию (4.4.3.1)
// @Description  Принимает именованное событие и инвалидирует кэш по тегам, префиксу или точным ключам.
// @Description  Используется для каскадной инвалидации при изменении данных в смежных сервисах.
// @Tags         Кэш
// @Accept       json
// @Produce      json
// @Param        body  body  domain.CacheEventRequest  true  "Описание события"
// @Success      200   {object}  domain.CacheEventResult
// @Failure      400   {object}  map[string]string
// @Router       /cache/events [post]
func (h *CacheHandler) PostEvent(w http.ResponseWriter, r *http.Request) {
	var req domain.CacheEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	total := 0
	for _, tag := range req.Tags {
		total += h.repo.DeleteByTag(tag)
	}
	if req.Prefix != "" {
		total += h.repo.DeleteByPrefix(req.Prefix)
	}
	for _, key := range req.Keys {
		total += h.repo.Delete(key)
	}
	respondWithJSON(w, http.StatusOK, domain.CacheEventResult{Deleted: total})
}

// DeleteByCondition godoc
// @Summary      Условная инвалидация
// @Description  Удаляет записи, удовлетворяющие условию.
// @Description  `stale_only` — только истёкшие записи.
// @Description  `tags` — хотя бы один из тегов должен совпасть (OR).
// @Tags         Кэш
// @Accept       json
// @Produce      json
// @Param        body  body  domain.CacheCondition  true  "Условие инвалидации"
// @Success      200   {object}  invalidateResponse
// @Failure      400   {object}  map[string]string
// @Router       /cache/condition [delete]
func (h *CacheHandler) DeleteByCondition(w http.ResponseWriter, r *http.Request) {
	var cond domain.CacheCondition
	if err := json.NewDecoder(r.Body).Decode(&cond); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, invalidateResponse{Deleted: h.repo.DeleteByCondition(cond)})
}
