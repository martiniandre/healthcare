package exam_analyzer

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var unsafeExamPathChars = regexp.MustCompile(`[^A-Za-z0-9-.]`)

func patientIDOrDefault(patientFhirID *string) string {
	if patientFhirID == nil || strings.TrimSpace(*patientFhirID) == "" {
		return "unspecified"
	}
	return *patientFhirID
}

func examFileExtension(fileName string) string {
	extension := filepath.Ext(strings.ToLower(fileName))
	if extension == "" {
		return ".png"
	}
	return extension
}

func buildExamStorageKey(analysisID uuid.UUID, patientFhirID, fileName string) string {
	sanitizedPatientID := unsafeExamPathChars.ReplaceAllString(patientFhirID, "_")
	return fmt.Sprintf("exams/%s/%s%s", sanitizedPatientID, analysisID.String(), examFileExtension(fileName))
}

func examContentType(fileName string) string {
	switch strings.ToLower(filepath.Ext(fileName)) {
	case ".pdf":
		return "application/pdf"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	default:
		return "image/png"
	}
}
