package handlers

import (
	"encoding/json"
	"github.com/CyberGeo335/prak_eleven/internal/core"
	"github.com/CyberGeo335/prak_eleven/internal/repo"
	"net/http"
)

type Handler struct {
	Repo *repo.NoteRepoMem
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var n core.Note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	id, _ := h.Repo.Create(n)
	n.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(n)
}
