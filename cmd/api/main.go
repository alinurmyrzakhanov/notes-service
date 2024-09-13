package main

import (
	"log"
	"net/http"
	"notes-service/internal/auth"
	"notes-service/internal/config"
	"notes-service/internal/handlers"
	"notes-service/internal/repository"
	"notes-service/internal/spellcheck"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	postgresRepo, err := repository.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer postgresRepo.Close()

	userRepo := repository.NewUserRepository(postgresRepo.GetDB())
	spellchecker := spellcheck.NewYandexSpellchecker(cfg.YandexSpellcheckerURL)
	authService := auth.NewAuthService(userRepo, cfg.JWTSecret)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	noteHandler := handlers.NewNoteHandler(postgresRepo, spellchecker, authService)

	r.Post("/register", authService.Register)
	r.Post("/login", authService.Login)

	r.Group(func(r chi.Router) {
		r.Use(authService.Authenticate)
		r.Post("/notes", noteHandler.CreateNote)
		r.Get("/notes", noteHandler.ListNotes)
	})

	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
