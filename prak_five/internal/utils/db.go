package utils

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"time"
)

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// настройки пула — достаточно для локалки
	db.SetMaxOpenConns(2)                   // максимум активных соединений
	db.SetMaxIdleConns(2)                   // соединений в простое
	db.SetConnMaxLifetime(30 * time.Minute) // будет жить 30 минут

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	log.Println("Connected to PostgreSQL")
	return db, nil

}
