package storage

import (
	"context"
	"fmt"
	"io"
	"time"
)

type StorageClient interface {
	Upload(ctx context.Context, bucketName, objectPath string, reader io.Reader) (string, error)
	GetPresignedURL(ctx context.Context, bucketName, objectPath string, expiration time.Duration) (string, error)
}

type dummyStorageClient struct{}

func NewStorageClient() StorageClient {
	return &dummyStorageClient{}
}

func (client *dummyStorageClient) Upload(ctx context.Context, bucketName, objectPath string, reader io.Reader) (string, error) {
	_, readError := io.Copy(io.Discard, reader)
	if readError != nil {
		return "", readError
	}
	return fmt.Sprintf("gs://%s/%s", bucketName, objectPath), nil
}

func (client *dummyStorageClient) GetPresignedURL(ctx context.Context, bucketName, objectPath string, expiration time.Duration) (string, error) {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s?expiration=%d", bucketName, objectPath, time.Now().Add(expiration).Unix()), nil
}
