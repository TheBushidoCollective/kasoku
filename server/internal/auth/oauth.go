package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/thebushidocollective/kasoku/server/internal/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type OAuthProvider string

const (
	ProviderGitHub OAuthProvider = "github"
	ProviderGitLab OAuthProvider = "gitlab"
)

type OAuthConfig struct {
	GitHub *oauth2.Config
	GitLab *oauth2.Config
}

// NewOAuthConfig creates OAuth configurations for providers
func NewOAuthConfig(baseURL string) *OAuthConfig {
	return &OAuthConfig{
		GitHub: &oauth2.Config{
			ClientID:     "", // Set via env var GITHUB_CLIENT_ID
			ClientSecret: "", // Set via env var GITHUB_CLIENT_SECRET
			RedirectURL:  baseURL + "/auth/callback/github",
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
		GitLab: &oauth2.Config{
			ClientID:     "", // Set via env var GITLAB_CLIENT_ID
			ClientSecret: "", // Set via env var GITLAB_CLIENT_SECRET
			RedirectURL:  baseURL + "/auth/callback/gitlab",
			Scopes:       []string{"read_user"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://gitlab.com/oauth/authorize",
				TokenURL: "https://gitlab.com/oauth/token",
			},
		},
	}
}

// GitHubUser represents GitHub user API response
type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// GitLabUser represents GitLab user API response
type GitLabUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// HandleOAuthLogin handles OAuth login for GitHub/GitLab
func (s *Service) HandleOAuthLogin(provider OAuthProvider, code string, oauthConfig *oauth2.Config) (string, *models.User, error) {
	ctx := context.Background()

	// Exchange code for token
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return "", nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Get user info from provider
	var email, name, providerID, avatarURL string
	switch provider {
	case ProviderGitHub:
		ghUser, err := s.fetchGitHubUser(ctx, token)
		if err != nil {
			return "", nil, fmt.Errorf("failed to fetch GitHub user: %w", err)
		}
		email = ghUser.Email
		name = ghUser.Name
		if name == "" {
			name = ghUser.Login
		}
		providerID = fmt.Sprintf("%d", ghUser.ID)
		avatarURL = ghUser.AvatarURL

	case ProviderGitLab:
		glUser, err := s.fetchGitLabUser(ctx, token)
		if err != nil {
			return "", nil, fmt.Errorf("failed to fetch GitLab user: %w", err)
		}
		email = glUser.Email
		name = glUser.Name
		if name == "" {
			name = glUser.Username
		}
		providerID = fmt.Sprintf("%d", glUser.ID)
		avatarURL = glUser.AvatarURL

	default:
		return "", nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	// Check if user exists by email
	user, err := s.db.GetUserByEmail(email)
	if err != nil {
		// User doesn't exist, create new one
		user = &models.User{
			Email:      email,
			Name:       name,
			Provider:   string(provider),
			ProviderID: providerID,
			AvatarURL:  &avatarURL,
			Plan:       "free",
		}

		if err := s.db.CreateUser(user); err != nil {
			return "", nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		// User exists, update OAuth info if needed
		if user.Provider == "" {
			user.Provider = string(provider)
			user.ProviderID = providerID
			user.AvatarURL = &avatarURL
			if err := s.db.UpdateUser(user); err != nil {
				return "", nil, fmt.Errorf("failed to update user: %w", err)
			}
		}
	}

	// Generate JWT token
	jwtToken, err := s.generateJWT(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	return jwtToken, user, nil
}

// fetchGitHubUser retrieves user info from GitHub API
func (s *Service) fetchGitHubUser(ctx context.Context, token *oauth2.Token) (*GitHubUser, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s", string(body))
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	// If email is not public, fetch from emails endpoint
	if user.Email == "" {
		emails, err := s.fetchGitHubEmails(ctx, token)
		if err == nil && len(emails) > 0 {
			user.Email = emails[0]
		}
	}

	return &user, nil
}

// fetchGitHubEmails retrieves user emails from GitHub API
func (s *Service) fetchGitHubEmails(ctx context.Context, token *oauth2.Token) ([]string, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	var emailsResp []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emailsResp); err != nil {
		return nil, err
	}

	var emails []string
	for _, e := range emailsResp {
		if e.Primary && e.Verified {
			emails = append(emails, e.Email)
		}
	}

	return emails, nil
}

// fetchGitLabUser retrieves user info from GitLab API
func (s *Service) fetchGitLabUser(ctx context.Context, token *oauth2.Token) (*GitLabUser, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", "https://gitlab.com/api/v4/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitLab API error: %s", string(body))
	}

	var user GitLabUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
