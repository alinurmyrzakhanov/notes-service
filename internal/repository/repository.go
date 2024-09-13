package repository

import (
	"context"
	"notes-service/internal/models"
)

type NoteRepository interface {
	CreateNote(ctx context.Context, note *models.Note) error
	ListNotes(ctx context.Context, userID int64) ([]*models.Note, error)
	Close() error
}

type UserRepository interface {
	CreateUser(ctx context.Context, username, password string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	ValidateUser(ctx context.Context, username, password string) (*User, error)
}
