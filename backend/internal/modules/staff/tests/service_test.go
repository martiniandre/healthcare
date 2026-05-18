package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/auth"
	"github.com/healthcare/backend/internal/modules/staff"
	"github.com/healthcare/backend/internal/modules/staff/mocks"
	"github.com/stretchr/testify/assert"
)

func TestStaffService_CreateEmployee(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	staffService := staff.NewService(mockRepository)
	contextParam := context.Background()
	userID := uuid.New()

	employee, err := staffService.CreateEmployee(contextParam, userID, "Dr. João Silva", "joao@clinic.com", string(auth.RoleDoctor), "CRM-12345")

	assert.NoError(testingInstance, err)
	assert.NotNil(testingInstance, employee)
	assert.Equal(testingInstance, "Dr. João Silva", employee.FullName)
	assert.Equal(testingInstance, auth.RoleDoctor, employee.Role)
	assert.Equal(testingInstance, "CRM-12345", employee.CRMNumber)
	assert.True(testingInstance, employee.IsActive)
}

func TestStaffService_GetEmployee(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	staffService := staff.NewService(mockRepository)
	contextParam := context.Background()
	userID := uuid.New()

	createdEmployee, _ := staffService.CreateEmployee(contextParam, userID, "Enf. Maria Costa", "maria@clinic.com", string(auth.RoleNurse), "")

	foundEmployee, err := staffService.GetEmployee(contextParam, createdEmployee.ID)

	assert.NoError(testingInstance, err)
	assert.Equal(testingInstance, createdEmployee.ID, foundEmployee.ID)

	_, errNotFound := staffService.GetEmployee(contextParam, uuid.New())
	assert.ErrorIs(testingInstance, errNotFound, staff.ErrEmployeeNotFound)
}

func TestStaffService_DeactivateEmployee(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	staffService := staff.NewService(mockRepository)
	contextParam := context.Background()
	userID := uuid.New()

	createdEmployee, _ := staffService.CreateEmployee(contextParam, userID, "Recep. Ana Lima", "ana@clinic.com", string(auth.RoleReception), "")

	err := staffService.DeactivateEmployee(contextParam, createdEmployee.ID)
	assert.NoError(testingInstance, err)

	errNotFound := staffService.DeactivateEmployee(contextParam, uuid.New())
	assert.ErrorIs(testingInstance, errNotFound, staff.ErrEmployeeNotFound)
}

func TestStaffService_ListEmployees(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	staffService := staff.NewService(mockRepository)
	contextParam := context.Background()

	staffService.CreateEmployee(contextParam, uuid.New(), "Dr. A", "a@clinic.com", string(auth.RoleDoctor), "CRM-1")
	staffService.CreateEmployee(contextParam, uuid.New(), "Dr. B", "b@clinic.com", string(auth.RoleDoctor), "CRM-2")

	employees, err := staffService.ListEmployees(contextParam)

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, employees, 2)
}
