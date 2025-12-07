package middleware

import (
	"net/http"

	"github.com/CyberGeo335/prak_ten/internal/http/httputil"
)

func AuthZRoles(allowed ...string) func(http.Handler) http.Handler {
	set := map[string]struct{}{}
	for _, a := range allowed {
		set[a] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := ClaimsFromContext(r.Context())
			if claims == nil {
				httputil.Error(w, http.StatusUnauthorized, "unauthorized", nil)
				return
			}
			role, _ := claims["role"].(string)
			if _, ok := set[role]; !ok {
				httputil.Error(w, http.StatusForbidden, "forbidden", "role not allowed")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
