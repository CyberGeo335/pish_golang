package api

import (
	"encoding/json"
	"errors"
	"example.com/pz3_http/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

type Handlers struct {
	Store *storage.MemoryStore
}

func NewHandlers(store *storage.MemoryStore) *Handlers {
	return &Handlers{Store: store}
}

// GET /tasks
func (h *Handlers) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.Store.List()

	// Простые фильтры через query: ?q=text
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q != "" {
		filtered := tasks[:0]
		for _, t := range tasks {
			if strings.Contains(strings.ToLower(t.Title), strings.ToLower(q)) {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}

	JSON(w, http.StatusOK, tasks)
}

type createTaskRequest struct {
	Title string `json:"title"`
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "" && !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		BadRequest(w, "Content-Type must be application/json")
		return
	}

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid json: "+err.Error())
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		BadRequest(w, "title is required")
		return
	}

	if len(req.Title) < 3 {
		Unprocessable(w, "title is too short")
		return
	}

	if len(req.Title) > 140 {
		Unprocessable(w, "title is too long")
		return
	}

	t := h.Store.Create(req.Title)
	JSON(w, http.StatusCreated, t)
}

func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	// Ожидаем путь вида </tasks/123>
	id, err := extractIDFromPath(w, r)
	if err != nil {
		return
	}

	t, err := h.Store.Get(id)
	if err != nil {
		if err.Error() == "not found" {
			NotFound(w, "task not found")
			return
		}
		Internal(w, "unexpected error")
		return
	}
	JSON(w, http.StatusOK, t)
}

type updateTaskRequest struct {
	Done bool `json:"done"`
}

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "" && !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		BadRequest(w, "Content-Type must be application/json")
		return
	}

	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid json: "+err.Error())
		return
	}

	id, err := extractIDFromPath(w, r)
	if err != nil {
		return
	}

	t := h.Store.Update(id, storage.TaskUpdatePayload{Done: req.Done})
	if t == nil {
		NotFound(w, "task not found")
		return
	}

	JSON(w, http.StatusOK, t)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(w, r)
	if err != nil {
		return
	}

	h.Store.Delete(id)
	w.WriteHeader(http.StatusOK)
}

func extractIDFromPath(w http.ResponseWriter, r *http.Request) (int64, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		NotFound(w, "invalid path")
		return 0, errors.New("invalid path")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		BadRequest(w, "invalid id")
		return 0, errors.New("invalid id")
	}
	return id, nil
}
