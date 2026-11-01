package mocks

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/imaging"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type MockRepository struct {
	Studies    map[uuid.UUID]*imaging.ImagingStudy
	Instances  map[uuid.UUID][]*imaging.InstanceIndex
	Err        error
	WriteCount int
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		Studies:   make(map[uuid.UUID]*imaging.ImagingStudy),
		Instances: make(map[uuid.UUID][]*imaging.InstanceIndex),
	}
}

func (mockRepository *MockRepository) CreateImagingStudy(ctx context.Context, study *imaging.ImagingStudy, instances []*imaging.InstanceIndex) error {
	if mockRepository.Err != nil {
		return mockRepository.Err
	}
	mockRepository.WriteCount++
	for _, existingStudy := range mockRepository.Studies {
		if existingStudy.StudyInstanceUID == study.StudyInstanceUID {
			return apperrors.ErrDICOMDuplicateInstance
		}
	}
	if hasRepeatedSOP(instances) {
		return apperrors.ErrDICOMDuplicateInstance
	}
	mockRepository.Studies[study.ID] = study
	mockRepository.Instances[study.ID] = append([]*imaging.InstanceIndex{}, instances...)
	return nil
}

func (mockRepository *MockRepository) AppendInstances(ctx context.Context, studyID uuid.UUID, instances []*imaging.InstanceIndex, updatedAt time.Time) error {
	if mockRepository.Err != nil {
		return mockRepository.Err
	}
	mockRepository.WriteCount++
	study, exists := mockRepository.Studies[studyID]
	if !exists {
		return apperrors.ErrImagingStudyNotFound
	}
	if hasRepeatedSOP(instances) {
		return apperrors.ErrDICOMDuplicateInstance
	}
	for _, storedInstance := range mockRepository.Instances[studyID] {
		for _, incomingInstance := range instances {
			if storedInstance.SOPInstanceUID == incomingInstance.SOPInstanceUID {
				return apperrors.ErrDICOMDuplicateInstance
			}
		}
	}
	mockRepository.Instances[studyID] = append(mockRepository.Instances[studyID], instances...)
	study.Status = imaging.ImagingStudyStatusReceiving
	study.UpdatedAt = updatedAt
	return nil
}

func hasRepeatedSOP(instances []*imaging.InstanceIndex) bool {
	seen := make(map[string]bool)
	for _, instance := range instances {
		if seen[instance.SOPInstanceUID] {
			return true
		}
		seen[instance.SOPInstanceUID] = true
	}
	return false
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

func (mockRepository *MockRepository) HasInstance(ctx context.Context, studyID uuid.UUID, sopInstanceUID string) (bool, error) {
	if mockRepository.Err != nil {
		return false, mockRepository.Err
	}
	for _, storedInstance := range mockRepository.Instances[studyID] {
		if storedInstance.SOPInstanceUID == sopInstanceUID {
			return true, nil
		}
	}
	return false, nil
}

func (mockRepository *MockRepository) ListInstances(ctx context.Context, studyID uuid.UUID) ([]*imaging.InstanceIndex, error) {
	if mockRepository.Err != nil {
		return nil, mockRepository.Err
	}
	return append([]*imaging.InstanceIndex{}, mockRepository.Instances[studyID]...), nil
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
	mockRepository.WriteCount++
	mockRepository.Studies[study.ID] = study
	return nil
}
