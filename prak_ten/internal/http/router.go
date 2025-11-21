package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/CyberGeo335/prak_ten/internal/core"
	"github.com/CyberGeo335/prak_ten/internal/http/middleware"
	"github.com/CyberGeo335/prak_ten/internal/platform/config"
	"github.com/CyberGeo335/prak_ten/internal/platform/jwt"
	"github.com/CyberGeo335/prak_ten/internal/repo"
)

func Build(cfg config.Config) http.Handler {
	r := chi.NewRouter()

	// DI
	userRepo := repo.NewUserMem()
	jwtv := jwt.NewRS256(cfg.JWTSecret, cfg.JWTTTL)
	svc := core.NewService(userRepo, jwtv)

	// Публичные маршруты
	r.Post("/api/v1/login", svc.LoginHandler) // выдаёт access + refresh
	r.Post("/api/v1/refresh", svc.RefreshHandler)

	// Защищённые маршруты (admin + user)
	r.Group(func(priv chi.Router) {
		priv.Use(middleware.AuthN(jwtv))
		priv.Use(middleware.AuthZRoles("admin", "user"))

		priv.Get("/api/v1/me", svc.MeHandler)
		priv.Get("/api/v1/users/{id}", svc.UserByIDHandler) // ABAC
	})

	// Только для админов
	r.Group(func(admin chi.Router) {
		admin.Use(middleware.AuthN(jwtv))
		admin.Use(middleware.AuthZRoles("admin"))

		admin.Get("/api/v1/admin/stats", svc.AdminStats)
	})

	return r
}
