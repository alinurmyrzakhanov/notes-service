package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("testuser", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).
			AddRow(1, "testuser"))

	user, err := repo.CreateUser(context.Background(), "testuser", "password")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "testuser", user.Username)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "username", "password"}).
		AddRow(1, "testuser", "hashedpassword")

	mock.ExpectQuery("SELECT (.+) FROM users WHERE username = ?").
		WithArgs("testuser").
		WillReturnRows(rows)

	user, err := repo.GetUserByUsername(context.Background(), "testuser")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "testuser", user.Username)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestValidateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	// Хешированный пароль "password"
	hashedPassword := "$2a$10$1234567890123456789012abcdefghijklmnopqrstuvwxyzABCDEF"

	rows := sqlmock.NewRows([]string{"id", "username", "password"}).
		AddRow(1, "testuser", hashedPassword)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE username = ?").
		WithArgs("testuser").
		WillReturnRows(rows)

	user, err := repo.ValidateUser(context.Background(), "testuser", "password")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "testuser", user.Username)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
