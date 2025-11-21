package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/thebushidocollective/kasoku/server/internal/db"
	"github.com/thebushidocollective/kasoku/server/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db        *db.DB
	jwtSecret []byte
}

func NewService(database *db.DB, jwtSecret string) *Service {
	return &Service{
		db:        database,
		jwtSecret: []byte(jwtSecret),
	}
}

// Register creates a new user account
func (s *Service) Register(email, name, password string) (*models.User, error) {
	// Check if user already exists
	existing, _ := s.db.GetUserByEmail(email)
	if existing != nil {
		return nil, fmt.Errorf("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:    email,
		Name:     name,
		Password: string(hashedPassword),
		Plan:     "free",
	}

	if err := s.db.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(email, password string) (string, *models.User, error) {
	// Get user
	user, err := s.db.GetUserByEmail(email)
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

// CreateAPIToken generates a new API token for CLI usage
func (s *Service) CreateAPIToken(userID uuid.UUID, name string, expiresAt *time.Time) (string, *models.Token, error) {
	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	rawToken := hex.EncodeToString(tokenBytes)

	// Hash token for storage
	hash := sha256.Sum256([]byte(rawToken))
	hashedToken := hex.EncodeToString(hash[:])

	// Create token record
	token := &models.Token{
		UserID:    userID,
		Name:      name,
		Token:     hashedToken,
		ExpiresAt: expiresAt,
	}

	if err := s.db.CreateToken(token); err != nil {
		return "", nil, fmt.Errorf("failed to create token: %w", err)
	}

	// Return raw token (only shown once)
	return rawToken, token, nil
}

// ValidateAPIToken checks if an API token is valid
func (s *Service) ValidateAPIToken(rawToken string) (*models.Token, *models.User, error) {
	// Hash the provided token
	hash := sha256.Sum256([]byte(rawToken))
	hashedToken := hex.EncodeToString(hash[:])

	// Get token from database
	token, err := s.db.GetTokenByValue(hashedToken)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid token")
	}

	// Check expiration
	if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now()) {
		return nil, nil, fmt.Errorf("token expired")
	}

	// Get user
	user, err := s.db.GetUserByID(token.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("user not found")
	}

	return token, user, nil
}

// RevokeAPIToken deletes an API token
func (s *Service) RevokeAPIToken(userID, tokenID uuid.UUID) error {
	// Verify token belongs to user
	tokens, err := s.db.ListTokens(userID)
	if err != nil {
		return err
	}

	found := false
	for _, t := range tokens {
		if t.ID == tokenID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("token not found")
	}

	return s.db.DeleteToken(tokenID)
}

// ListAPITokens returns all API tokens for a user
func (s *Service) ListAPITokens(userID uuid.UUID) ([]models.Token, error) {
	return s.db.ListTokens(userID)
}

// generateJWT creates a JWT token for a user
func (s *Service) generateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"plan":    user.Plan,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	if user.TeamID != nil {
		claims["team_id"] = user.TeamID.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// Middleware returns a middleware that validates JWT tokens or API tokens
func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if it's a Bearer token (JWT) or API token
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Try JWT first
			if user, teamID, err := s.validateJWT(tokenString); err == nil {
				ctx := context.WithValue(r.Context(), "user_id", user.ID)
				if teamID != nil {
					ctx = context.WithValue(ctx, "team_id", *teamID)
				}
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Try API token
			if _, user, err := s.ValidateAPIToken(tokenString); err == nil {
				ctx := context.WithValue(r.Context(), "user_id", user.ID)
				if user.TeamID != nil {
					ctx = context.WithValue(ctx, "team_id", *user.TeamID)
				}
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})
}

// validateJWT verifies a JWT token and returns the user
func (s *Service) validateJWT(tokenString string) (*models.User, *uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, fmt.Errorf("invalid claims")
	}

	// Get user ID
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("invalid user_id")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid user_id format")
	}

	// Get user from database
	user, err := s.db.GetUserByID(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("user not found")
	}

	// Get team ID if present
	var teamID *uuid.UUID
	if teamIDStr, ok := claims["team_id"].(string); ok {
		tid, err := uuid.Parse(teamIDStr)
		if err == nil {
			teamID = &tid
		}
	}

	return user, teamID, nil
}
