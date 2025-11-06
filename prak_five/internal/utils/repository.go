package utils

import (
	"context"
	"database/sql"
	"errors"
	"github.com/CyberGeo335/prak_five/internal/config"
)

type Repo struct {
	DB *sql.DB
}

var ErrNotFound = errors.New("task not found")

func NewRepo(db *sql.DB) *Repo { return &Repo{DB: db} }

// CreateTask — параметризованный INSERT с возвратом id
// просто такой стиль записи, кринж но ок.
func (r *Repo) CreateTask(ctx context.Context, title string) (int, error) {
	var id int
	const q = `INSERT INTO tasks (title) VALUES ($1) RETURNING id;`
	err := r.DB.QueryRowContext(ctx, q, title).Scan(&id)
	return id, err
}

// ListTasks — базовый SELECT всех задач
func (r *Repo) ListTasks(ctx context.Context) ([]config.Task, error) {
	const q = `SELECT id, title, done, created_at FROM tasks ORDER BY id;`
	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []config.Task
	for rows.Next() {
		var t config.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// ListDone - список выполниненных/невыполненных
func (r *Repo) ListDone(ctx context.Context, done bool) ([]config.Task, error) {
	const q = `SELECT 
    			id, title, done, created_at 
				FROM 
				    tasks 
				WHERE 
				    done=true;`
	rows, err := r.DB.QueryContext(ctx, q, done)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []config.Task
	for rows.Next() {
		var t config.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// FindByID - подробности по задаче с конкретным id
func (r *Repo) FindByID(ctx context.Context, id int) (*config.Task, error) {
	const q = `SELECT id, title, done, created_at
             FROM tasks
             WHERE id = $1;`
	var t config.Task
	err := r.DB.QueryRowContext(ctx, q, id).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

// CreateMany - массовая вставка через транзакцию
func (r *Repo) CreateMany(ctx context.Context, titles []string) error {
	if len(titles) == 0 {
		return nil
	}

	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO tasks (title) VALUES ($1);`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, title := range titles {
		if _, err := stmt.ExecContext(ctx, title); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
