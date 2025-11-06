package db

import (
	"github.com/joho/godotenv"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func BuildPostgresURL() (string, error) {
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, pass),
		Host:   net.JoinHostPort(host, port),
		Path:   "/" + name,
	}

	// Опциональные параметры — добавляем только если заданы
	q := url.Values{}
	if v := os.Getenv("DB_SSLMODE"); v != "" {
		q.Set("sslmode", v) // напр. disable / require / verify-full
	}
	if v := os.Getenv("DB_SEARCH_PATH"); v != "" {
		q.Set("search_path", v)
	}
	if len(q) > 0 {
		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}

func Connect() *gorm.DB {

	_ = godotenv.Load()

	dsn, err := BuildPostgresURL()
	if err != nil {
		log.Fatal("build DSN:", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal("connect db:", err)
	}

	// Пул соединений
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(2)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}
