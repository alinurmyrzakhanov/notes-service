package repository

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Пароль не должен сериализоваться в JSON
}

type SQLUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) CreateUser(ctx context.Context, username, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var user User
	err = r.db.QueryRowContext(ctx,
		"INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, username",
		username, hashedPassword).Scan(&user.ID, &user.Username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *SQLUserRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, username, password FROM users WHERE username = $1",
		username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Пользователь не найден
		}
		return nil, err
	}
	return &user, nil
}

func (r *SQLUserRepository) ValidateUser(ctx context.Context, username, password string) (*User, error) {
	user, err := r.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil // Пользователь не найден
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err // Неверный пароль
	}

	return user, nil
}
