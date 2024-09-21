package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gopher-market/services/auth-service/internal/models"
	"gopher-market/services/auth-service/internal/repository"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo *repository.AuthRepository
}

func NewAuthHandler(repo *repository.AuthRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Save to Auth Service database/Redis
	ctx := r.Context()
	if err := h.repo.SaveUser(ctx, &user); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	// Prepare request body for User Service
	userProfileBody := map[string]string{
		"username": user.Username,
		"email":    user.Email,
	}
	userProfileJSON, err := json.Marshal(userProfileBody)
	if err != nil {
		log.Printf("Failed to marshal user profile: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Call User Service to create detailed profile
	userProfileReq, err := http.NewRequestWithContext(ctx, "POST", "http://user-service:8081/users", bytes.NewReader(userProfileJSON))
	if err != nil {
		log.Printf("Failed to create request to User Service: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userProfileReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(userProfileReq)
	if err != nil {
		log.Printf("Failed to send request to User Service: %v", err)
		// Consider whether to rollback the auth creation or just log the error
		http.Error(w, "Failed to complete registration", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("User Service responded with status: %d", resp.StatusCode)
		http.Error(w, "Failed to complete registration", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginUser models.User
	if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUser(r.Context(), loginUser.Username)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	log.Printf("Retrieved user: %+v", user)
	log.Printf("Login attempt for user: %s", loginUser.Username)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password)); err != nil {
		log.Printf("Password comparison failed: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	log.Printf("user login successfully: %s", user.Username)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	if err := h.repo.SaveToken(r.Context(), user.Username, tokenString); err != nil {
		log.Printf("Error saving token: %v", err)
		http.Error(w, "Failed to save token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
