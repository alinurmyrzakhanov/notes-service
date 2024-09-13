package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"notes-service/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateNote(ctx context.Context, note *models.Note) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MockRepository) ListNotes(ctx context.Context, userID int64) ([]*models.Note, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.Note), args.Error(1)
}

func (m *MockRepository) Close() error {
	return nil
}

type MockSpellchecker struct {
	mock.Mock
}

func (m *MockSpellchecker) CheckSpelling(text string) (string, error) {
	args := m.Called(text)
	return args.String(0), args.Error(1)
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "user_id", int64(1))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TestCreateNote(t *testing.T) {
	mockRepo := new(MockRepository)
	mockSpellchecker := new(MockSpellchecker)
	mockAuthService := new(MockAuthService)

	handler := NewNoteHandler(mockRepo, mockSpellchecker, mockAuthService)

	mockSpellchecker.On("CheckSpelling", "This is a test note.").Return("This is a test note.", nil)
	mockRepo.On("CreateNote", mock.Anything, mock.AnythingOfType("*models.Note")).Return(nil)

	reqBody := bytes.NewBufferString(`{"title":"Test Note","content":"This is a test note."}`)
	req, _ := http.NewRequest("POST", "/notes", reqBody)
	rr := httptest.NewRecorder()

	// Оборачиваем обработчик в middleware аутентификации
	authenticatedHandler := mockAuthService.Authenticate(http.HandlerFunc(handler.CreateNote))

	authenticatedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response models.Note
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, "Test Note", response.Title)
	assert.Equal(t, "This is a test note.", response.Content)
}

func TestListNotes(t *testing.T) {
	mockRepo := new(MockRepository)
	mockSpellchecker := new(MockSpellchecker)
	mockAuthService := new(MockAuthService)

	handler := NewNoteHandler(mockRepo, mockSpellchecker, mockAuthService)

	mockNotes := []*models.Note{
		{ID: 1, Title: "Note 1", Content: "Content 1"},
		{ID: 2, Title: "Note 2", Content: "Content 2"},
	}

	mockRepo.On("ListNotes", mock.Anything, int64(1)).Return(mockNotes, nil)

	req, _ := http.NewRequest("GET", "/notes", nil)
	rr := httptest.NewRecorder()

	// Оборачиваем обработчик в middleware аутентификации
	authenticatedHandler := mockAuthService.Authenticate(http.HandlerFunc(handler.ListNotes))

	authenticatedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []*models.Note
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Len(t, response, 2)
	assert.Equal(t, "Note 1", response[0].Title)
	assert.Equal(t, "Note 2", response[1].Title)
}
