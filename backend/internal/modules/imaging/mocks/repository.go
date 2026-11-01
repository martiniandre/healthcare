package mocks

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/imaging"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type MockRepository struct {
	Studies         map[uuid.UUID]*imaging.ImagingStudy
	InstanceIndexes map[uuid.UUID]*imaging.InstanceIndex
	Err             error
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		Studies:         make(map[uuid.UUID]*imaging.ImagingStudy),
		InstanceIndexes: make(map[uuid.UUID]*imaging.InstanceIndex),
	}
}

func (mockRepository *MockRepository) CreateImagingStudy(ctx context.Context, study *imaging.ImagingStudy, instanceIndex *imaging.InstanceIndex) error {
	if mockRepository.Err != nil {
		return mockRepository.Err
	}
	for _, existingStudy := range mockRepository.Studies {
		if existingStudy.StudyInstanceUID == study.StudyInstanceUID {
			return apperrors.ErrDICOMDuplicateInstance
		}
	}
	mockRepository.Studies[study.ID] = study
	if instanceIndex != nil {
		mockRepository.InstanceIndexes[study.ID] = instanceIndex
	}
	return nil
}

func (mockRepository *MockRepository) GetImagingStudy(ctx context.Context, id uuid.UUID) (*imaging.ImagingStudy, error) {
	if mockRepository.Err != nil {
		return nil, mockRepository.Err
	}
	study, exists := mockRepository.Studies[id]
	if !exists {
		return nil, apperrors.ErrImagingStudyNotFound
	}
	return study, nil
}

func (mockRepository *MockRepository) FindImagingStudyByStudyInstanceUID(ctx context.Context, studyInstanceUID string) (*imaging.ImagingStudy, error) {
	if mockRepository.Err != nil {
		return nil, mockRepository.Err
	}
	for _, study := range mockRepository.Studies {
		if study.StudyInstanceUID == studyInstanceUID {
			return study, nil
		}
	}
	return nil, fmt.Errorf("find imaging study: %w", apperrors.ErrImagingStudyNotFound)
}

func (mockRepository *MockRepository) GetInstanceIndex(ctx context.Context, studyID uuid.UUID) (*imaging.InstanceIndex, error) {
	if mockRepository.Err != nil {
		return nil, mockRepository.Err
	}
	instanceIndex, exists := mockRepository.InstanceIndexes[studyID]
	if !exists {
		return nil, apperrors.ErrImagingStudyNotFound
	}
	return instanceIndex, nil
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
