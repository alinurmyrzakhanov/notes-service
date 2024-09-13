package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"notes-service/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, username, password string) (*repository.User, error) {
	args := m.Called(ctx, username, password)
	return args.Get(0).(*repository.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*repository.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*repository.User), args.Error(1)
}

func (m *MockUserRepository) ValidateUser(ctx context.Context, username, password string) (*repository.User, error) {
	args := m.Called(ctx, username, password)
	return args.Get(0).(*repository.User), args.Error(1)
}

func TestRegister(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, "secret")

	mockRepo.On("CreateUser", mock.Anything, "testuser", "password").Return(&repository.User{
		ID:       1,
		Username: "testuser",
	}, nil)

	reqBody := bytes.NewBufferString(`{"username":"testuser","password":"password"}`)
	req, _ := http.NewRequest("POST", "/register", reqBody)
	rr := httptest.NewRecorder()

	authService.Register(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, "testuser", response["username"])
}

func TestLogin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, "secret")

	mockRepo.On("ValidateUser", mock.Anything, "testuser", "password").Return(&repository.User{
		ID:       1,
		Username: "testuser",
	}, nil)

	reqBody := bytes.NewBufferString(`{"username":"testuser","password":"password"}`)
	req, _ := http.NewRequest("POST", "/login", reqBody)
	rr := httptest.NewRecorder()

	authService.Login(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.NotEmpty(t, response["token"])
}
