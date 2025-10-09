package main

import (
	"example.com/pz3_http/internal/api"
	"example.com/pz3_http/internal/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	store := storage.NewMemoryStore()
	h := api.NewHandlers(store)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		api.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /tasks", h.ListTasks)
	mux.HandleFunc("POST /tasks", h.CreateTask)

	mux.HandleFunc("GET /tasks/{id}", h.GetTask)
	mux.HandleFunc("PATCH /tasks/{id}", h.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", h.DeleteTask)

	handler := api.WithCORS(api.WithLogging(mux))
	addr := getAddr()
	log.Println("listening on", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func getAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
