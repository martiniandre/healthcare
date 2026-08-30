package storage

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"
)

type StorageClient interface {
	Upload(ctx context.Context, key string, reader io.Reader, contentType string) error
	GetPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error)
	Delete(ctx context.Context, key string) error
	Read(ctx context.Context, key string) (io.ReadCloser, error)
}

func joinObjectPath(segments ...string) string {
	nonEmptySegments := make([]string, 0, len(segments))
	for _, segment := range segments {
		if trimmedSegment := strings.Trim(segment, "/"); trimmedSegment != "" {
			nonEmptySegments = append(nonEmptySegments, trimmedSegment)
		}
	}
	return strings.Join(nonEmptySegments, "/")
}

type dummyStorageClient struct {
	bucketName string
	rootPrefix string
}

func NewStorageClient(bucketName, rootPrefix string) StorageClient {
	return &dummyStorageClient{
		bucketName: bucketName,
		rootPrefix: strings.Trim(rootPrefix, "/"),
	}
}

func (client *dummyStorageClient) Upload(ctx context.Context, key string, reader io.Reader, contentType string) error {
	_, readError := io.Copy(io.Discard, reader)
	if readError != nil {
		return readError
	}
	return nil
}

func (client *dummyStorageClient) GetPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s?expiration=%d", client.bucketName, joinObjectPath(client.rootPrefix, key), time.Now().Add(expiration).Unix()), nil
}

func (client *dummyStorageClient) Delete(ctx context.Context, key string) error {
	return nil
}

func (client *dummyStorageClient) Read(ctx context.Context, key string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("")), nil
}
