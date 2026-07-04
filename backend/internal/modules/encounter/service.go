package encounter

import (
	"context"
)

type Service interface {
	CreateEncounter(ctx context.Context, encounter *Encounter) (*Encounter, error)
	GetEncounter(ctx context.Context, fhirResourceID string) (*Encounter, error)
	GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*Encounter, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (encounterService *service) CreateEncounter(ctx context.Context, encounter *Encounter) (*Encounter, error) {
	if encounter.PatientFHIRID == "" {
		return nil, ErrEncounterNotFound
	}
	return encounterService.repo.CreateEncounter(ctx, encounter)
}

func (encounterService *service) GetEncounter(ctx context.Context, fhirResourceID string) (*Encounter, error) {
	return encounterService.repo.GetEncounterByID(ctx, fhirResourceID)
}

func (encounterService *service) GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*Encounter, error) {
	return encounterService.repo.GetEncountersByPatient(ctx, patientFHIRID)
}
