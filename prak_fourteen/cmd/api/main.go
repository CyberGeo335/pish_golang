package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/CyberGeo335/prak_fourteen/internal/config"
	"github.com/CyberGeo335/prak_fourteen/internal/storage/postgres"
	rediscache "github.com/CyberGeo335/prak_fourteen/internal/storage/redis"
	httptransport "github.com/CyberGeo335/prak_fourteen/internal/transport/http"
)

func main() {
	_ = godotenv.Load()

	cfg := config.FromEnv()

	// Postgres pool
	pgxCfg, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	pgxCfg.MaxConns = 20
	pgxCfg.MinConns = 5
	pgxCfg.MaxConnLifetime = time.Hour
	// Statement cache (pgx автоматически "готовит" и кеширует планы на соединение)
	pgxCfg.ConnConfig.StatementCacheCapacity = 256

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := postgres.NewRepo(pool)

	// Redis cache
	cache, err := rediscache.New(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, cfg.CacheTTL)
	if err != nil {
		log.Fatal(err)
	}
	defer cache.Close()

	srv := httptransport.NewServer(repo, cache)

	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, srv.Router()); err != nil {
		log.Fatal(err)
	}
}
