package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/CyberGeo335/prak_sixteen/internal/db"
	"github.com/CyberGeo335/prak_sixteen/internal/httpapi"
	"github.com/CyberGeo335/prak_sixteen/internal/repo"
	"github.com/CyberGeo335/prak_sixteen/internal/service"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is empty")
	}
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	dbx, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbx.Close()

	dbx.SetMaxOpenConns(10)
	dbx.SetMaxIdleConns(5)
	dbx.SetConnMaxLifetime(30 * time.Minute)

	if err := dbx.Ping(); err != nil {
		log.Fatal(err)
	}

	db.MustApplyMigrations(dbx)

	r := gin.Default()
	svc := service.Service{Notes: repo.NoteRepo{DB: dbx}}
	httpapi.Router{Svc: &svc}.Register(r)

	log.Printf("listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
