package repo

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/CyberGeo335/prak_nine/internal/core"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailTaken = errors.New("email already in use")

type UserRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) AutoMigrate() error {
	return r.db.AutoMigrate(&core.User{})
}

func (r *UserRepo) Create(ctx context.Context, u *core.User) error {
	err := r.db.WithContext(ctx).Create(u).Error
	if err == nil {
		return nil
	}

	// Логируем "сырую" ошибку и её тип
	log.Printf("Create user raw error (%T): %v\n", err, err)

	// 1) Стандартный путь gorm для дубликата
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		log.Println("mapped duplicate via gorm.ErrDuplicatedKey")
		return ErrEmailTaken
	}

	// 2) Прямо парсим postgres-ошибку
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		log.Printf("pg error code: %s, message: %s\n", pgErr.Code, pgErr.Message)
		if pgErr.Code == "23505" {
			log.Println("mapped duplicate via pgErr 23505")
			return ErrEmailTaken
		}
	}

	// 3) Жёсткий, но надёжный fallback по тексту
	msg := err.Error()
	if strings.Contains(msg, "SQLSTATE 23505") ||
		strings.Contains(msg, "duplicate key value violates unique constraint") {
		log.Println("mapped duplicate via string contains")
		return ErrEmailTaken
	}

	// Всё остальное — реальная DB-ошибка
	return err
}

func (r *UserRepo) ByEmail(ctx context.Context, email string) (core.User, error) {
	var u core.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return core.User{}, ErrUserNotFound
	}
	return u, err
}
