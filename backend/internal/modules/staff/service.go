package staff

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/healthcare/backend/internal/shared/role"
	"github.com/healthcare/backend/internal/shared/validator"
)

type Service interface {
	CreateEmployee(ctx context.Context, input CreateEmployeeInput) (*Employee, error)
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

func (staffService *service) CreateEmployee(ctx context.Context, input CreateEmployeeInput) (*Employee, error) {
	if fieldViolations := validateEmployeeFields(input); len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid employee input", fieldViolations)
	}

	parsedRole, _ := role.ParseRole(input.Role)
	parsedCreatedBy, _ := uuid.Parse(input.CreatedBy)
	employee := &Employee{
		ID:        uuid.New(),
		FullName:  input.FullName,
		Email:     input.Email,
		Role:      parsedRole,
		CRMNumber: nil,
		CreatedBy: &parsedCreatedBy,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if input.CRMNumber != "" {
		employee.CRMNumber = &input.CRMNumber
	}

	err := staffService.repo.CreateEmployee(ctx, employee)
	if err != nil {
		return nil, err
	}

	if staffService.fhirClient != nil {
		practitionerResource := fhir.NewPractitionerResource(input.FullName, input.CRMNumber)
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
	return staffService.repo.GetEmployeeByID(ctx, employeeID)
}

func (staffService *service) ListEmployees(ctx context.Context, search string, role string) ([]*Employee, error) {
	return staffService.repo.ListEmployees(ctx, search, role)
}

func (staffService *service) DeactivateEmployee(ctx context.Context, employeeID uuid.UUID) error {
	_, err := staffService.repo.GetEmployeeByID(ctx, employeeID)
	if err != nil {
		return err
	}
	return staffService.repo.DeactivateEmployee(ctx, employeeID)
}

func validateEmployeeFields(input CreateEmployeeInput) map[string]string {
	fieldViolations := make(map[string]string)
	if strings.TrimSpace(input.CreatedBy) == "" {
		fieldViolations["created_by"] = "is required"
	} else if _, err := uuid.Parse(input.CreatedBy); err != nil {
		fieldViolations["created_by"] = "invalid UUID format"
	}
	if strings.TrimSpace(input.FullName) == "" {
		fieldViolations["full_name"] = "is required"
	}
	if strings.TrimSpace(input.Email) == "" || !validator.IsValidEmail(input.Email) {
		fieldViolations["email"] = "invalid email format"
	}
	if _, roleIsValid := role.ParseRole(input.Role); !roleIsValid {
		fieldViolations["role"] = "invalid role"
	}
	if input.CRMNumber != "" && !validator.IsValidCRMNumber(input.CRMNumber) {
		fieldViolations["crm_number"] = "invalid CRM format"
	}
	return fieldViolations
}
