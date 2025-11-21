package remote

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/thebushidocollective/kasoku/internal/client"
	"github.com/thebushidocollective/kasoku/internal/config"
)

// CustomCache implements remote caching using a custom Kasoku server
type CustomCache struct {
	client *client.Client
}

// NewCustomCache creates a new custom cache backend
func NewCustomCache(cfg config.RemoteCacheConfig) (*CustomCache, error) {
	url := cfg.URL
	if url == "" && cfg.Endpoint != "" {
		// Fallback to legacy endpoint field
		url = cfg.Endpoint
	}

	if url == "" {
		return nil, fmt.Errorf("remote cache URL is required")
	}

	if cfg.Token == "" {
		return nil, fmt.Errorf("remote cache token is required")
	}

	c := client.NewClient(url, cfg.Token)

	return &CustomCache{
		client: c,
	}, nil
}

// Get retrieves an artifact from remote cache
func (c *CustomCache) Get(ctx context.Context, hash string) (io.ReadCloser, error) {
	buf := new(bytes.Buffer)
	if err := c.client.GetCache(ctx, hash, buf); err != nil {
		return nil, err
	}

	return io.NopCloser(buf), nil
}

// Put stores an artifact in remote cache
func (c *CustomCache) Put(ctx context.Context, hash string, command string, reader io.Reader) error {
	// Read all data to get size
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read artifact data: %w", err)
	}

	size := int64(len(data))

	// Upload to remote cache
	return c.client.PutCache(ctx, hash, command, bytes.NewReader(data), size)
}

// Exists checks if an artifact exists in remote cache
func (c *CustomCache) Exists(ctx context.Context, hash string) bool {
	// Try to get metadata - if it succeeds, the artifact exists
	// For now, we'll try a HEAD request-equivalent by attempting to get it
	buf := new(bytes.Buffer)
	err := c.client.GetCache(ctx, hash, buf)
	return err == nil
}

// Delete removes an artifact from remote cache
func (c *CustomCache) Delete(ctx context.Context, hash string) error {
	return c.client.DeleteCache(ctx, hash)
}
