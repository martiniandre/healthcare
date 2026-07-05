package staff

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/healthcare/backend/internal/shared/role"
	"github.com/healthcare/backend/internal/shared/validator"
)

var ErrEmployeeNotFound = errors.New("employee not found")

type Service interface {
	CreateEmployee(ctx context.Context, createdBy uuid.UUID, fullName, email, requestedRole, crmNumber string) (*Employee, error)
	GetEmployee(ctx context.Context, employeeID uuid.UUID) (*Employee, error)
	ListEmployees(ctx context.Context, search string, role string) ([]*Employee, error)
	DeactivateEmployee(ctx context.Context, employeeID uuid.UUID) error
}

type service struct {
	repo       Repository
	fhirClient healthcare.FHIRClient
}

func NewService(repo Repository, fhirClient healthcare.FHIRClient) Service {
	return &service{repo: repo, fhirClient: fhirClient}
}

func (staffService *service) CreateEmployee(ctx context.Context, createdBy uuid.UUID, fullName, email, requestedRole, crmNumber string) (*Employee, error) {
	parsedRole, roleIsValid := role.ParseRole(requestedRole)
	if !roleIsValid {
		return nil, role.ErrInvalidRole
	}

	if !validator.IsValidEmail(email) {
		return nil, errors.New("invalid email format")
	}

	var crmNumberPtr *string
	if crmNumber != "" {
		if !validator.IsValidCRMNumber(crmNumber) {
			return nil, errors.New("invalid CRM format")
		}
		crmNumberPtr = &crmNumber
	}

	parsedCreatedBy := createdBy
	employee := &Employee{
		ID:        uuid.New(),
		FullName:  fullName,
		Email:     email,
		Role:      parsedRole,
		CRMNumber: crmNumberPtr,
		CreatedBy: &parsedCreatedBy,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := staffService.repo.CreateEmployee(ctx, employee)
	if err != nil {
		return nil, err
	}

	if staffService.fhirClient != nil {
		practitionerResource := fhir.NewPractitionerResource(fullName, crmNumber)
		responseBody, fhirErr := staffService.fhirClient.CreateResource(ctx, "Practitioner", practitionerResource)
		if fhirErr != nil {
			return nil, fmt.Errorf("failed to create practitioner in healthcare api: %w", fhirErr)
		}

		var createdResource map[string]interface{}
		if parseErr := json.Unmarshal(responseBody, &createdResource); parseErr != nil {
			return nil, fmt.Errorf("failed to parse practitioner response: %w", parseErr)
		}

		fhirID, _ := createdResource["id"].(string)
		if fhirID != "" {
			employee.FHIRResourceID = &fhirID
			if updateErr := staffService.repo.UpdateEmployeeFHIRResourceID(ctx, employee.ID, fhirID); updateErr != nil {
				return nil, fmt.Errorf("failed to update employee fhir resource id: %w", updateErr)
			}
		}
	}

	return employee, nil
}

func (staffService *service) GetEmployee(ctx context.Context, employeeID uuid.UUID) (*Employee, error) {
	employee, err := staffService.repo.GetEmployeeByID(ctx, employeeID)
	if err != nil {
		return nil, ErrEmployeeNotFound
	}
	return employee, nil
}

func (staffService *service) ListEmployees(ctx context.Context, search string, role string) ([]*Employee, error) {
	return staffService.repo.ListEmployees(ctx, search, role)
}

func (staffService *service) DeactivateEmployee(ctx context.Context, employeeID uuid.UUID) error {
	_, err := staffService.repo.GetEmployeeByID(ctx, employeeID)
	if err != nil {
		return ErrEmployeeNotFound
	}
	return staffService.repo.DeactivateEmployee(ctx, employeeID)
}
