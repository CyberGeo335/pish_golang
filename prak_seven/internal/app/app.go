package app

import (
	"github.com/CyberGeo335/prak_seven/internal/cache"
	"github.com/CyberGeo335/prak_seven/internal/handlers"
	"log"
	"net/http"
)

var globalCache *cache.Cache

func Run() {
	c := cache.New("127.0.0.1:6379")

	mux := http.NewServeMux()

	//mux.HandleFunc("/set", handlers.SetHandler) -> Не работает плак плак

	// "/set"
	mux.HandleFunc("/set", func(writer http.ResponseWriter, request *http.Request) {
		handlers.SetHandler(c, writer, request)
	})

	// "/get"
	mux.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		handlers.GetHandler(c, writer, request)
	})

	mux.HandleFunc("/ttl", func(writer http.ResponseWriter, request *http.Request) {
		handlers.TtlHandler(c, writer, request)
	})

	log.Println("Prak seven server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
