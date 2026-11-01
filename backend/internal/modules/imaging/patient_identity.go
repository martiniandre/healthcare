package imaging

import (
	"context"
	"fmt"
	"strings"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/healthcare/backend/internal/shared/healthcare"
)

type PatientIdentity struct {
	FullNames   []string
	DocumentIDs []string
}

type PatientIdentityReader interface {
	ReadPatientIdentity(ctx context.Context, patientFhirID string) (PatientIdentity, error)
}

type fhirPatientIdentityReader struct {
	fhirClient healthcare.FHIRClient
}

func NewFHIRPatientIdentityReader(fhirClient healthcare.FHIRClient) PatientIdentityReader {
	return &fhirPatientIdentityReader{fhirClient: fhirClient}
}

func (reader *fhirPatientIdentityReader) ReadPatientIdentity(ctx context.Context, patientFhirID string) (PatientIdentity, error) {
	rawPatient, getError := reader.fhirClient.GetResource(ctx, "Patient", patientFhirID)
	if getError != nil {
		if healthcare.IsNotFound(getError) {
			return PatientIdentity{}, fmt.Errorf("failed to read patient identity: %w", apperrors.ErrPatientNotFound)
		}
		return PatientIdentity{}, fmt.Errorf("failed to read patient identity: %w", getError)
	}

	decodedPatient, decodeError := fhir.DecodeResource[fhir.Patient](rawPatient)
	if decodeError != nil {
		return PatientIdentity{}, fmt.Errorf("failed to decode patient identity: %w", decodeError)
	}

	identity := PatientIdentity{}
	for _, humanName := range decodedPatient.Name {
		nameParts := append(append([]string{}, humanName.Given...), humanName.Family)
		if fullName := strings.TrimSpace(strings.Join(nameParts, " ")); fullName != "" {
			identity.FullNames = append(identity.FullNames, fullName)
		}
	}
	for _, patientIdentifier := range decodedPatient.Identifier {
		if documentID := strings.TrimSpace(patientIdentifier.Value); documentID != "" {
			identity.DocumentIDs = append(identity.DocumentIDs, documentID)
		}
	}
	return identity, nil
}
