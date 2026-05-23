package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"time"
)

type gcsClient struct {
	client *storage.Client
}

func NewGCSClient(ctx context.Context) (StorageClient, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}
	return &gcsClient{client: client}, nil
}

func (gcs *gcsClient) Upload(ctx context.Context, bucketName, objectPath string, reader io.Reader) (string, error) {
	writer := gcs.client.Bucket(bucketName).Object(objectPath).NewWriter(ctx)
	writer.ContentType = "application/dicom"
	writer.ChunkSize = 8 * 1024 * 1024 // 8MB chunks

	if _, err := io.Copy(writer, reader); err != nil {
		writer.Close()
		return "", fmt.Errorf("upload to GCS failed: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to finalize GCS upload: %w", err)
	}
	return fmt.Sprintf("gs://%s/%s", bucketName, objectPath), nil
}

func (gcs *gcsClient) GetPresignedURL(ctx context.Context, bucketName, objectPath string, expiration time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(expiration),
	}
	url, err := gcs.client.Bucket(bucketName).SignedURL(objectPath, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return url, nil
}
