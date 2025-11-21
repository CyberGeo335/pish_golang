package core

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/CyberGeo335/prak_ten/internal/http/middleware"
	"github.com/CyberGeo335/prak_ten/internal/platform/jwt"
	"github.com/CyberGeo335/prak_ten/internal/repo"
)

// userRepo — абстракция над репозиторием пользователей.
// Здесь мы используем конкретный тип repo.UserRecord,
// у которого должны быть экспортируемые поля ID, Email, Role.
type userRepo interface {
	CheckPassword(email, pass string) (repo.UserRecord, error)
	ByID(id int64) (repo.UserRecord, error)
}

type Service struct {
	repo      userRepo
	jwt       jwt.Validator
	blacklist *refreshBlacklist
	rl        *rateLimiter
}

func NewService(r userRepo, j jwt.Validator) *Service {
	return &Service{
		repo:      r,
		jwt:       j,
		blacklist: newRefreshBlacklist(),
		rl:        newRateLimiter(5, 5*time.Minute), // 5 попыток за 5 минут с одного IP
	}
}

// POST /api/v1/login
// Вход: {"Email":"...","Password":"..."}
// Выход: {"access_token":"...","refresh_token":"..."}
func (s *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ip := clientIP(r)
	if !s.rl.allow(ip) {
		httpError(w, r, http.StatusTooManyRequests, "too_many_attempts",
			"too many login attempts from this IP, please try later")
		return
	}

	var in struct {
		Email    string
		Password string
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Email == "" || in.Password == "" {
		httpError(w, r, http.StatusBadRequest, "invalid_credentials", "email and password are required")
		return
	}

	u, err := s.repo.CheckPassword(in.Email, in.Password)
	if err != nil {
		httpError(w, r, http.StatusUnauthorized, "unauthorized", "wrong email or password")
		return
	}

	access, err := s.jwt.SignAccess(u.ID, u.Email, u.Role)
	if err != nil {
		httpError(w, r, http.StatusInternalServerError, "token_error", "failed to issue access token")
		return
	}

	refresh, err := s.jwt.SignRefresh(u.ID, u.Email, u.Role)
	if err != nil {
		httpError(w, r, http.StatusInternalServerError, "token_error", "failed to issue refresh token")
		return
	}

	jsonOK(w, http.StatusOK, map[string]any{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

// POST /api/v1/refresh
// Вход:  {"refresh_token":"..."}
// Выход: новая пара {"access_token":"...","refresh_token":"..."}
func (s *Service) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var in struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.RefreshToken == "" {
		httpError(w, r, http.StatusBadRequest, "invalid_refresh", "refresh_token is required")
		return
	}

	if s.blacklist.isBlacklisted(in.RefreshToken) {
		httpError(w, r, http.StatusUnauthorized, "refresh_revoked", "refresh token has been revoked")
		return
	}

	claims, err := s.jwt.Parse(in.RefreshToken)
	if err != nil {
		httpError(w, r, http.StatusUnauthorized, "invalid_refresh", "cannot parse refresh token")
		return
	}

	if typ, _ := claims["typ"].(string); typ != "refresh" {
		httpError(w, r, http.StatusBadRequest, "invalid_refresh_type", "expected refresh token")
		return
	}

	subID := int64FromClaim(claims["sub"])
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	expUnix := int64FromClaim(claims["exp"])
	if expUnix > 0 {
		s.blacklist.add(in.RefreshToken, expUnix)
	}

	access, err := s.jwt.SignAccess(subID, email, role)
	if err != nil {
		httpError(w, r, http.StatusInternalServerError, "token_error", "failed to issue new access token")
		return
	}

	refresh, err := s.jwt.SignRefresh(subID, email, role)
	if err != nil {
		httpError(w, r, http.StatusInternalServerError, "token_error", "failed to issue new refresh token")
		return
	}

	jsonOK(w, http.StatusOK, map[string]any{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

// GET /api/v1/me
func (s *Service) MeHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.CtxClaimsKey).(map[string]any)
	if !ok {
		httpError(w, r, http.StatusInternalServerError, "claims_missing", "AuthN middleware is not configured")
		return
	}

	jsonOK(w, http.StatusOK, map[string]any{
		"id":    claims["sub"],
		"email": claims["email"],
		"role":  claims["role"],
	})
}

// GET /api/v1/users/{id}
// ABAC: user может читать только себя (id == sub), admin — любой.
func (s *Service) UserByIDHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.CtxClaimsKey).(map[string]any)
	if !ok {
		httpError(w, r, http.StatusInternalServerError, "claims_missing", nil)
		return
	}

	role, _ := claims["role"].(string)
	subID := int64FromClaim(claims["sub"])

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		httpError(w, r, http.StatusBadRequest, "bad_id", "id is required")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpError(w, r, http.StatusBadRequest, "bad_id", "id must be integer")
		return
	}

	if role == "user" && id != subID {
		httpError(w, r, http.StatusForbidden, "forbidden", "user can access only own profile")
		return
	}

	u, err := s.repo.ByID(id)
	if err != nil {
		httpError(w, r, http.StatusNotFound, "not_found", "user not found")
		return
	}

	jsonOK(w, http.StatusOK, map[string]any{
		"id":    u.ID,
		"email": u.Email,
		"role":  u.Role,
	})
}

// GET /api/v1/admin/stats
func (s *Service) AdminStats(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, http.StatusOK, map[string]any{
		"users":   2,
		"version": "1.0",
	})
}

type refreshBlacklist struct {
	mu     sync.RWMutex
	tokens map[string]int64 // token -> expUnix
}

func newRefreshBlacklist() *refreshBlacklist {
	return &refreshBlacklist{
		tokens: make(map[string]int64),
	}
}

func (b *refreshBlacklist) isBlacklisted(token string) bool {
	b.mu.RLock()
	exp, ok := b.tokens[token]
	b.mu.RUnlock()
	if !ok {
		return false
	}

	now := time.Now().Unix()
	if now > exp {
		b.mu.Lock()
		delete(b.tokens, token)
		b.mu.Unlock()
		return false
	}
	return true
}

func (b *refreshBlacklist) add(token string, exp int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tokens[token] = exp
}

// Простейший rate limiter: не более limit попыток за window с одного IP.
type rateLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time // ip -> timestamps
	limit    int
	window   time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		attempts: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (r *rateLimiter) allow(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	threshold := now.Add(-r.window)

	times := r.attempts[ip]
	// фильтруем старые попытки
	j := 0
	for _, t := range times {
		if t.After(threshold) {
			times[j] = t
			j++
		}
	}
	times = times[:j]

	if len(times) >= r.limit {
		r.attempts[ip] = times
		return false
	}

	times = append(times, now)
	r.attempts[ip] = times
	return true
}

type errorResponse struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

func jsonOK(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, r *http.Request, code int, msg string, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(errorResponse{
		Error:   msg,
		Details: details,
	})

	log.Printf("HTTP %d %s %s: %s (%v)", code, r.Method, r.URL.Path, msg, details)
}

func int64FromClaim(v any) int64 {
	switch t := v.(type) {
	case float64:
		return int64(t)
	case int64:
		return t
	case int:
		return int64(t)
	case string:
		id, err := strconv.ParseInt(t, 10, 64)
		if err == nil {
			return id
		}
	}
	return 0
}

func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		parts := strings.Split(fwd, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
