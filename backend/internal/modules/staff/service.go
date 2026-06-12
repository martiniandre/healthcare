package staff

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/auth"
	"github.com/healthcare/backend/internal/shared/validator"
)

var ErrEmployeeNotFound = errors.New("employee not found")

type Service interface {
	CreateEmployee(ctx context.Context, userID uuid.UUID, fullName, email, role, crmNumber string) (*Employee, error)
	GetEmployee(ctx context.Context, employeeID uuid.UUID) (*Employee, error)
	ListEmployees(ctx context.Context, search string, role string) ([]*Employee, error)
	DeactivateEmployee(ctx context.Context, employeeID uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (staffService *service) CreateEmployee(ctx context.Context, userID uuid.UUID, fullName, email, role, crmNumber string) (*Employee, error) {
	parsedRole, roleIsValid := auth.ParseRole(role)
	if !roleIsValid {
		return nil, auth.ErrInvalidRole
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

	employee := &Employee{
		ID:        uuid.New(),
		UserID:    userID,
		FullName:  fullName,
		Email:     email,
		Role:      parsedRole,
		CRMNumber: crmNumberPtr,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := staffService.repo.CreateEmployee(ctx, employee)
	if err != nil {
		return nil, err
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
