package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/staff"
	"github.com/healthcare/backend/internal/modules/staff/mocks"
	"github.com/healthcare/backend/internal/modules/staff/pb"
	"github.com/healthcare/backend/internal/shared/role"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGRPCCreateEmployee_AssignsDefaultDepartment(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	defaultDepartment := mockRepository.SeedDepartment("Clínica Geral")
	grpcHandler := staff.NewGRPCHandler(staff.NewService(mockRepository, nil))

	response, err := grpcHandler.CreateEmployee(context.Background(), &pb.CreateEmployeeRequest{
		CreatedBy: uuid.New().String(),
		FullName:  "Dr. Paulo Nunes",
		Email:     "paulo@clinica.com",
		Role:      string(role.RoleDoctor),
		CrmNumber: "CRM-98765",
	})

	require.NoError(testingInstance, err)
	assert.NotEmpty(testingInstance, response.EmployeeId)

	createdEmployee, findErr := mockRepository.GetEmployeeByID(context.Background(), uuid.MustParse(response.EmployeeId))
	require.NoError(testingInstance, findErr)
	assert.Equal(testingInstance, defaultDepartment.ID, createdEmployee.DepartmentID)
	assert.Equal(testingInstance, "Clínica Geral", createdEmployee.DepartmentName)
}

func TestGRPCCreateEmployee_WithoutDepartments_ReturnsError(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	grpcHandler := staff.NewGRPCHandler(staff.NewService(mockRepository, nil))

	response, err := grpcHandler.CreateEmployee(context.Background(), &pb.CreateEmployeeRequest{
		CreatedBy: uuid.New().String(),
		FullName:  "Dr. Paulo Nunes",
		Email:     "paulo@clinica.com",
		Role:      string(role.RoleDoctor),
		CrmNumber: "CRM-98765",
	})

	assert.Nil(testingInstance, response)
	assert.Error(testingInstance, err)
	assert.Empty(testingInstance, mockRepository.Employees)
}

func TestGRPCDeactivateEmployee_SetsEmployeeInactive(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	department := mockRepository.SeedDepartment("Pediatria")
	grpcHandler := staff.NewGRPCHandler(staff.NewService(mockRepository, nil))
	contextParam := context.Background()

	createdEmployee, err := grpcHandler.CreateEmployee(contextParam, &pb.CreateEmployeeRequest{
		CreatedBy: uuid.New().String(),
		FullName:  "Enf. Carla Dias",
		Email:     "carla@clinica.com",
		Role:      string(role.RoleNurse),
		CrmNumber: "COREN-55555",
	})
	require.NoError(testingInstance, err)
	require.Equal(testingInstance, department.ID, mockRepository.Employees[uuid.MustParse(createdEmployee.EmployeeId)].DepartmentID)

	_, deactivateErr := grpcHandler.DeactivateEmployee(contextParam, &pb.DeactivateEmployeeRequest{
		EmployeeId: createdEmployee.EmployeeId,
	})

	assert.NoError(testingInstance, deactivateErr)
	assert.False(testingInstance, mockRepository.Employees[uuid.MustParse(createdEmployee.EmployeeId)].IsActive)
}
