package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string          `json:"token"`
	User  json.RawMessage `json:"user"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	ID      string `json:"id"`
	Name    string `json:"name"`
}

type CacheEntry struct {
	ID        string    `json:"id"`
	Hash      string    `json:"hash"`
	Command   string    `json:"command"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

type ListCacheResponse struct {
	Entries []CacheEntry `json:"entries"`
	Limit   int          `json:"limit"`
	Offset  int          `json:"offset"`
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Login authenticates with the server and returns a token
func (c *Client) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	req := LoginRequest{
		Email:    email,
		Password: password,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, "POST", "/auth/login", bytes.NewReader(body), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed: %s", resp.Status)
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, err
	}

	return &loginResp, nil
}

// Register creates a new account
func (c *Client) Register(ctx context.Context, email, name, password string) error {
	req := RegisterRequest{
		Email:    email,
		Name:     name,
		Password: password,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := c.doRequest(ctx, "POST", "/auth/register", bytes.NewReader(body), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("registration failed: %s", resp.Status)
	}

	return nil
}

// CreateToken creates a new API token
func (c *Client) CreateToken(ctx context.Context, name string) (*TokenResponse, error) {
	req := map[string]string{"name": name}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, "POST", "/auth/tokens", bytes.NewReader(body), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create token: %s", resp.Status)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// PutCache uploads a cache artifact
func (c *Client) PutCache(ctx context.Context, hash string, command string, reader io.Reader, size int64) error {
	url := fmt.Sprintf("/cache/%s?command=%s", hash, command)

	headers := map[string]string{
		"Content-Length": fmt.Sprintf("%d", size),
	}

	resp, err := c.doRequest(ctx, "PUT", url, reader, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload cache: %s - %s", resp.Status, string(body))
	}

	return nil
}

// GetCache downloads a cache artifact
func (c *Client) GetCache(ctx context.Context, hash string, writer io.Writer) error {
	url := fmt.Sprintf("/cache/%s", hash)

	resp, err := c.doRequest(ctx, "GET", url, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("cache miss")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download cache: %s", resp.Status)
	}

	if _, err := io.Copy(writer, resp.Body); err != nil {
		return fmt.Errorf("failed to read cache data: %w", err)
	}

	return nil
}

// ListCache lists cache entries
func (c *Client) ListCache(ctx context.Context, limit, offset int) (*ListCacheResponse, error) {
	url := fmt.Sprintf("/cache?limit=%d&offset=%d", limit, offset)

	resp, err := c.doRequest(ctx, "GET", url, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list cache: %s", resp.Status)
	}

	var listResp ListCacheResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, err
	}

	return &listResp, nil
}

// DeleteCache deletes a cache entry
func (c *Client) DeleteCache(ctx context.Context, hash string) error {
	url := fmt.Sprintf("/cache/%s", hash)

	resp, err := c.doRequest(ctx, "DELETE", url, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete cache: %s", resp.Status)
	}

	return nil
}

// GetAnalytics retrieves analytics summary
func (c *Client) GetAnalytics(ctx context.Context, since time.Time) (map[string]interface{}, error) {
	url := fmt.Sprintf("/analytics?since=%s", since.Format(time.RFC3339))

	resp, err := c.doRequest(ctx, "GET", url, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get analytics: %s", resp.Status)
	}

	var analytics map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&analytics); err != nil {
		return nil, err
	}

	return analytics, nil
}

// doRequest is a helper method to make HTTP requests
func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader, headers map[string]string) (*http.Response, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")

	// Set auth token if available
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.httpClient.Do(req)
}
