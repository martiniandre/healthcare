package exam_analyzer

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestExamFileExtension(t *testing.T) {
	testCases := []struct {
		name      string
		fileName  string
		expected  string
	}{
		{name: "png extension", fileName: "chest_xray.png", expected: ".png"},
		{name: "uppercase pdf extension", fileName: "result.PDF", expected: ".pdf"},
		{name: "missing extension defaults to png", fileName: "scan", expected: ".png"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if result := examFileExtension(testCase.fileName); result != testCase.expected {
				t.Fatalf("expected %q, got %q", testCase.expected, result)
			}
		})
	}
}

func TestExamContentType(t *testing.T) {
	testCases := []struct {
		name     string
		fileName string
		expected string
	}{
		{name: "pdf", fileName: "report.pdf", expected: "application/pdf"},
		{name: "jpg", fileName: "photo.jpg", expected: "image/jpeg"},
		{name: "jpeg", fileName: "photo.jpeg", expected: "image/jpeg"},
		{name: "png default", fileName: "photo.png", expected: "image/png"},
		{name: "unknown defaults to png", fileName: "photo.xyz", expected: "image/png"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if result := examContentType(testCase.fileName); result != testCase.expected {
				t.Fatalf("expected %q, got %q", testCase.expected, result)
			}
		})
	}
}

func TestPatientIDOrDefault(t *testing.T) {
	patientID := "patient-123"
	if result := patientIDOrDefault(&patientID); result != "patient-123" {
		t.Fatalf("expected patient id to be preserved, got %q", result)
	}
	if result := patientIDOrDefault(nil); result != "unspecified" {
		t.Fatalf("expected nil patient id to default to unspecified, got %q", result)
	}
	emptyPatientID := ""
	if result := patientIDOrDefault(&emptyPatientID); result != "unspecified" {
		t.Fatalf("expected empty patient id to default to unspecified, got %q", result)
	}
}

func TestBuildExamStorageKey(t *testing.T) {
	analysisID := uuid.New()
	key := buildExamStorageKey(analysisID, "patient-123", "chest_xray.png")

	expected := "exams/patient-123/" + analysisID.String() + ".png"
	if key != expected {
		t.Fatalf("expected key %q, got %q", expected, key)
	}
}

func TestBuildExamStorageKeySanitizesUnsafeSegments(t *testing.T) {
	analysisID := uuid.New()
	key := buildExamStorageKey(analysisID, "patient/123", "chest xray.png")

	segments := strings.Split(key, "/")
	if len(segments) != 3 {
		t.Fatalf("expected 3 key segments, got %d from %q", len(segments), key)
	}

	for _, segment := range segments[1:] {
		if strings.ContainsAny(segment, "/ ") {
			t.Fatalf("expected segment %q to contain no path separators or spaces, key %q", segment, key)
		}
	}
}
