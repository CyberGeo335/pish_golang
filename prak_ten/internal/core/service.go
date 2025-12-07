package core

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	jwtlib "github.com/golang-jwt/jwt/v5"

	"github.com/CyberGeo335/prak_ten/internal/http/httputil"
	"github.com/CyberGeo335/prak_ten/internal/http/middleware"
	"github.com/CyberGeo335/prak_ten/internal/repo"
)

type userRepo interface {
	CheckPassword(email, pass string) (repo.UserRecord, error)
	ByID(id int64) (repo.UserRecord, error)
}

type tokenManager interface {
	SignAccess(userID int64, email, role string) (string, error)
	SignRefresh(userID int64, email, role string) (string, error)
	Parse(tokenStr string) (jwtlib.MapClaims, error)
}

type Service struct {
	repo         userRepo
	jwt          tokenManager
	refreshBL    *refreshBlacklist
	loginLimiter *loginLimiter
}

func NewService(r userRepo, tm tokenManager, loginMax int, loginWindow time.Duration) *Service {
	return &Service{
		repo:         r,
		jwt:          tm,
		refreshBL:    newRefreshBlacklist(),
		loginLimiter: newLoginLimiter(loginMax, loginWindow),
	}
}

// ---------- Handlers ----------

func (s *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ip := clientIP(r)
	if ok, retry := s.loginLimiter.Allow(ip); !ok {
		httputil.Error(w, http.StatusTooManyRequests, "too_many_requests", map[string]any{
			"retry_after_seconds": int(retry.Seconds()),
		})
		return
	}

	var in struct {
		Email    string `json:"Email"`
		Password string `json:"Password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Email == "" || in.Password == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_payload", "Email and Password are required")
		return
	}

	u, err := s.repo.CheckPassword(in.Email, in.Password)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "invalid_credentials", nil)
		return
	}

	access, err := s.jwt.SignAccess(u.ID, u.Email, u.Role)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "token_error", err.Error())
		return
	}
	refresh, err := s.jwt.SignRefresh(u.ID, u.Email, u.Role)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "token_error", err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"access_token":  access,
		"refresh_token": refresh,
		"token_type":    "Bearer",
	})
}

func (s *Service) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var in struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.RefreshToken == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_payload", "refresh_token is required")
		return
	}

	if s.refreshBL.IsRevoked(in.RefreshToken) {
		httputil.Error(w, http.StatusUnauthorized, "refresh_revoked", nil)
		return
	}

	claims, err := s.jwt.Parse(in.RefreshToken)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "invalid_refresh", err.Error())
		return
	}

	if typ, _ := claims["typ"].(string); typ != "refresh" {
		httputil.Error(w, http.StatusUnauthorized, "invalid_token_type", "expected refresh token")
		return
	}

	userID, ok := claimToInt64(claims["sub"])
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "invalid_sub", nil)
		return
	}
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	expUnix, ok := claimToInt64(claims["exp"])
	if !ok {
		httputil.Error(w, http.StatusInternalServerError, "invalid_exp", nil)
		return
	}

	// Отзываем старый refresh
	s.refreshBL.Revoke(in.RefreshToken, expUnix)

	// Новая пара токенов
	access, err := s.jwt.SignAccess(userID, email, role)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "token_error", err.Error())
		return
	}
	refresh, err := s.jwt.SignRefresh(userID, email, role)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "token_error", err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"access_token":  access,
		"refresh_token": refresh,
		"token_type":    "Bearer",
	})
}

func (s *Service) MeHandler(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]any{
		"id":    claims["sub"],
		"email": claims["email"],
		"role":  claims["role"],
	})
}

func (s *Service) AdminStats(w http.ResponseWriter, r *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]any{
		"users":   2,
		"version": "1.0",
	})
}

// ABAC: user может только своего /users/{id}, admin — любого
func (s *Service) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	role, _ := claims["role"].(string)
	subID, ok := claimToInt64(claims["sub"])
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "invalid_sub", nil)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_id", err.Error())
		return
	}

	if role == "user" && id != subID {
		httputil.Error(w, http.StatusForbidden, "forbidden", "user can access only own profile")
		return
	}

	u, err := s.repo.ByID(id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "not_found", nil)
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"id":    u.ID,
		"email": u.Email,
		"role":  u.Role,
	})
}

// ---------- вспомогательные структуры и функции ----------

type refreshBlacklist struct {
	mu     sync.RWMutex
	tokens map[string]int64 // token -> exp (unix)
}

func newRefreshBlacklist() *refreshBlacklist {
	return &refreshBlacklist{tokens: make(map[string]int64)}
}

func (b *refreshBlacklist) IsRevoked(tok string) bool {
	if tok == "" {
		return false
	}
	b.mu.RLock()
	exp, ok := b.tokens[tok]
	b.mu.RUnlock()
	if !ok {
		return false
	}
	if time.Now().Unix() > exp {
		b.mu.Lock()
		delete(b.tokens, tok)
		b.mu.Unlock()
		return false
	}
	return true
}

func (b *refreshBlacklist) Revoke(tok string, exp int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tokens[tok] = exp
}

type loginLimiter struct {
	mu       sync.Mutex
	attempts map[string]*clientCounter
	max      int
	window   time.Duration
}

type clientCounter struct {
	count int
	reset time.Time
}

func newLoginLimiter(max int, window time.Duration) *loginLimiter {
	return &loginLimiter{
		attempts: make(map[string]*clientCounter),
		max:      max,
		window:   window,
	}
}

func (l *loginLimiter) Allow(ip string) (bool, time.Duration) {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()

	c, ok := l.attempts[ip]
	if !ok || now.After(c.reset) {
		l.attempts[ip] = &clientCounter{
			count: 1,
			reset: now.Add(l.window),
		}
		return true, 0
	}
	if c.count >= l.max {
		return false, c.reset.Sub(now)
	}
	c.count++
	return true, 0
}

func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func claimToInt64(v any) (int64, bool) {
	switch t := v.(type) {
	case float64:
		return int64(t), true
	case int64:
		return t, true
	default:
		return 0, false
	}
}
