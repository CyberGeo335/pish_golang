package service

import (
	"context"
	"strings"

	"github.com/CyberGeo335/prak_sixteen/internal/models"
	"github.com/CyberGeo335/prak_sixteen/internal/repo"
)

type Service struct{ Notes repo.NotesRepository }

func (s Service) Create(ctx context.Context, n *models.Note) error {
	n.Title = strings.TrimSpace(n.Title)
	n.Content = strings.TrimSpace(n.Content)
	if n.Title == "" || n.Content == "" {
		return repo.ErrNotFound // можно завести отдельную ErrValidation
	}
	return s.Notes.Create(ctx, n)
}

func (s Service) Get(ctx context.Context, id int64) (models.Note, error) {
	return s.Notes.Get(ctx, id)
}
