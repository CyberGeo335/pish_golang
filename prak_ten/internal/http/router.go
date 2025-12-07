package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/CyberGeo335/prak_ten/internal/core"
	"github.com/CyberGeo335/prak_ten/internal/http/middleware"
	"github.com/CyberGeo335/prak_ten/internal/platform/config"
	jwtplatform "github.com/CyberGeo335/prak_ten/internal/platform/jwt"
	"github.com/CyberGeo335/prak_ten/internal/repo"
)

func Build(cfg config.Config) (http.Handler, error) {
	r := chi.NewRouter()

	r.Use(middleware.Logging)

	userRepo := repo.NewUserMem()
	jwtm, err := jwtplatform.NewRS256(cfg)
	if err != nil {
		return nil, err
	}

	svc := core.NewService(userRepo, jwtm, cfg.RateLimitLoginMax, cfg.RateLimitLoginWindow)

	// Публичные маршруты
	r.Post("/api/v1/login", svc.LoginHandler)
	r.Post("/api/v1/refresh", svc.RefreshHandler)

	// Защищённые маршруты (user + admin)
	r.Group(func(priv chi.Router) {
		priv.Use(middleware.AuthN(jwtm))
		priv.Use(middleware.AuthZRoles("admin", "user"))
		priv.Get("/api/v1/me", svc.MeHandler)
		priv.Get("/api/v1/users/{id}", svc.GetUserHandler)
	})

	// Только админы
	r.Group(func(admin chi.Router) {
		admin.Use(middleware.AuthN(jwtm))
		admin.Use(middleware.AuthZRoles("admin"))
		admin.Get("/api/v1/admin/stats", svc.AdminStats)
	})

	return r, nil
}
