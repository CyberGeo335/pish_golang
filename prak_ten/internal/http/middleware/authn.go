package middleware

import (
	"context"
	"encoding/json"
	"github.com/CyberGeo335/prak_ten/internal/platform/jwt"
	"log"
	"net/http"
	"strings"
)

type CtxKey int

const CtxClaimsKey CtxKey = iota

type errorResponse struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

func jsonError(w http.ResponseWriter, r *http.Request, code int, msg string, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(errorResponse{
		Error:   msg,
		Details: details,
	})
	log.Printf("HTTP %d %s %s: %s (%v)", code, r.Method, r.URL.Path, msg, details)
}

// AuthN: достаём Bearer-токен, валидируем RS256, кладём клеймы в context.
// Принимаем только access-токены (typ == "access").
func AuthN(v jwt.Validator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h == "" || !strings.HasPrefix(h, "Bearer ") {
				jsonError(w, r, http.StatusUnauthorized, "unauthorized", "missing bearer token")
				return
			}

			raw := strings.TrimPrefix(h, "Bearer ")
			claims, err := v.Parse(raw)
			if err != nil {
				jsonError(w, r, http.StatusUnauthorized, "unauthorized", "invalid or expired token")
				return
			}

			if typ, _ := claims["typ"].(string); typ != "access" {
				jsonError(w, r, http.StatusUnauthorized, "unauthorized", "access token required")
				return
			}

			ctx := context.WithValue(r.Context(), CtxClaimsKey, map[string]any(claims))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
