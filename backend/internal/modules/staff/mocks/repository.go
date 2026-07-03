package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/staff"
)

type MockStaffRepository struct {
	Employees map[uuid.UUID]*staff.Employee
	Err       error
}

func NewMockStaffRepository() *MockStaffRepository {
	return &MockStaffRepository{
		Employees: make(map[uuid.UUID]*staff.Employee),
	}
}

func (mockRepo *MockStaffRepository) CreateEmployee(contextParam context.Context, employee *staff.Employee) error {
	if mockRepo.Err != nil {
		return mockRepo.Err
	}
	mockRepo.Employees[employee.ID] = employee
	return nil
}

func (mockRepo *MockStaffRepository) GetEmployeeByID(contextParam context.Context, employeeID uuid.UUID) (*staff.Employee, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}
	employee, exists := mockRepo.Employees[employeeID]
	if !exists {
		return nil, staff.ErrEmployeeNotFound
	}
	return employee, nil
}

func (mockRepo *MockStaffRepository) GetEmployeeByUserID(contextParam context.Context, userID uuid.UUID) (*staff.Employee, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}
	for _, employee := range mockRepo.Employees {
		if employee.UserID == userID {
			return employee, nil
		}
	}
	return nil, staff.ErrEmployeeNotFound
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

func (mockRepo *MockStaffRepository) DeactivateEmployee(contextParam context.Context, employeeID uuid.UUID) error {
	if mockRepo.Err != nil {
		return mockRepo.Err
	}
	employee, exists := mockRepo.Employees[employeeID]
	if !exists {
		return staff.ErrEmployeeNotFound
	}
	employee.IsActive = false
	return nil
}
