package custom

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/thebushidocollective/kasoku/internal/client"
)

// Backend implements storage.Backend for custom Kasoku server
type Backend struct {
	client *client.Client
}

// New creates a new custom storage backend
func New(url, token string) (*Backend, error) {
	if url == "" {
		return nil, fmt.Errorf("custom backend URL is required")
	}

	if token == "" {
		return nil, fmt.Errorf("custom backend token is required")
	}

	c := client.NewClient(url, token)

	return &Backend{
		client: c,
	}, nil
}

// Get retrieves data from remote cache
// The key format is expected to be: {hash}/metadata.json or {hash}/artifacts.tar.gz
func (b *Backend) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	// Extract hash from key (format: hash/filename)
	hash, _, err := parseKey(key)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err := b.client.GetCache(ctx, hash, buf); err != nil {
		return nil, err
	}

	return io.NopCloser(buf), nil
}

// Put stores data in remote cache
func (b *Backend) Put(ctx context.Context, key string, data io.Reader) error {
	// Extract hash from key
	hash, filename, err := parseKey(key)
	if err != nil {
		return err
	}

	// Only store artifacts.tar.gz - metadata is embedded in the server's database
	if filename != "artifacts.tar.gz" {
		// Skip metadata.json - it's handled separately by the server
		return nil
	}

	// Read all data to get size
	dataBytes, err := io.ReadAll(data)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	size := int64(len(dataBytes))

	// Upload to remote cache
	return b.client.PutCache(ctx, hash, "", bytes.NewReader(dataBytes), size)
}

// Exists checks if data exists in remote cache
func (b *Backend) Exists(ctx context.Context, key string) (bool, error) {
	hash, _, err := parseKey(key)
	if err != nil {
		return false, err
	}

	// Try to get the artifact
	buf := new(bytes.Buffer)
	err = b.client.GetCache(ctx, hash, buf)
	return err == nil, nil
}

// Delete removes data from remote cache
func (b *Backend) Delete(ctx context.Context, key string) error {
	hash, _, err := parseKey(key)
	if err != nil {
		return err
	}

	return b.client.DeleteCache(ctx, hash)
}

// List is not fully implemented for custom backend
func (b *Backend) List(ctx context.Context, prefix string) ([]string, error) {
	// This would require server-side support for listing
	// For now, return empty list
	return []string{}, nil
}

// parseKey extracts hash and filename from a key like "abc123/artifacts.tar.gz"
func parseKey(key string) (hash string, filename string, err error) {
	var i int
	for i = 0; i < len(key); i++ {
		if key[i] == '/' {
			break
		}
	}

	if i == len(key) {
		return "", "", fmt.Errorf("invalid key format: %s", key)
	}

	hash = key[:i]
	filename = key[i+1:]

	return hash, filename, nil
}
