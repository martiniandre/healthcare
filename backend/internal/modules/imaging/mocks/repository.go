package mocks

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/imaging"
)

type MockRepository struct {
	Studies map[uuid.UUID]*imaging.ImagingStudy
	Err     error
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		Studies: make(map[uuid.UUID]*imaging.ImagingStudy),
	}
}

func (mockRepository *MockRepository) CreateImagingStudy(ctx context.Context, study *imaging.ImagingStudy) error {
	if mockRepository.Err != nil {
		return mockRepository.Err
	}
	mockRepository.Studies[study.ID] = study
	return nil
}

func (mockRepository *MockRepository) GetImagingStudy(ctx context.Context, id uuid.UUID) (*imaging.ImagingStudy, error) {
	if mockRepository.Err != nil {
		return nil, mockRepository.Err
	}
	study, exists := mockRepository.Studies[id]
	if !exists {
		return nil, errors.New("imaging study not found")
	}
	return study, nil
}

func (mockRepository *MockRepository) ListImagingStudiesByPatient(ctx context.Context, patientFhirID string) ([]*imaging.ImagingStudy, error) {
	if mockRepository.Err != nil {
		return nil, mockRepository.Err
	}
	var patientStudies []*imaging.ImagingStudy
	for _, study := range mockRepository.Studies {
		if study.PatientFhirID == patientFhirID {
			patientStudies = append(patientStudies, study)
		}
	}
	return patientStudies, nil
}

func (mockRepository *MockRepository) UpdateImagingStudy(ctx context.Context, study *imaging.ImagingStudy) error {
	if mockRepository.Err != nil {
		return mockRepository.Err
	}
	mockRepository.Studies[study.ID] = study
	return nil
}
