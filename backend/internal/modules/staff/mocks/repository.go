package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/staff"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/role"
)

type MockStaffRepository struct {
	Employees   map[uuid.UUID]*staff.Employee
	Departments map[uuid.UUID]*staff.Department
	RoleCatalog map[string]*staff.Role
	Err         error
}

func NewMockStaffRepository() *MockStaffRepository {
	mockRepository := &MockStaffRepository{
		Employees:   make(map[uuid.UUID]*staff.Employee),
		Departments: make(map[uuid.UUID]*staff.Department),
		RoleCatalog: make(map[string]*staff.Role),
	}
	roleCodes := []role.Role{role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception, role.RolePatient}
	for roleIndex, roleCode := range roleCodes {
		seededRole := &staff.Role{
			ID:        uuid.New(),
			Code:      roleCode.String(),
			Name:      roleCode.String(),
			SortOrder: roleIndex + 1,
			IsActive:  true,
		}
		mockRepository.RoleCatalog[seededRole.Code] = seededRole
	}
	return mockRepository
}

func (mockRepo *MockStaffRepository) GetRoleIDByCode(contextParam context.Context, roleCode string) (uuid.UUID, error) {
	if mockRepo.Err != nil {
		return uuid.Nil, mockRepo.Err
	}
	catalogRole, roleExists := mockRepo.RoleCatalog[roleCode]
	if !roleExists {
		return uuid.Nil, apperrors.ErrRoleNotFound
	}
	return catalogRole.ID, nil
}

func (mockRepo *MockStaffRepository) SeedDepartment(departmentName string) *staff.Department {
	department := &staff.Department{
		ID:        uuid.New(),
		Name:      departmentName,
		SortOrder: len(mockRepo.Departments) + 1,
		IsActive:  true,
	}
	mockRepo.Departments[department.ID] = department
	return department
}

func (mockRepo *MockStaffRepository) CreateEmployee(contextParam context.Context, employee *staff.Employee) error {
	if mockRepo.Err != nil {
		return mockRepo.Err
	}
	department, departmentExists := mockRepo.Departments[employee.DepartmentID]
	if !departmentExists {
		return apperrors.InvalidArgument("invalid employee input", map[string]string{"department_id": "unknown department"})
	}
	catalogRole, roleExists := mockRepo.RoleCatalog[string(employee.Role)]
	if !roleExists {
		return apperrors.InvalidArgument("invalid employee input", map[string]string{"role": "unknown role"})
	}
	employee.RoleID = catalogRole.ID
	employee.DepartmentName = department.Name
	mockRepo.Employees[employee.ID] = employee
	return nil
}

func (mockRepo *MockStaffRepository) GetEmployeeByID(contextParam context.Context, employeeID uuid.UUID) (*staff.Employee, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}
	employee, exists := mockRepo.Employees[employeeID]
	if !exists {
		return nil, apperrors.ErrEmployeeNotFound
	}
	return employee, nil
}

func (mockRepo *MockStaffRepository) ListEmployees(contextParam context.Context, search string, role string) ([]*staff.Employee, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}
	result := make([]*staff.Employee, 0, len(mockRepo.Employees))
	for _, employee := range mockRepo.Employees {
		result = append(result, employee)
	}
	return result, nil
}

func (mockRepo *MockStaffRepository) SetEmployeeActive(contextParam context.Context, employeeID uuid.UUID, isActive bool) error {
	if mockRepo.Err != nil {
		return mockRepo.Err
	}
	employee, exists := mockRepo.Employees[employeeID]
	if !exists {
		return apperrors.ErrEmployeeNotFound
	}
	employee.IsActive = isActive
	return nil
}

func (mockRepo *MockStaffRepository) GetDepartmentByID(contextParam context.Context, departmentID uuid.UUID) (*staff.Department, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}
	department, exists := mockRepo.Departments[departmentID]
	if !exists {
		return nil, apperrors.ErrDepartmentNotFound
	}
	return department, nil
}

func (mockRepo *MockStaffRepository) GetDefaultDepartment(contextParam context.Context) (*staff.Department, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}
	defaultDepartment := &staff.Department{}
	foundAny := false
	for _, department := range mockRepo.Departments {
		if !department.IsActive {
			continue
		}
		if !foundAny || department.SortOrder < defaultDepartment.SortOrder {
			defaultDepartment = department
			foundAny = true
		}
	}
	if !foundAny {
		return nil, apperrors.ErrDepartmentNotFound
	}
	return defaultDepartment, nil
}

func (mockRepo *MockStaffRepository) ListDepartments(contextParam context.Context) ([]*staff.Department, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}
	departments := make([]*staff.Department, 0, len(mockRepo.Departments))
	for _, department := range mockRepo.Departments {
		departments = append(departments, department)
	}
	return departments, nil
}

func (mockRepo *MockStaffRepository) UpdateEmployeeFHIRResourceID(contextParam context.Context, employeeID uuid.UUID, fhirResourceID string) error {
	if mockRepo.Err != nil {
		return mockRepo.Err
	}
	employee, exists := mockRepo.Employees[employeeID]
	if !exists {
		return apperrors.ErrEmployeeNotFound
	}
	employee.FHIRResourceID = &fhirResourceID
	return nil
}
