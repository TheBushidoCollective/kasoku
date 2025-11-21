package azure

import (
	"context"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type Backend struct {
	client    *azblob.Client
	container string
}

func New(accountName, container, accountKey string) (*Backend, error) {
	if accountName == "" || container == "" {
		return nil, fmt.Errorf("account name and container are required")
	}

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &Backend{
		client:    client,
		container: container,
	}, nil
}

func (b *Backend) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	resp, err := b.client.DownloadStream(ctx, b.container, key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download blob: %w", err)
	}

	return resp.Body, nil
}

func (b *Backend) Put(ctx context.Context, key string, data io.Reader) error {
	_, err := b.client.UploadStream(ctx, b.container, key, data, nil)
	if err != nil {
		return fmt.Errorf("failed to upload blob: %w", err)
	}

	return nil
}

func (b *Backend) Exists(ctx context.Context, key string) (bool, error) {
	_, err := b.client.DownloadStream(ctx, b.container, key, nil)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (b *Backend) Delete(ctx context.Context, key string) error {
	_, err := b.client.DeleteBlob(ctx, b.container, key, nil)
	if err != nil {
		return fmt.Errorf("failed to delete blob: %w", err)
	}

	return nil
}

func (b *Backend) List(ctx context.Context, prefix string) ([]string, error) {
	var keys []string

	pager := b.client.NewListBlobsFlatPager(b.container, &azblob.ListBlobsFlatOptions{
		Prefix: &prefix,
	})

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list blobs: %w", err)
		}

		for _, blob := range page.Segment.BlobItems {
			keys = append(keys, *blob.Name)
		}
	}

	return keys, nil
}
