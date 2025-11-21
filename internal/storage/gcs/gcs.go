package gcs

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Backend struct {
	client *storage.Client
	bucket string
}

func New(bucket, credentialsFile string) (*Backend, error) {
	if bucket == "" {
		return nil, fmt.Errorf("bucket is required")
	}

	ctx := context.Background()

	var opts []option.ClientOption
	if credentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(credentialsFile))
	}

	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	return &Backend{
		client: client,
		bucket: bucket,
	}, nil
}

func (b *Backend) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	obj := b.client.Bucket(b.bucket).Object(key)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return reader, nil
}

func (b *Backend) Put(ctx context.Context, key string, data io.Reader) error {
	obj := b.client.Bucket(b.bucket).Object(key)
	writer := obj.NewWriter(ctx)

	if _, err := io.Copy(writer, data); err != nil {
		writer.Close()
		return fmt.Errorf("failed to write object: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}

func (b *Backend) Exists(ctx context.Context, key string) (bool, error) {
	obj := b.client.Bucket(b.bucket).Object(key)
	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (b *Backend) Delete(ctx context.Context, key string) error {
	obj := b.client.Bucket(b.bucket).Object(key)
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

func (b *Backend) List(ctx context.Context, prefix string) ([]string, error) {
	var keys []string

	query := &storage.Query{Prefix: prefix}
	it := b.client.Bucket(b.bucket).Objects(ctx, query)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}

		keys = append(keys, attrs.Name)
	}

	return keys, nil
}

func (b *Backend) Close() error {
	return b.client.Close()
}
