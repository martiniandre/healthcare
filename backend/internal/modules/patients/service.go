package patients

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/healthcare/backend/internal/shared/validator"
)

type Service interface {
	CreatePatient(ctx context.Context, input CreatePatientInput) (*Patient, error)
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

func (patientService *service) CreatePatient(ctx context.Context, input CreatePatientInput) (*Patient, error) {
	if fieldViolations := validatePatientFields(input); len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid patient input", fieldViolations)
	}

	existingPatient, _ := patientService.repo.GetPatientByDocumentID(ctx, input.DocumentID)
	if existingPatient != nil {
		return nil, apperrors.ErrPatientAlreadyExists
	}

	parsedBirthDate, _ := time.Parse("2006-01-02", input.BirthDate)

	patient := &Patient{
		ID:          uuid.New(),
		FullName:    input.FullName,
		BirthDate:   parsedBirthDate,
		DocumentID:  input.DocumentID,
		PhoneNumber: input.PhoneNumber,
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
	return patientService.repo.GetPatientByID(ctx, fhirResourceID)
}

func (patientService *service) GetPatientByDocument(ctx context.Context, documentID string) (*Patient, error) {
	return patientService.repo.GetPatientByDocumentID(ctx, documentID)
}

func (patientService *service) ListPatients(ctx context.Context, search string, sortField string, sortDirection string, page int, limit int) ([]*Patient, error) {
	return patientService.repo.ListPatients(ctx, search, sortField, sortDirection, page, limit)
}

func validatePatientFields(input CreatePatientInput) map[string]string {
	fieldViolations := make(map[string]string)
	if strings.TrimSpace(input.FullName) == "" {
		fieldViolations["full_name"] = "is required"
	}
	if strings.TrimSpace(input.BirthDate) == "" {
		fieldViolations["birth_date"] = "is required"
	} else if parsedBirthDate, err := time.Parse("2006-01-02", input.BirthDate); err != nil {
		fieldViolations["birth_date"] = "invalid date format, expected YYYY-MM-DD"
	} else if !parsedBirthDate.Before(time.Now()) {
		fieldViolations["birth_date"] = "must be in the past"
	}
	if strings.TrimSpace(input.DocumentID) == "" || !validator.IsValidCPF(input.DocumentID) {
		fieldViolations["document_id"] = "invalid CPF format"
	}
	if strings.TrimSpace(input.PhoneNumber) == "" || !validator.IsValidPhone(input.PhoneNumber) {
		fieldViolations["phone_number"] = "invalid phone format"
	}
	return fieldViolations
}
