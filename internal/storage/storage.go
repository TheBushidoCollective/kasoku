package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/thebushidocollective/kasoku/internal/storage/azure"
	"github.com/thebushidocollective/kasoku/internal/storage/custom"
	"github.com/thebushidocollective/kasoku/internal/storage/filesystem"
	"github.com/thebushidocollective/kasoku/internal/storage/gcs"
	"github.com/thebushidocollective/kasoku/internal/storage/s3"
)

type Backend interface {
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Put(ctx context.Context, key string, data io.Reader) error
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) ([]string, error)
}

type Config struct {
	Type string

	// Custom Kasoku server
	CustomURL   string
	CustomToken string

	FilesystemPath string

	S3Bucket   string
	S3Region   string
	S3Endpoint string

	GCSBucket      string
	GCSCredentials string

	AzureAccount   string
	AzureContainer string
	AzureKey       string
}

func NewBackend(cfg Config) (Backend, error) {
	switch cfg.Type {
	case "custom":
		if cfg.CustomURL == "" {
			return nil, fmt.Errorf("custom backend URL is required")
		}
		if cfg.CustomToken == "" {
			return nil, fmt.Errorf("custom backend token is required")
		}
		return custom.New(cfg.CustomURL, cfg.CustomToken)

	case "filesystem":
		if cfg.FilesystemPath == "" {
			return nil, fmt.Errorf("filesystem path is required")
		}
		return filesystem.New(cfg.FilesystemPath)

	case "s3":
		if cfg.S3Bucket == "" {
			return nil, fmt.Errorf("S3 bucket is required")
		}
		if cfg.S3Region == "" {
			cfg.S3Region = "us-east-1"
		}
		return s3.New(cfg.S3Bucket, cfg.S3Region, cfg.S3Endpoint)

	case "gcs":
		if cfg.GCSBucket == "" {
			return nil, fmt.Errorf("GCS bucket is required")
		}
		return gcs.New(cfg.GCSBucket, cfg.GCSCredentials)

	case "azure":
		if cfg.AzureAccount == "" || cfg.AzureContainer == "" {
			return nil, fmt.Errorf("Azure account and container are required")
		}
		return azure.New(cfg.AzureAccount, cfg.AzureContainer, cfg.AzureKey)

	default:
		return nil, fmt.Errorf("unsupported storage backend type: %s", cfg.Type)
	}
}
