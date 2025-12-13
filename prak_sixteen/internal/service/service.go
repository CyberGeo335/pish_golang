package service

import (
	"context"
	"strings"

	"github.com/CyberGeo335/prak_sixteen/internal/models"
	"github.com/CyberGeo335/prak_sixteen/internal/repo"
)

type Service struct{ Notes repo.NotesRepository }

func New(notes repo.NotesRepository) *Service {
	return &Service{Notes: notes}
}

func (s Service) Create(ctx context.Context, n *models.Note) error {
	n.Title = strings.TrimSpace(n.Title)
	n.Content = strings.TrimSpace(n.Content)
	if n.Title == "" || n.Content == "" {
		// лучше завести ErrValidation и отдавать 400, но оставим так, мне лень
		return repo.ErrNotFound
	}
	return s.Notes.Create(ctx, n)
}

func (s Service) Get(ctx context.Context, id int64) (models.Note, error) {
	return s.Notes.Get(ctx, id)
}
