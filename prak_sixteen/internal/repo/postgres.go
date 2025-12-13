package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/CyberGeo335/prak_sixteen/internal/models"
)

var ErrNotFound = errors.New("not found")

type NotesRepository interface {
	Create(ctx context.Context, n *models.Note) error
	Get(ctx context.Context, id int64) (models.Note, error)
	Update(ctx context.Context, n *models.Note) error
	Delete(ctx context.Context, id int64) error
}

type NoteRepo struct{ DB *sql.DB }

func (r NoteRepo) Create(ctx context.Context, n *models.Note) error {
	return r.DB.QueryRowContext(ctx,
		`INSERT INTO notes(title, content)
		 VALUES($1,$2)
		 RETURNING id, created_at, updated_at`,
		n.Title, n.Content,
	).Scan(&n.ID, &n.CreatedAt, &n.UpdatedAt)
}

func (r NoteRepo) Get(ctx context.Context, id int64) (models.Note, error) {
	var n models.Note
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, title, content, created_at, updated_at
		 FROM notes
		 WHERE id=$1`,
		id,
	).Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)

	if err == sql.ErrNoRows {
		return models.Note{}, ErrNotFound
	}
	return n, err
}

func (r NoteRepo) Update(ctx context.Context, n *models.Note) error {
	err := r.DB.QueryRowContext(ctx,
		`UPDATE notes
		 SET title=$2, content=$3
		 WHERE id=$1
		 RETURNING created_at, updated_at`,
		n.ID, n.Title, n.Content,
	).Scan(&n.CreatedAt, &n.UpdatedAt)

	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	return err
}

func (r NoteRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.DB.ExecContext(ctx, `DELETE FROM notes WHERE id=$1`, id)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return ErrNotFound
	}
	return nil
}
