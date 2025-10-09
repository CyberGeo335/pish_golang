package app

import (
	"github.com/CyberGeo335/myapp/internal/app/handlers"
	"github.com/CyberGeo335/myapp/utils"
	"net/http"
	"os"
)

func Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Root)
	mux.HandleFunc("/ping", handlers.Ping)
	mux.HandleFunc("/fail", handlers.Fail)

	handler := withRequestID(mux)
	addr := getAddr()

	utils.LogInfo("Server is starting on " + addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		utils.LogError("server error: " + err.Error())
	}

}

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-Id")
		if id == "" {
			id = utils.NewID16()
		}
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r)
	})
}

func getAddr() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
