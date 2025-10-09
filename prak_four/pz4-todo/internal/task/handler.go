package task

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Handler инкапсулирует зависимости для работы HTTP-слоя.
type Handler struct {
	repo *Repo
}

// NewHandler конструирует обработчик задач.
func NewHandler(repo *Repo) *Handler { return &Handler{repo: repo} }

// Routes описывает под-маршруты для ресурса /tasks.
// База: /api/v1/tasks
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)          // GET    /api/v1/tasks
	r.Post("/", h.create)       // POST   /api/v1/tasks
	r.Get("/{id}", h.get)       // GET    /api/v1/tasks/{id}
	r.Put("/{id}", h.update)    // PUT    /api/v1/tasks/{id}
	r.Delete("/{id}", h.delete) // DELETE /api/v1/tasks/{id}
	return r
}

// list возвращает список задач с простыми фильтрами и пагинацией.
// Query: ?title=foo&page=1&limit=10
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	title := strings.TrimSpace(q.Get("title"))

	page := 1
	if p := q.Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}

	limit := 10
	if l := q.Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			if val > 100 {
				val = 100
			}
			limit = val
		}
	}

	list, err := h.repo.List(title, page, limit)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// get возвращает одну задачу по id.
func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	t, err := h.repo.Get(id)
	if err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

// create создаёт задачу. Ожидается JSON: {"title":"..."}.
type createReq struct {
	Title string `json:"title"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	if !requireJSON(w, r) {
		return
	}

	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Title) == "" {
		httpError(w, http.StatusBadRequest, "invalid json: require non-empty title")
		return
	}

	if !validateTitle(w, req.Title) {
		return
	}

	t, err := h.repo.Create(req.Title)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

// update обновляет заголовок и статус done. Ожидается JSON: {"title":"...","done":true}.
type updateReq struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	if !requireJSON(w, r) {
		return
	}

	id, bad := parseID(w, r)
	if bad {
		return
	}

	var req updateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Title) == "" {
		httpError(w, http.StatusBadRequest, "invalid json: require non-empty title")
		return
	}

	if !validateTitle(w, req.Title) {
		return
	}

	t, err := h.repo.Update(id, req.Title, req.Done)
	if err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

// delete удаляет задачу по id.
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	if err := h.repo.Delete(id); err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// requireJSON проверяет корректный Content-Type для JSON-запросов.
func requireJSON(w http.ResponseWriter, r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	if ct != "" && !strings.Contains(ct, "application/json") {
		httpError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return false
	}
	return true
}

// validateTitle проверяет длину и непустоту заголовка.
func validateTitle(w http.ResponseWriter, title string) bool {
	title = strings.TrimSpace(title)
	if title == "" {
		httpError(w, http.StatusBadRequest, "invalid title")
		return false
	}
	if len(title) < 3 {
		httpError(w, http.StatusBadRequest, "too short title")
		return false
	}
	if len(title) > 100 {
		httpError(w, http.StatusBadRequest, "too long title")
		return false
	}
	return true
}

// parseID достаёт параметр {id} из пути.
func parseID(w http.ResponseWriter, r *http.Request) (string, bool) {
	raw := chi.URLParam(r, "id")
	if raw == "" {
		httpError(w, http.StatusBadRequest, "invalid id")
		return "", true
	}
	return raw, false
}

// writeJSON и httpError — вспомогательные функции ответа.
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
