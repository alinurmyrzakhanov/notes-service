package repository

import (
	"context"
	"database/sql"
	"notes-service/internal/models"

	_ "github.com/lib/pq"
)

// PostgresRepository реализует методы для работы с PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository создает новый экземпляр PostgresRepository
func NewPostgresRepository(dbURL string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresRepository{db: db}, nil
}

// GetDB возвращает экземпляр *sql.DB
func (r *PostgresRepository) GetDB() *sql.DB {
	return r.db
}

// Close закрывает соединение с базой данных
func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

// CreateNote создает новую заметку в базе данных
func (r *PostgresRepository) CreateNote(ctx context.Context, note *models.Note) error {
	query := `
		INSERT INTO notes (user_id, title, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		note.UserID, note.Title, note.Content, note.CreatedAt, note.UpdatedAt).
		Scan(&note.ID)

	return err
}

// ListNotes возвращает список заметок пользователя
func (r *PostgresRepository) ListNotes(ctx context.Context, userID int64) ([]*models.Note, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM notes
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(
			&note.ID, &note.UserID, &note.Title, &note.Content,
			&note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}
