package models

import "time"

type Note struct {
	ID        int64     `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
