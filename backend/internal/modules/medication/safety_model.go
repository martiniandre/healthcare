package medication

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	AlertKindAllergy     = "allergy"
	AlertKindInteraction = "interaction"

	SeverityInfo            = "info"
	SeverityModerate        = "moderate"
	SeverityMajor           = "major"
	SeverityContraindicated = "contraindicated"

	AlertMessageKeyAllergyIngredient = "medicationSafety.allergy.ingredient"
	AlertMessageKeyAllergyClass      = "medicationSafety.allergy.class"

	MinimumOverrideJustificationLength = 10
)

type SafetyAlert struct {
	Kind              string
	Severity          string
	Blocking          bool
	MessageKey        string
	Params            map[string]string
	RelatedResourceID string
}

type CheckSafetyInput struct {
	PatientFHIRID  string
	MedicationCode string
	MedicationName string
}

type PatientAllergy struct {
	FHIRResourceID  string
	AllergenCode    string
	AllergenDisplay string
	ClinicalStatus  string
}

type SafetyOverride struct {
	ID                      uuid.UUID
	PatientFHIRID           string
	MedicationRequestFHIRID string
	MedicationCode          string
	PractitionerUserID      uuid.UUID
	AlertKinds              []string
	Justification           string
	CreatedAt               time.Time
}

type OverrideAuditEvent struct {
	CorrelationID           string
	CallerUserID            string
	CallerRole              string
	PatientFHIRID           string
	MedicationRequestFHIRID string
	MedicationCode          string
	AlertKinds              []string
	RelatedResourceIDs      []string
}

type AllergyReader interface {
	ListPatientAllergies(ctx context.Context, patientFHIRID string) ([]PatientAllergy, error)
}

type SafetyCatalog interface {
	LoadSnapshot(ctx context.Context) (*CatalogSnapshot, error)
}

type OverrideRecorder interface {
	RecordOverride(ctx context.Context, override SafetyOverride) error
	LinkMedicationRequest(ctx context.Context, overrideID uuid.UUID, medicationRequestFHIRID string) error
	DiscardOverride(ctx context.Context, overrideID uuid.UUID) error
}

type AuditRecorder interface {
	RecordOverrideAudit(ctx context.Context, event OverrideAuditEvent) error
}
