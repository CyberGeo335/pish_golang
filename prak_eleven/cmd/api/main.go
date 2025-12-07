package main

import (
	"github.com/CyberGeo335/prak_eleven/internal/http"
	"github.com/CyberGeo335/prak_eleven/internal/http/handlers"
	"github.com/CyberGeo335/prak_eleven/internal/repo"
	"log"
	"net/http"
)

func main() {
	repo := repo.NewNoteRepoMem()
	h := &handlers.Handler{Repo: repo}
	r := httpx.NewRouter(h)

	log.Println("Server started at :8085")
	log.Fatal(http.ListenAndServe(":8085", r))
}
