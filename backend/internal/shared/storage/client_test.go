package storage

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"
)

func TestDummyStorageClientUpload(t *testing.T) {
	client := NewStorageClient("test-bucket", "prod")

	err := client.Upload(context.Background(), "pacs/patient-1/abc.dcm", strings.NewReader("dicom-data"), "application/dicom")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDummyStorageClientGetPresignedURL(t *testing.T) {
	client := NewStorageClient("test-bucket", "prod")

	url, err := client.GetPresignedURL(context.Background(), "pacs/patient-1/abc.dcm", 15*time.Minute)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.HasPrefix(url, "https://storage.googleapis.com/test-bucket/prod/pacs/patient-1/abc.dcm?") {
		t.Fatalf("unexpected presigned url: %s", url)
	}
}

func TestDummyStorageClientDelete(t *testing.T) {
	client := NewStorageClient("test-bucket", "prod")

	if err := client.Delete(context.Background(), "exams/patient-1/abc.png"); err != nil {
		t.Fatalf("expected no error from delete, got %v", err)
	}
}

func TestDummyStorageClientReadReturnsReader(t *testing.T) {
	client := NewStorageClient("test-bucket", "prod")

	reader, err := client.Read(context.Background(), "exams/patient-1/abc.png")
	if err != nil {
		t.Fatalf("expected no error from read, got %v", err)
	}
	defer reader.Close()

	content, readError := io.ReadAll(reader)
	if readError != nil {
		t.Fatalf("expected no read error, got %v", readError)
	}
	if len(content) != 0 {
		t.Fatalf("expected empty dummy content, got %q", string(content))
	}
}

func TestJoinObjectPath(t *testing.T) {
	testCases := []struct {
		name     string
		segments []string
		expected string
	}{
		{name: "single segment", segments: []string{"pacs"}, expected: "pacs"},
		{name: "multiple segments", segments: []string{"prod", "pacs", "patient-1", "abc.dcm"}, expected: "prod/pacs/patient-1/abc.dcm"},
		{name: "empty middle segment", segments: []string{"prod", "", "pacs"}, expected: "prod/pacs"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if result := joinObjectPath(testCase.segments...); result != testCase.expected {
				t.Fatalf("expected %q, got %q", testCase.expected, result)
			}
		})
	}
}
