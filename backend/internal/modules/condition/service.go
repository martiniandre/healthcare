package condition

import (
	"context"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/validator"
)

type Service interface {
	CreateCondition(ctx context.Context, input CreateConditionInput) (*Condition, error)
	GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]*Condition, error)
	UpdateCondition(ctx context.Context, fhirResourceID string, input UpdateConditionInput) (*Condition, error)
	DeleteCondition(ctx context.Context, fhirResourceID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (conditionService *service) CreateCondition(ctx context.Context, input CreateConditionInput) (*Condition, error) {
	if fieldViolations := validateConditionFields(input.PatientFHIRID, input.ICD10Code, input.CodeDisplay, input.ClinicalStatus); len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid condition input", fieldViolations)
	}

	newCondition := &Condition{
		PatientFHIRID:   input.PatientFHIRID,
		EncounterFHIRID: input.EncounterFHIRID,
		ICD10Code:       input.ICD10Code,
		CodeDisplay:     input.CodeDisplay,
		ClinicalStatus:  normalizeClinicalStatus(input.ClinicalStatus),
	}

	return conditionService.repo.CreateCondition(ctx, newCondition)
}

func (conditionService *service) UpdateCondition(ctx context.Context, fhirResourceID string, input UpdateConditionInput) (*Condition, error) {
	currentCondition, fetchErr := conditionService.repo.GetConditionByID(ctx, fhirResourceID)
	if fetchErr != nil {
		return nil, fetchErr
	}

	mergedCondition := mergeConditionInput(currentCondition, input)

	if fieldViolations := validateConditionFields(mergedCondition.PatientFHIRID, mergedCondition.ICD10Code, mergedCondition.CodeDisplay, mergedCondition.ClinicalStatus); len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid condition input", fieldViolations)
	}
	mergedCondition.ClinicalStatus = normalizeClinicalStatus(mergedCondition.ClinicalStatus)

	return conditionService.repo.UpdateCondition(ctx, fhirResourceID, mergedCondition)
}

func (conditionService *service) DeleteCondition(ctx context.Context, fhirResourceID string) error {
	return conditionService.repo.DeleteCondition(ctx, fhirResourceID)
}

func (conditionService *service) GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]*Condition, error) {
	return conditionService.repo.GetConditionsByPatient(ctx, patientFHIRID)
}

func validateConditionFields(patientFHIRID string, icd10Code string, codeDisplay string, clinicalStatus string) map[string]string {
	fieldViolations := make(map[string]string)
	if patientFHIRID == "" {
		fieldViolations["patient_fhir_id"] = "is required"
	}
	if icd10Code == "" {
		fieldViolations["icd10_code"] = "is required"
	} else if !validator.IsValidICD10(icd10Code) {
		fieldViolations["icd10_code"] = "is not a valid ICD-10 code"
	}
	if codeDisplay == "" {
		fieldViolations["code_display"] = "is required"
	}
	if clinicalStatus != "" && !validator.IsValidClinicalStatus(clinicalStatus) {
		fieldViolations["clinical_status"] = "is not a valid clinical status"
	}
	return fieldViolations
}

func normalizeClinicalStatus(clinicalStatus string) string {
	if clinicalStatus == "" {
		return "active"
	}
	return clinicalStatus
}

func mergeConditionInput(currentCondition *Condition, input UpdateConditionInput) *Condition {
	mergedCondition := *currentCondition
	if input.ICD10Code != nil {
		mergedCondition.ICD10Code = *input.ICD10Code
	}
	if input.CodeDisplay != nil {
		mergedCondition.CodeDisplay = *input.CodeDisplay
	}
	if input.ClinicalStatus != nil {
		mergedCondition.ClinicalStatus = *input.ClinicalStatus
	}
	if input.EncounterFHIRID != nil {
		mergedCondition.EncounterFHIRID = *input.EncounterFHIRID
	}
	return &mergedCondition
}
