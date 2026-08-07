package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/staff"
	"github.com/healthcare/backend/internal/modules/staff/mocks"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/role"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEmployee_ValidInput_CreatesEmployee(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	staffService := staff.NewService(mockRepository, nil)

	input := staff.CreateEmployeeInput{
		CreatedBy: uuid.New().String(),
		FullName:  "Dr. João Silva",
		Email:     "joao@clinic.com",
		Role:      string(role.RoleDoctor),
		CRMNumber: "CRM-12345",
	}

	employee, err := staffService.CreateEmployee(context.Background(), input)

	require.NoError(testingInstance, err)
	assert.NotNil(testingInstance, employee)
	assert.Equal(testingInstance, "Dr. João Silva", employee.FullName)
	assert.Equal(testingInstance, role.RoleDoctor, employee.Role)
	require.NotNil(testingInstance, employee.CRMNumber)
	assert.Equal(testingInstance, "CRM-12345", *employee.CRMNumber)
	assert.True(testingInstance, employee.IsActive)
}

func TestCreateEmployee_MissingFields_ReturnsErrorAndDoesNotCallRepository(testingInstance *testing.T) {
	testCases := []struct {
		name             string
		input            staff.CreateEmployeeInput
		expectedFieldKey string
	}{
		{
			name: "invalid created by",
			input: staff.CreateEmployeeInput{
				CreatedBy: "not-a-uuid",
				FullName:  "Dr. João Silva",
				Email:     "joao@clinic.com",
				Role:      string(role.RoleDoctor),
			},
			expectedFieldKey: "created_by",
		},
		{
			name: "missing full name",
			input: staff.CreateEmployeeInput{
				CreatedBy: uuid.New().String(),
				Email:     "joao@clinic.com",
				Role:      string(role.RoleDoctor),
			},
			expectedFieldKey: "full_name",
		},
		{
			name: "invalid email",
			input: staff.CreateEmployeeInput{
				CreatedBy: uuid.New().String(),
				FullName:  "Dr. João Silva",
				Email:     "not-an-email",
				Role:      string(role.RoleDoctor),
			},
			expectedFieldKey: "email",
		},
		{
			name: "invalid role",
			input: staff.CreateEmployeeInput{
				CreatedBy: uuid.New().String(),
				FullName:  "Dr. João Silva",
				Email:     "joao@clinic.com",
				Role:      "invalid-role",
			},
			expectedFieldKey: "role",
		},
		{
			name: "invalid crm number",
			input: staff.CreateEmployeeInput{
				CreatedBy: uuid.New().String(),
				FullName:  "Dr. João Silva",
				Email:     "joao@clinic.com",
				Role:      string(role.RoleDoctor),
				CRMNumber: "invalid-crm",
			},
			expectedFieldKey: "crm_number",
		},
	}

	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(subTest *testing.T) {
			mockRepository := mocks.NewMockStaffRepository()
			staffService := staff.NewService(mockRepository, nil)

			result, err := staffService.CreateEmployee(context.Background(), testCase.input)

			require.Error(subTest, err)
			assert.Nil(subTest, result)
			var appError apperrors.AppError
			require.True(subTest, errors.As(err, &appError))
			assert.Contains(subTest, appError.Message, testCase.expectedFieldKey)
		})
	}
}

func TestCreateEmployee_RepositoryFailure_ReturnsError(testingInstance *testing.T) {
	expectedErr := errors.New("database unavailable")
	mockRepository := mocks.NewMockStaffRepository()
	mockRepository.Err = expectedErr
	staffService := staff.NewService(mockRepository, nil)

	input := staff.CreateEmployeeInput{
		CreatedBy: uuid.New().String(),
		FullName:  "Dr. João Silva",
		Email:     "joao@clinic.com",
		Role:      string(role.RoleDoctor),
	}

	result, err := staffService.CreateEmployee(context.Background(), input)

	assert.Nil(testingInstance, result)
	assert.ErrorIs(testingInstance, err, expectedErr)
}

func TestGetEmployee_ReturnsEmployee(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	staffService := staff.NewService(mockRepository, nil)

	createdEmployee, _ := staffService.CreateEmployee(context.Background(), staff.CreateEmployeeInput{
		CreatedBy: uuid.New().String(),
		FullName:  "Enf. Maria Costa",
		Email:     "maria@clinic.com",
		Role:      string(role.RoleNurse),
	})

	foundEmployee, err := staffService.GetEmployee(context.Background(), createdEmployee.ID)

	assert.NoError(testingInstance, err)
	assert.Equal(testingInstance, createdEmployee.ID, foundEmployee.ID)
}

func TestGetEmployee_NotFound_ReturnsAppError(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	staffService := staff.NewService(mockRepository, nil)

	_, errNotFound := staffService.GetEmployee(context.Background(), uuid.New())

	var appError apperrors.AppError
	require.True(testingInstance, errors.As(errNotFound, &appError))
	assert.Equal(testingInstance, "employee not found", appError.Message)
}

func TestDeactivateEmployee(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	staffService := staff.NewService(mockRepository, nil)

	createdEmployee, _ := staffService.CreateEmployee(context.Background(), staff.CreateEmployeeInput{
		CreatedBy: uuid.New().String(),
		FullName:  "Recep. Ana Lima",
		Email:     "ana@clinic.com",
		Role:      string(role.RoleReception),
	})

	err := staffService.DeactivateEmployee(context.Background(), createdEmployee.ID)
	assert.NoError(testingInstance, err)
}

func TestDeactivateEmployee_NotFound_ReturnsAppError(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	staffService := staff.NewService(mockRepository, nil)

	errNotFound := staffService.DeactivateEmployee(context.Background(), uuid.New())

	var appError apperrors.AppError
	require.True(testingInstance, errors.As(errNotFound, &appError))
	assert.Equal(testingInstance, "employee not found", appError.Message)
}

func TestListEmployees_ReturnsActiveEmployees(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	staffService := staff.NewService(mockRepository, nil)
	contextParam := context.Background()

	staffService.CreateEmployee(contextParam, staff.CreateEmployeeInput{
		CreatedBy: uuid.New().String(),
		FullName:  "Dr. A",
		Email:     "a@clinic.com",
		Role:      string(role.RoleDoctor),
		CRMNumber: "CRM-1",
	})
	staffService.CreateEmployee(contextParam, staff.CreateEmployeeInput{
		CreatedBy: uuid.New().String(),
		FullName:  "Dr. B",
		Email:     "b@clinic.com",
		Role:      string(role.RoleDoctor),
		CRMNumber: "CRM-2",
	})

	employees, err := staffService.ListEmployees(contextParam, "", "")

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, employees, 2)
}
