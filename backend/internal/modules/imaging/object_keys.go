package imaging

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

const defaultDICOMExtension = ".dcm"

var unsafeObjectPathChars = regexp.MustCompile(`[^A-Za-z0-9-.]`)

func sanitizeLeafName(rawName string) string {
	baseName := strings.TrimSpace(filepath.Base(rawName))
	if baseName == "" || baseName == "." || baseName == string(filepath.Separator) {
		return ""
	}
	return unsafeObjectPathChars.ReplaceAllString(baseName, "_")
}

func sanitizePatientSegment(patientFhirID string) string {
	return unsafeObjectPathChars.ReplaceAllString(patientFhirID, "_")
}

func instanceObjectKey(patientFhirID string, studyID uuid.UUID, instanceID uuid.UUID) string {
	return fmt.Sprintf("pacs/%s/%s/%s%s", sanitizePatientSegment(patientFhirID), studyID.String(), instanceID.String(), defaultDICOMExtension)
}

func legacyStudyObjectKey(study *ImagingStudy) string {
	studyExtension := filepath.Ext(study.FileName)
	if studyExtension == "" {
		studyExtension = defaultDICOMExtension
	}
	return fmt.Sprintf("pacs/%s/%s%s", sanitizePatientSegment(study.PatientFhirID), study.ID.String(), studyExtension)
}
