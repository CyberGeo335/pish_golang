package main

import (
	"log"
	"net/http"

	router "github.com/CyberGeo335/prak_ten/internal/http"
	"github.com/CyberGeo335/prak_ten/internal/platform/config"
)

func main() {
	cfg := config.Load()
	mux, err := router.Build(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening on", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
