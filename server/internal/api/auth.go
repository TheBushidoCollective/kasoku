package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/thebushidocollective/kasoku/server/internal/auth"
)

type AuthHandler struct {
	auth *auth.Service
}

func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{
		auth: authService,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateTokenRequest struct {
	Name      string     `json:"name"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	if len(req.Password) < 8 {
		respondError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	// Create user
	user, err := h.auth.Register(req.Email, req.Name, req.Password)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "user created successfully",
		"user":    user,
	})
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Authenticate
	token, user, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// CreateToken handles POST /auth/tokens
func (h *AuthHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "name is required")
		return
	}

	// Create API token
	rawToken, token, err := h.auth.CreateAPIToken(userID, req.Name, req.ExpiresAt)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "token created successfully",
		"token":   rawToken, // Only shown once!
		"id":      token.ID,
		"name":    token.Name,
	})
}

// ListTokens handles GET /auth/tokens
func (h *AuthHandler) ListTokens(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	tokens, err := h.auth.ListAPITokens(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"tokens": tokens,
	})
}

// RevokeToken handles DELETE /auth/tokens/:id
func (h *AuthHandler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get token ID from URL
	tokenIDStr := chi.URLParam(r, "id")
	tokenID, err := uuid.Parse(tokenIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid token id")
		return
	}

	// Revoke token
	if err := h.auth.RevokeAPIToken(userID, tokenID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "token revoked",
	})
}

// Me handles GET /auth/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get user details (this is set in the auth middleware)
	// For now, just return the user_id from context
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"user_id": userID,
	})
}
