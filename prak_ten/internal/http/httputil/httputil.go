package httputil

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Details any    `json:"details,omitempty"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Error(w http.ResponseWriter, status int, msg string, details any) {
	JSON(w, status, ErrorResponse{
		Error:   msg,
		Details: details,
	})
}
