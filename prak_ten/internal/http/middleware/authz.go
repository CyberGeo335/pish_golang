package middleware

import "net/http"

// AuthZRoles — простой RBAC: проверяет, что роль пользователя входит в allowed.
func AuthZRoles(allowed ...string) func(http.Handler) http.Handler {
	set := map[string]struct{}{}
	for _, a := range allowed {
		set[a] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, _ := r.Context().Value(CtxClaimsKey).(map[string]any)
			role, _ := claims["role"].(string)
			if _, ok := set[role]; !ok {
				jsonError(w, r, http.StatusForbidden, "forbidden", "insufficient role")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
