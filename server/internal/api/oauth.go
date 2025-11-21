package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/thebushidocollective/kasoku/server/internal/auth"
	"golang.org/x/oauth2"
)

type OAuthHandler struct {
	auth        *auth.Service
	oauthConfig *auth.OAuthConfig
	states      map[string]bool // Simple state management (use Redis in production)
}

func NewOAuthHandler(authService *auth.Service, oauthConfig *auth.OAuthConfig) *OAuthHandler {
	return &OAuthHandler{
		auth:        authService,
		oauthConfig: oauthConfig,
		states:      make(map[string]bool),
	}
}

// GitHubLogin handles GET /auth/github
func (h *OAuthHandler) GitHubLogin(w http.ResponseWriter, r *http.Request) {
	if h.oauthConfig.GitHub == nil || h.oauthConfig.GitHub.ClientID == "" {
		respondError(w, http.StatusServiceUnavailable, "GitHub OAuth not configured")
		return
	}

	state := generateState()
	h.states[state] = true

	url := h.oauthConfig.GitHub.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GitLabLogin handles GET /auth/gitlab
func (h *OAuthHandler) GitLabLogin(w http.ResponseWriter, r *http.Request) {
	if h.oauthConfig.GitLab == nil || h.oauthConfig.GitLab.ClientID == "" {
		respondError(w, http.StatusServiceUnavailable, "GitLab OAuth not configured")
		return
	}

	state := generateState()
	h.states[state] = true

	url := h.oauthConfig.GitLab.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GitHubCallback handles GET /auth/callback/github
func (h *OAuthHandler) GitHubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	// Validate state
	if !h.states[state] {
		respondError(w, http.StatusBadRequest, "invalid state")
		return
	}
	delete(h.states, state)

	if code == "" {
		respondError(w, http.StatusBadRequest, "code not provided")
		return
	}

	// Exchange code for token and create/login user
	token, user, err := h.auth.HandleOAuthLogin(auth.ProviderGitHub, code, h.oauthConfig.GitHub)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("OAuth login failed: %s", err.Error()))
		return
	}

	// Redirect to frontend with token
	redirectURL := fmt.Sprintf("http://localhost:3000/auth/callback?token=%s&user_id=%s", token, user.ID)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// GitLabCallback handles GET /auth/callback/gitlab
func (h *OAuthHandler) GitLabCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	// Validate state
	if !h.states[state] {
		respondError(w, http.StatusBadRequest, "invalid state")
		return
	}
	delete(h.states, state)

	if code == "" {
		respondError(w, http.StatusBadRequest, "code not provided")
		return
	}

	// Exchange code for token and create/login user
	token, user, err := h.auth.HandleOAuthLogin(auth.ProviderGitLab, code, h.oauthConfig.GitLab)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("OAuth login failed: %s", err.Error()))
		return
	}

	// Redirect to frontend with token
	redirectURL := fmt.Sprintf("http://localhost:3000/auth/callback?token=%s&user_id=%s", token, user.ID)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// generateState generates a random state string for CSRF protection
func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
