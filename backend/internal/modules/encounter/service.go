package encounter

import (
	"context"

	"github.com/healthcare/backend/internal/shared/eventbus"
)

type Service interface {
	CreateEncounter(ctx context.Context, encounter *Encounter) (*Encounter, error)
	GetEncounter(ctx context.Context, fhirResourceID string) (*Encounter, error)
	GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*Encounter, error)
}

type service struct {
	repo     Repository
	eventBus eventbus.Bus
}

func NewService(repo Repository, eventBus eventbus.Bus) Service {
	return &service{repo: repo, eventBus: eventBus}
}

func (encounterService *service) CreateEncounter(ctx context.Context, encounter *Encounter) (*Encounter, error) {
	if encounter.PatientFHIRID == "" {
		return nil, ErrEncounterNotFound
	}
	createdEncounter, err := encounterService.repo.CreateEncounter(ctx, encounter)
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

func (encounterService *service) GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*Encounter, error) {
	return encounterService.repo.GetEncountersByPatient(ctx, patientFHIRID)
}
