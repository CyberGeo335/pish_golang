package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Port      string
	JWTSecret []byte        // остался для совместимости, RS256 здесь его не использует
	JWTTTL    time.Duration // TTL access-токена
}

func Load() Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8083"
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// для RS256 секрет не обязателен, но оставим проверку,
		secret = "dev-secret"
	}

	ttl := os.Getenv("JWT_TTL")
	if ttl == "" {
		ttl = "15m" // по заданию: access TTL 15 минут
	}
	dur, err := time.ParseDuration(ttl)
	if err != nil {
		log.Fatal("bad JWT_TTL")
	}

	return Config{
		Port:      ":" + port,
		JWTSecret: []byte(secret),
		JWTTTL:    dur,
	}
}
