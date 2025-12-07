package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/CyberGeo335/prak_ten/internal/http/httputil"
	jwtplatform "github.com/CyberGeo335/prak_ten/internal/platform/jwt"
)

type ctxKey int

const ctxClaimsKey ctxKey = iota

func AuthN(v jwtplatform.Validator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h == "" || !strings.HasPrefix(h, "Bearer ") {
				httputil.Error(w, http.StatusUnauthorized, "unauthorized", "missing Bearer token")
				return
			}
			raw := strings.TrimPrefix(h, "Bearer ")
			claims, err := v.Parse(raw)
			if err != nil {
				httputil.Error(w, http.StatusUnauthorized, "unauthorized", err.Error())
				return
			}
			if typ, _ := claims["typ"].(string); typ != "access" {
				httputil.Error(w, http.StatusUnauthorized, "unauthorized", "access token required")
				return
			}
			ctx := context.WithValue(r.Context(), ctxClaimsKey, map[string]any(claims))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ClaimsFromContext(ctx context.Context) map[string]any {
	claims, _ := ctx.Value(ctxClaimsKey).(map[string]any)
	return claims
}
