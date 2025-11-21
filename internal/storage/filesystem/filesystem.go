package filesystem

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Backend struct {
	basePath string
}

func New(basePath string) (*Backend, error) {
	if basePath == "" {
		return nil, fmt.Errorf("base path is required")
	}

	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	return &Backend{basePath: absPath}, nil
}

func (b *Backend) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	path := filepath.Join(b.basePath, key)

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("key not found: %s", key)
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return f, nil
}

func (b *Backend) Put(ctx context.Context, key string, data io.Reader) error {
	path := filepath.Join(b.basePath, key)

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, data); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (b *Backend) Exists(ctx context.Context, key string) (bool, error) {
	path := filepath.Join(b.basePath, key)

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (b *Backend) Delete(ctx context.Context, key string) error {
	path := filepath.Join(b.basePath, key)

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (b *Backend) List(ctx context.Context, prefix string) ([]string, error) {
	var keys []string

	prefixPath := filepath.Join(b.basePath, prefix)
	baseDir := filepath.Dir(prefixPath)

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(b.basePath, path)
		if err != nil {
			return err
		}

		if filepath.HasPrefix(relPath, prefix) {
			keys = append(keys, relPath)
		}

		return nil
	})

	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return keys, nil
}
