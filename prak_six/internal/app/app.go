package app

import (
	"github.com/CyberGeo335/prak_six/internal/db"
	"github.com/CyberGeo335/prak_six/internal/httpapi"
	"github.com/CyberGeo335/prak_six/internal/models"
	"log"
	"net/http"
)

func Run() {
	d := db.Connect()

	if err := d.AutoMigrate(&models.User{}, &models.Note{}, &models.Tag{}); err != nil {
		log.Fatal("migrate:", err)
	}

	r := httpapi.BuildRouter(d)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
