package patients

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/eventbus"
)

var ErrPatientNotFound = errors.New("patient not found")
var ErrPatientAlreadyExists = errors.New("patient with this document already exists")

type Service interface {
	CreatePatient(ctx context.Context, fullName, birthDate, documentID, phoneNumber string) (*Patient, error)
	GetPatient(ctx context.Context, fhirResourceID string) (*Patient, error)
	GetPatientByDocument(ctx context.Context, documentID string) (*Patient, error)
	ListPatients(ctx context.Context, search string, sortField string, sortDirection string, page int, limit int) ([]*Patient, error)
}

type service struct {
	repo     Repository
	eventBus eventbus.Bus
}

func NewService(repo Repository, eventBus eventbus.Bus) Service {
	return &service{repo: repo, eventBus: eventBus}
}

func (patientService *service) CreatePatient(ctx context.Context, fullName, birthDate, documentID, phoneNumber string) (*Patient, error) {
	existingPatient, _ := patientService.repo.GetPatientByDocumentID(ctx, documentID)
	if existingPatient != nil {
		return nil, ErrPatientAlreadyExists
	}

	parsedBirthDate, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return nil, errors.New("invalid birth date format, expected YYYY-MM-DD")
	}
	if !parsedBirthDate.Before(time.Now()) {
		return nil, errors.New("birth date must be in the past")
	}

	patient := &Patient{
		ID:          uuid.New(),
		FullName:    fullName,
		BirthDate:   parsedBirthDate,
		DocumentID:  documentID,
		PhoneNumber: phoneNumber,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdPatient, err := patientService.repo.CreatePatient(ctx, patient)
	if err != nil {
		return nil, err
	}

	if patientService.eventBus != nil {
		patientService.eventBus.Publish(ctx, eventbus.Event{
		Name: "patient.created",
		Data: map[string]any{
			"title":         "Novo Paciente Cadastrado",
			"body":          "Paciente " + createdPatient.FullName + " foi cadastrado no sistema.",
			"resource_type": "patient",
			"resource_id":   createdPatient.ID.String(),
		},
	})
	}

	return createdPatient, nil
}

func (patientService *service) GetPatient(ctx context.Context, fhirResourceID string) (*Patient, error) {
	patient, err := patientService.repo.GetPatientByID(ctx, fhirResourceID)
	if err != nil {
		return nil, ErrPatientNotFound
	}
	return patient, nil
}

func (patientService *service) GetPatientByDocument(ctx context.Context, documentID string) (*Patient, error) {
	patient, err := patientService.repo.GetPatientByDocumentID(ctx, documentID)
	if err != nil {
		return nil, ErrPatientNotFound
	}
	return patient, nil
}

func (patientService *service) ListPatients(ctx context.Context, search string, sortField string, sortDirection string, page int, limit int) ([]*Patient, error) {
	return patientService.repo.ListPatients(ctx, search, sortField, sortDirection, page, limit)
}
