package repo

import (
	"context"
	"errors"

	"github.com/CyberGeo335/prak_nine/internal/core"
	"github.com/jackc/pgconn" // ← вот это добавить в imports
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

	// 1. Сначала пробуем нормальный путь GORM
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrEmailTaken
	}

	// 2. Явно распарсим ошибку postgres по коду SQLSTATE
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrEmailTaken
	}

	// 3. Всё остальное трушная db-ошибка
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
