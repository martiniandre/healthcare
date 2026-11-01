package imaging

import (
	"strings"
	"unicode"

	"github.com/healthcare/backend/internal/shared/apperrors"
)

const (
	cpfDigitCount        = 11
	nameOverlapThreshold = 0.5
)

type matchOutcome int

const (
	matchInconclusive matchOutcome = iota
	matchConfirmed
	matchRejected
)

var accentFolder = strings.NewReplacer(
	"á", "a", "à", "a", "â", "a", "ã", "a", "ä", "a",
	"é", "e", "è", "e", "ê", "e", "ë", "e",
	"í", "i", "ì", "i", "î", "i", "ï", "i",
	"ó", "o", "ò", "o", "ô", "o", "õ", "o", "ö", "o",
	"ú", "u", "ù", "u", "û", "u", "ü", "u",
	"ç", "c", "ñ", "n",
)

var nameParticles = map[string]bool{
	"de": true, "da": true, "do": true, "das": true, "dos": true,
	"e": true, "di": true, "del": true, "la": true,
}

func evaluatePatientMatch(patientFhirID string, source dicomPatientIdentity, identity PatientIdentity) ([]UploadWarning, error) {
	if source.isAnonymized() {
		return []UploadWarning{newUploadWarning(WarningCodeAnonymizedStudy)}, nil
	}

	identifierOutcome := compareIdentifier(patientFhirID, source.patientID, identity.DocumentIDs)
	if identifierOutcome == matchRejected {
		return nil, apperrors.ErrDICOMPatientMismatch
	}
	if identifierOutcome == matchConfirmed {
		return nil, nil
	}

	switch compareNames(source.patientName, identity.FullNames) {
	case matchRejected:
		return nil, apperrors.ErrDICOMPatientMismatch
	case matchConfirmed:
		return nil, nil
	}
	return []UploadWarning{newUploadWarning(WarningCodePatientIdentityMissing)}, nil
}

type dicomPatientIdentity struct {
	patientID   string
	patientName string
}

func (identity dicomPatientIdentity) isAnonymized() bool {
	normalizedID := strings.ToUpper(strings.TrimSpace(identity.patientID))
	if strings.HasPrefix(normalizedID, "ANON") {
		return true
	}
	if normalizedID != "" {
		return false
	}
	nameTokens := strings.Fields(normalizeText(identity.patientName))
	return len(nameTokens) == 0 || strings.HasPrefix(nameTokens[0], "anon")
}

func compareIdentifier(patientFhirID string, dicomPatientID string, patientDocumentIDs []string) matchOutcome {
	trimmedPatientID := strings.TrimSpace(dicomPatientID)
	if trimmedPatientID == "" {
		return matchInconclusive
	}
	if strings.EqualFold(trimmedPatientID, strings.TrimSpace(patientFhirID)) {
		return matchConfirmed
	}

	dicomDigits, looksLikeDocument := documentDigits(trimmedPatientID)
	if !looksLikeDocument || len(dicomDigits) != cpfDigitCount {
		return matchInconclusive
	}

	comparableDocumentFound := false
	for _, patientDocumentID := range patientDocumentIDs {
		patientDigits, isDocument := documentDigits(patientDocumentID)
		if !isDocument || len(patientDigits) != cpfDigitCount {
			continue
		}
		comparableDocumentFound = true
		if patientDigits == dicomDigits {
			return matchConfirmed
		}
	}
	if comparableDocumentFound {
		return matchRejected
	}
	return matchInconclusive
}

func documentDigits(rawValue string) (string, bool) {
	var digitBuilder strings.Builder
	for _, character := range rawValue {
		switch {
		case unicode.IsDigit(character):
			digitBuilder.WriteRune(character)
		case character == '.' || character == '-' || character == '/' || unicode.IsSpace(character):
		default:
			return "", false
		}
	}
	return digitBuilder.String(), digitBuilder.Len() > 0
}

func compareNames(dicomPatientName string, patientFullNames []string) matchOutcome {
	dicomTokens := significantTokens(dicomNameWithoutExtras(dicomPatientName))
	if len(dicomTokens) == 0 {
		return matchInconclusive
	}

	comparedAny := false
	for _, patientFullName := range patientFullNames {
		patientTokens := significantTokens(patientFullName)
		if len(patientTokens) == 0 {
			continue
		}
		comparedAny = true
		if namesOverlap(dicomTokens, patientTokens) {
			return matchConfirmed
		}
	}
	if comparedAny {
		return matchRejected
	}
	return matchInconclusive
}

func dicomNameWithoutExtras(rawName string) string {
	firstRepresentation := strings.SplitN(rawName, "=", 2)[0]
	nameComponents := strings.Split(firstRepresentation, "^")
	if len(nameComponents) > 2 {
		nameComponents = nameComponents[:2]
	}
	return strings.Join(nameComponents, " ")
}

func namesOverlap(firstTokens map[string]bool, secondTokens map[string]bool) bool {
	sharedCount := 0
	for token := range firstTokens {
		if secondTokens[token] {
			sharedCount++
		}
	}
	smallerSize := len(firstTokens)
	if len(secondTokens) < smallerSize {
		smallerSize = len(secondTokens)
	}
	if sharedCount == smallerSize {
		return true
	}
	return float64(sharedCount)/float64(smallerSize) > nameOverlapThreshold
}

func significantTokens(rawName string) map[string]bool {
	tokens := make(map[string]bool)
	for _, token := range strings.Fields(normalizeText(rawName)) {
		if len(token) < 2 || nameParticles[token] {
			continue
		}
		tokens[token] = true
	}
	return tokens
}

func normalizeText(rawText string) string {
	foldedText := accentFolder.Replace(strings.ToLower(rawText))
	var normalizedBuilder strings.Builder
	for _, character := range foldedText {
		if unicode.IsLetter(character) || unicode.IsDigit(character) {
			normalizedBuilder.WriteRune(character)
		} else {
			normalizedBuilder.WriteRune(' ')
		}
	}
	return normalizedBuilder.String()
}
