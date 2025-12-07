package httpx

import (
	"github.com/CyberGeo335/prak_eleven/internal/http/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/v1/notes", h.CreateNote)
	return r
}
