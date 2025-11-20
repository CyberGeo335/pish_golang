package app

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/CyberGeo335/prak_nine/internal/http/handlers"
	"github.com/CyberGeo335/prak_nine/internal/platform/config"
	"github.com/CyberGeo335/prak_nine/internal/repo"
)

func Run() {
	cfg := config.Load()
	db, err := repo.Open(cfg.DB_DSN)
	if err != nil {
		log.Fatal("db connect:", err)
	}

	if err := db.Exec("SET timezone TO 'UTC'").Error; err != nil { /* необязательно */
	}

	users := repo.NewUserRepo(db)
	if err := users.AutoMigrate(); err != nil {
		log.Fatal("migrate:", err)
	}

	auth := &handlers.AuthHandler{Users: users, BcryptCost: cfg.BcryptCost}

	r := chi.NewRouter()
	r.Post("/auth/register", auth.Register)
	r.Post("/auth/login", auth.Login)

	log.Println("listening on", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, r))
}
