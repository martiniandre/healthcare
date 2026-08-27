package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"strings"
	"time"
)

type gcsClient struct {
	client     *storage.Client
	bucketName string
	rootPrefix string
}

func NewGCSClient(ctx context.Context, bucketName, rootPrefix string) (StorageClient, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}
	return &gcsClient{
		client:     client,
		bucketName: bucketName,
		rootPrefix: strings.Trim(rootPrefix, "/"),
	}, nil
}

func (gcs *gcsClient) Upload(ctx context.Context, key string, reader io.Reader, contentType string) error {
	objectPath := joinObjectPath(gcs.rootPrefix, key)
	writer := gcs.client.Bucket(gcs.bucketName).Object(objectPath).NewWriter(ctx)
	writer.ContentType = contentType
	writer.ChunkSize = 8 * 1024 * 1024

	if _, err := io.Copy(writer, reader); err != nil {
		writer.Close()
		return fmt.Errorf("upload to GCS failed: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to finalize GCS upload: %w", err)
	}
	return nil
}

func (gcs *gcsClient) GetPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	objectPath := joinObjectPath(gcs.rootPrefix, key)
	opts := &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(expiration),
	}
	url, err := gcs.client.Bucket(gcs.bucketName).SignedURL(objectPath, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return url, nil
}

func (gcs *gcsClient) Delete(ctx context.Context, key string) error {
	objectPath := joinObjectPath(gcs.rootPrefix, key)
	if err := gcs.client.Bucket(gcs.bucketName).Object(objectPath).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete object from GCS: %w", err)
	}
	return nil
}

func (gcs *gcsClient) Read(ctx context.Context, key string) (io.ReadCloser, error) {
	objectPath := joinObjectPath(gcs.rootPrefix, key)
	reader, err := gcs.client.Bucket(gcs.bucketName).Object(objectPath).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read object from GCS: %w", err)
	}
	return reader, nil
}
