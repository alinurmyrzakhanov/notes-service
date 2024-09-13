package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"notes-service/internal/repository"

	"github.com/dgrijalva/jwt-go"
)

type AuthService interface {
	Authenticate(next http.Handler) http.Handler
}

type AuthServiceImpl struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *AuthServiceImpl) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.userRepo.CreateUser(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
	})
}

func (s *AuthServiceImpl) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.userRepo.ValidateUser(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	token, err := s.generateToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *AuthServiceImpl) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return s.jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", int64(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *AuthServiceImpl) generateToken(userID int64, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(s.jwtSecret)
}
