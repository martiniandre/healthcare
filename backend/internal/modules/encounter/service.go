package encounter

import (
	"context"
	"time"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/eventbus"
)

type Service interface {
	CreateEncounter(ctx context.Context, input CreateEncounterInput) (*Encounter, error)
	GetEncounter(ctx context.Context, fhirResourceID string) (*Encounter, error)
	GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*Encounter, error)
	UpdateEncounter(ctx context.Context, fhirResourceID string, input UpdateEncounterInput) (*Encounter, error)
	DeleteEncounter(ctx context.Context, fhirResourceID string) error
}

type service struct {
	repo     Repository
	eventBus eventbus.Bus
}

func NewService(repo Repository, eventBus eventbus.Bus) Service {
	return &service{repo: repo, eventBus: eventBus}
}

func (encounterService *service) CreateEncounter(ctx context.Context, input CreateEncounterInput) (*Encounter, error) {
	fieldViolations := make(map[string]string)
	if input.PatientFHIRID == "" {
		fieldViolations["patient_fhir_id"] = "is required"
	}
	if input.PractitionerID == "" {
		fieldViolations["practitioner_id"] = "is required"
	}
	if input.ReasonDisplay == "" {
		fieldViolations["reason"] = "is required"
	}
	if len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid encounter input", fieldViolations)
	}

	newEncounter := &Encounter{
		PatientFHIRID:  input.PatientFHIRID,
		PractitionerID: input.PractitionerID,
		ReasonCode:     input.ReasonCode,
		ReasonDisplay:  input.ReasonDisplay,
		Status:         "in-progress",
	}
	if newEncounter.StartedAt.IsZero() {
		newEncounter.StartedAt = time.Now()
	}

	createdEncounter, err := encounterService.repo.CreateEncounter(ctx, newEncounter)
	if err != nil {
		return nil, err
	}

	if encounterService.eventBus != nil {
		encounterService.eventBus.Publish(ctx, eventbus.Event{
			Name: "encounter.created",
			Data: map[string]any{
				"title":         "Novo Atendimento Criado",
				"body":          "Atendimento para paciente " + createdEncounter.PatientFHIRID + " foi registrado.",
				"resource_type": "encounter",
				"resource_id":   createdEncounter.FHIRResourceID,
			},
		})
	}

	return createdEncounter, nil
}

func (encounterService *service) GetEncounter(ctx context.Context, fhirResourceID string) (*Encounter, error) {
	return encounterService.repo.GetEncounterByID(ctx, fhirResourceID)
}

func (encounterService *service) UpdateEncounter(ctx context.Context, fhirResourceID string, input UpdateEncounterInput) (*Encounter, error) {
	currentEncounter, fetchErr := encounterService.repo.GetEncounterByID(ctx, fhirResourceID)
	if fetchErr != nil {
		return nil, fetchErr
	}

	mergedEncounter := mergeEncounterInput(currentEncounter, input)

	if input.Status != nil && *input.Status != currentEncounter.Status && !isAllowedStatusTransition(currentEncounter.Status, *input.Status) {
		return nil, apperrors.InvalidArgument("invalid encounter status transition", nil)
	}

	return encounterService.repo.UpdateEncounter(ctx, fhirResourceID, mergedEncounter)
}

func (encounterService *service) DeleteEncounter(ctx context.Context, fhirResourceID string) error {
	return encounterService.repo.DeleteEncounter(ctx, fhirResourceID)
}

func (encounterService *service) GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*Encounter, error) {
	return encounterService.repo.GetEncountersByPatient(ctx, patientFHIRID)
}

func mergeEncounterInput(currentEncounter *Encounter, input UpdateEncounterInput) *Encounter {
	mergedEncounter := *currentEncounter
	if input.ReasonCode != nil {
		mergedEncounter.ReasonCode = *input.ReasonCode
	}
	if input.ReasonDisplay != nil {
		mergedEncounter.ReasonDisplay = *input.ReasonDisplay
	}
	if input.PractitionerID != nil {
		mergedEncounter.PractitionerID = *input.PractitionerID
	}
	if input.Status != nil {
		mergedEncounter.Status = *input.Status
	}
	return &mergedEncounter
}

func isAllowedStatusTransition(currentStatus string, targetStatus string) bool {
	if currentStatus != "in-progress" {
		return false
	}
	return targetStatus == "finished" || targetStatus == "cancelled"
}
