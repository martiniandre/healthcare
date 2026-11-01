package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/imaging"
)

type MockPatientIdentityReader struct {
	Identity imaging.PatientIdentity
	Err      error
	Calls    int
}

func (mockReader *MockPatientIdentityReader) ReadPatientIdentity(ctx context.Context, patientFhirID string) (imaging.PatientIdentity, error) {
	mockReader.Calls++
	if mockReader.Err != nil {
		return imaging.PatientIdentity{}, mockReader.Err
	}
	return mockReader.Identity, nil
}
