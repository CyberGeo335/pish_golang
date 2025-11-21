package main

import (
	"log"
	"net/http"

	router "github.com/CyberGeo335/prak_ten/internal/http"
	"github.com/CyberGeo335/prak_ten/internal/platform/config"
)

func main() {
	cfg := config.Load()

	mux := router.Build(cfg)

	log.Println("listening on", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}
