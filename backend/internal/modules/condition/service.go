package condition

import (
	"context"
	"errors"

	"github.com/healthcare/backend/internal/shared/validator"
)

type Service interface {
	CreateCondition(ctx context.Context, condition *Condition) (*Condition, error)
	GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]*Condition, error)
	UpdateCondition(ctx context.Context, fhirResourceID string, condition *Condition) (*Condition, error)
	DeleteCondition(ctx context.Context, fhirResourceID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (conditionService *service) CreateCondition(ctx context.Context, condition *Condition) (*Condition, error) {
	if condition.PatientFHIRID == "" || condition.ICD10Code == "" {
		return nil, ErrConditionNotFound
	}
	if !validator.IsValidICD10(condition.ICD10Code) {
		return nil, errors.New("invalid ICD-10 format")
	}
	if condition.ClinicalStatus == "" {
		condition.ClinicalStatus = "active"
	} else if !validator.IsValidClinicalStatus(condition.ClinicalStatus) {
		return nil, errors.New("invalid clinical status")
	}
	return conditionService.repo.CreateCondition(ctx, condition)
}

func (conditionService *service) UpdateCondition(ctx context.Context, fhirResourceID string, condition *Condition) (*Condition, error) {
	if condition.PatientFHIRID == "" || condition.ICD10Code == "" {
		return nil, ErrConditionNotFound
	}
	if !validator.IsValidICD10(condition.ICD10Code) {
		return nil, errors.New("invalid ICD-10 format")
	}
	if condition.ClinicalStatus == "" {
		condition.ClinicalStatus = "active"
	} else if !validator.IsValidClinicalStatus(condition.ClinicalStatus) {
		return nil, errors.New("invalid clinical status")
	}
	return conditionService.repo.UpdateCondition(ctx, fhirResourceID, condition)
}

func (conditionService *service) DeleteCondition(ctx context.Context, fhirResourceID string) error {
	return conditionService.repo.DeleteCondition(ctx, fhirResourceID)
}

func (conditionService *service) GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]*Condition, error) {
	return conditionService.repo.GetConditionsByPatient(ctx, patientFHIRID)
}
