package handlers

import (
	"encoding/json"
	"net/http"
	"notes-service/internal/auth"
	"notes-service/internal/models"
	"notes-service/internal/repository"
	"notes-service/internal/spellcheck"
	"time"
)

// NoteHandler обрабатывает запросы, связанные с заметками
type NoteHandler struct {
	repo         repository.NoteRepository
	spellchecker spellcheck.Spellchecker // Change this to an interface
	authService  auth.AuthService        // Change this to an interface
}

// NewNoteHandler создает новый экземпляр NoteHandler
func NewNoteHandler(repo repository.NoteRepository, spellchecker spellcheck.Spellchecker, authService auth.AuthService) *NoteHandler {
	return &NoteHandler{
		repo:         repo,
		spellchecker: spellchecker,
		authService:  authService,
	}
}

// CreateNote обрабатывает создание новой заметки
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка орфографии
	correctedContent, err := h.spellchecker.CheckSpelling(note.Content)
	if err != nil {
		http.Error(w, "Failed to check spelling", http.StatusInternalServerError)
		return
	}
	note.Content = correctedContent

	// Получение ID пользователя из контекста (установленного middleware аутентификации)
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "Unauthorize", http.StatusUnauthorized)
		return
	}
	note.UserID = userID

	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	if err := h.repo.CreateNote(r.Context(), &note); err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

// ListNotes обрабатывает запрос на получение списка заметок пользователя
func (h *NoteHandler) ListNotes(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notes, err := h.repo.ListNotes(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch notes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}
