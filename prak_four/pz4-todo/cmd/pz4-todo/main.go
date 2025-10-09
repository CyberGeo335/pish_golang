package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/CyberGeo335/pz4-todo/internal/task"
	myMW "github.com/CyberGeo335/pz4-todo/pkg/middleware"
)

// envOr возвращает значение переменной окружения или дефолт,
// если переменная не задана.
func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	// Хранилище на базе JSON-файла. Путь можно переопределить через DATA_FILE.
	repo := task.NewRepo(envOr("DATA_FILE", "tasks.json"))
	handler := task.NewHandler(repo)

	// Маршрутизатор и стандартные middleware.
	r := chi.NewRouter()
	r.Use(chimw.RequestID) // уникальный ID запроса
	r.Use(chimw.RealIP)    // реальный IP из заголовков
	r.Use(chimw.Recoverer) // защита от паник
	r.Use(myMW.Logger)     // наш логгер со статусом и длительностью
	r.Use(myMW.SimpleCORS) // CORS (для фронтов)

	// Простейший health-check.
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Версионирование API.
	r.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(v1 chi.Router) {
			v1.Mount("/tasks", handler.Routes()) // /api/v1/tasks
		})
	})

	// HTTP-сервер с таймаутами.
	srv := &http.Server{
		Addr:         ":" + envOr("PORT", "8080"),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер в отдельной горутине.
	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown по Ctrl+C / SIGTERM.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
