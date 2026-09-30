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
	department := mockRepository.SeedDepartment("Cardiologia")
	staffService := staff.NewService(mockRepository, nil)

	input := staff.CreateEmployeeInput{
		CreatedBy:    uuid.New().String(),
		FullName:     "Dr. João Silva",
		Email:        "joao@clinic.com",
		Role:         string(role.RoleDoctor),
		CRMNumber:    "CRM-12345",
		DepartmentID: department.ID.String(),
	}

	employee, err := staffService.CreateEmployee(context.Background(), input)

	require.NoError(testingInstance, err)
	assert.NotNil(testingInstance, employee)
	assert.Equal(testingInstance, "Dr. João Silva", employee.FullName)
	assert.Equal(testingInstance, role.RoleDoctor, employee.Role)
	assert.Equal(testingInstance, mockRepository.RoleCatalog[string(role.RoleDoctor)].ID, employee.RoleID)
	require.NotNil(testingInstance, employee.CRMNumber)
	assert.Equal(testingInstance, "CRM-12345", *employee.CRMNumber)
	assert.Equal(testingInstance, department.ID, employee.DepartmentID)
	assert.Equal(testingInstance, "Cardiologia", employee.DepartmentName)
	assert.True(testingInstance, employee.IsActive)
}

func TestCreateEmployee_MissingFields_ReturnsErrorAndDoesNotCallRepository(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	department := mockRepository.SeedDepartment("Clínica Geral")
	departmentID := department.ID.String()

	testCases := []struct {
		name             string
		input            staff.CreateEmployeeInput
		expectedFieldKey string
	}{
		{
			name: "invalid created by",
			input: staff.CreateEmployeeInput{
				CreatedBy:    "not-a-uuid",
				FullName:     "Dr. João Silva",
				Email:        "joao@clinic.com",
				Role:         string(role.RoleDoctor),
				CRMNumber:    "CRM-12345",
				DepartmentID: departmentID,
			},
			expectedFieldKey: "created_by",
		},
		{
			name: "missing full name",
			input: staff.CreateEmployeeInput{
				CreatedBy:    uuid.New().String(),
				Email:        "joao@clinic.com",
				Role:         string(role.RoleDoctor),
				CRMNumber:    "CRM-12345",
				DepartmentID: departmentID,
			},
			expectedFieldKey: "full_name",
		},
		{
			name: "invalid email",
			input: staff.CreateEmployeeInput{
				CreatedBy:    uuid.New().String(),
				FullName:     "Dr. João Silva",
				Email:        "not-an-email",
				Role:         string(role.RoleDoctor),
				CRMNumber:    "CRM-12345",
				DepartmentID: departmentID,
			},
			expectedFieldKey: "email",
		},
		{
			name: "invalid role",
			input: staff.CreateEmployeeInput{
				CreatedBy:    uuid.New().String(),
				FullName:     "Dr. João Silva",
				Email:        "joao@clinic.com",
				Role:         "invalid-role",
				CRMNumber:    "CRM-12345",
				DepartmentID: departmentID,
			},
			expectedFieldKey: "role",
		},
		{
			name: "missing crm number",
			input: staff.CreateEmployeeInput{
				CreatedBy:    uuid.New().String(),
				FullName:     "Dr. João Silva",
				Email:        "joao@clinic.com",
				Role:         string(role.RoleDoctor),
				DepartmentID: departmentID,
			},
			expectedFieldKey: "crm_number",
		},
		{
			name: "invalid crm number",
			input: staff.CreateEmployeeInput{
				CreatedBy:    uuid.New().String(),
				FullName:     "Dr. João Silva",
				Email:        "joao@clinic.com",
				Role:         string(role.RoleDoctor),
				CRMNumber:    "invalid-crm",
				DepartmentID: departmentID,
			},
			expectedFieldKey: "crm_number",
		},
		{
			name: "missing department",
			input: staff.CreateEmployeeInput{
				CreatedBy: uuid.New().String(),
				FullName:  "Dr. João Silva",
				Email:     "joao@clinic.com",
				Role:      string(role.RoleDoctor),
				CRMNumber: "CRM-12345",
			},
			expectedFieldKey: "department_id",
		},
		{
			name: "invalid department",
			input: staff.CreateEmployeeInput{
				CreatedBy:    uuid.New().String(),
				FullName:     "Dr. João Silva",
				Email:        "joao@clinic.com",
				Role:         string(role.RoleDoctor),
				CRMNumber:    "CRM-12345",
				DepartmentID: "not-a-uuid",
			},
			expectedFieldKey: "department_id",
		},
	}

	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(subTest *testing.T) {
			isolatedRepository := mocks.NewMockStaffRepository()
			isolatedRepository.SeedDepartment("Clínica Geral")
			staffService := staff.NewService(isolatedRepository, nil)

			result, err := staffService.CreateEmployee(context.Background(), testCase.input)

			require.Error(subTest, err)
			assert.Nil(subTest, result)
			var appError apperrors.AppError
			require.True(subTest, errors.As(err, &appError))
			assert.Contains(subTest, appError.Message, testCase.expectedFieldKey)
		})
	}
}

func TestCreateEmployee_UnknownDepartment_ReturnsError(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	staffService := staff.NewService(mockRepository, nil)

	result, err := staffService.CreateEmployee(context.Background(), staff.CreateEmployeeInput{
		CreatedBy:    uuid.New().String(),
		FullName:     "Dr. João Silva",
		Email:        "joao@clinic.com",
		Role:         string(role.RoleDoctor),
		CRMNumber:    "CRM-12345",
		DepartmentID: uuid.New().String(),
	})

	assert.Nil(testingInstance, result)
	var appError apperrors.AppError
	require.True(testingInstance, errors.As(err, &appError))
	assert.Contains(testingInstance, appError.Message, "department_id")
}

func TestCreateEmployee_RoleMissingFromCatalog_ReturnsError(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	delete(mockRepository.RoleCatalog, string(role.RoleDoctor))
	department := mockRepository.SeedDepartment("Cardiologia")
	staffService := staff.NewService(mockRepository, nil)

	result, err := staffService.CreateEmployee(context.Background(), staff.CreateEmployeeInput{
		CreatedBy:    uuid.New().String(),
		FullName:     "Dr. João Silva",
		Email:        "joao@clinic.com",
		Role:         string(role.RoleDoctor),
		CRMNumber:    "CRM-12345",
		DepartmentID: department.ID.String(),
	})

	assert.Nil(testingInstance, result)
	var appError apperrors.AppError
	require.True(testingInstance, errors.As(err, &appError))
	assert.Contains(testingInstance, appError.Message, "role")
}

func TestCreateEmployee_RepositoryFailure_ReturnsError(testingInstance *testing.T) {
	expectedErr := errors.New("database unavailable")
	mockRepository := mocks.NewMockStaffRepository()
	department := mockRepository.SeedDepartment("Pediatria")
	mockRepository.Err = expectedErr
	staffService := staff.NewService(mockRepository, nil)

	input := staff.CreateEmployeeInput{
		CreatedBy:    uuid.New().String(),
		FullName:     "Dr. João Silva",
		Email:        "joao@clinic.com",
		Role:         string(role.RoleDoctor),
		CRMNumber:    "CRM-12345",
		DepartmentID: department.ID.String(),
	}

	result, err := staffService.CreateEmployee(context.Background(), input)

	assert.Nil(testingInstance, result)
	assert.ErrorIs(testingInstance, err, expectedErr)
}

func TestGetEmployee_ReturnsEmployee(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	department := mockRepository.SeedDepartment("Neurologia")
	staffService := staff.NewService(mockRepository, nil)

	createdEmployee, _ := staffService.CreateEmployee(context.Background(), staff.CreateEmployeeInput{
		CreatedBy:    uuid.New().String(),
		FullName:     "Enf. Maria Costa",
		Email:        "maria@clinic.com",
		Role:         string(role.RoleNurse),
		CRMNumber:    "COREN-12345",
		DepartmentID: department.ID.String(),
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

func TestSetEmployeeActive_TogglesStatusBothWays(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	department := mockRepository.SeedDepartment("Ortopedia")
	staffService := staff.NewService(mockRepository, nil)
	contextParam := context.Background()

	createdEmployee, _ := staffService.CreateEmployee(contextParam, staff.CreateEmployeeInput{
		CreatedBy:    uuid.New().String(),
		FullName:     "Recep. Ana Lima",
		Email:        "ana@clinic.com",
		Role:         string(role.RoleReception),
		CRMNumber:    "CRM-00001",
		DepartmentID: department.ID.String(),
	})

	disabledEmployee, errDisable := staffService.SetEmployeeActive(contextParam, staff.SetEmployeeActiveInput{
		EmployeeID: createdEmployee.ID.String(),
		IsActive:   false,
	})
	assert.NoError(testingInstance, errDisable)
	assert.False(testingInstance, disabledEmployee.IsActive)

	enabledEmployee, errEnable := staffService.SetEmployeeActive(contextParam, staff.SetEmployeeActiveInput{
		EmployeeID: createdEmployee.ID.String(),
		IsActive:   true,
	})
	assert.NoError(testingInstance, errEnable)
	assert.True(testingInstance, enabledEmployee.IsActive)
}

func TestSetEmployeeActive_InvalidIdentifier_ReturnsAppError(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	staffService := staff.NewService(mockRepository, nil)

	employee, errInvalid := staffService.SetEmployeeActive(context.Background(), staff.SetEmployeeActiveInput{
		EmployeeID: "not-a-uuid",
		IsActive:   false,
	})

	assert.Nil(testingInstance, employee)
	var appError apperrors.AppError
	require.True(testingInstance, errors.As(errInvalid, &appError))
	assert.Contains(testingInstance, appError.Message, "employee_id")
}

func TestSetEmployeeActive_NotFound_ReturnsAppError(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	staffService := staff.NewService(mockRepository, nil)

	employee, errNotFound := staffService.SetEmployeeActive(context.Background(), staff.SetEmployeeActiveInput{
		EmployeeID: uuid.New().String(),
		IsActive:   false,
	})

	assert.Nil(testingInstance, employee)
	var appError apperrors.AppError
	require.True(testingInstance, errors.As(errNotFound, &appError))
	assert.Equal(testingInstance, "employee not found", appError.Message)
}

func TestListDepartments_ReturnsCatalog(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	mockRepository.SeedDepartment("Clínica Geral")
	mockRepository.SeedDepartment("Cardiologia")
	staffService := staff.NewService(mockRepository, nil)

	departments, err := staffService.ListDepartments(context.Background())

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, departments, 2)
}

func TestListDepartments_RepositoryFailure_ReturnsError(testingInstance *testing.T) {
	expectedErr := errors.New("database unavailable")
	mockRepository := mocks.NewMockStaffRepository()
	mockRepository.Err = expectedErr
	staffService := staff.NewService(mockRepository, nil)

	departments, err := staffService.ListDepartments(context.Background())

	assert.Nil(testingInstance, departments)
	assert.ErrorIs(testingInstance, err, expectedErr)
}

func TestGetDepartment_ReturnsDepartment(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	department := mockRepository.SeedDepartment("Dermatologia")
	staffService := staff.NewService(mockRepository, nil)

	foundDepartment, err := staffService.GetDepartment(context.Background(), department.ID)

	assert.NoError(testingInstance, err)
	assert.Equal(testingInstance, "Dermatologia", foundDepartment.Name)
}

func TestGetDepartment_NotFound_ReturnsAppError(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	staffService := staff.NewService(mockRepository, nil)

	_, errNotFound := staffService.GetDepartment(context.Background(), uuid.New())

	var appError apperrors.AppError
	require.True(testingInstance, errors.As(errNotFound, &appError))
	assert.Equal(testingInstance, "department not found", appError.Message)
}

func TestGetDefaultDepartmentID_ReturnsFirstActiveDepartmentBySortOrder(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	neurology := mockRepository.SeedDepartment("Neurologia")
	generalPractice := mockRepository.SeedDepartment("Clínica Geral")
	cardiology := mockRepository.SeedDepartment("Cardiologia")
	mockRepository.Departments[neurology.ID].IsActive = false
	staffService := staff.NewService(mockRepository, nil)

	departmentID, err := staffService.GetDefaultDepartmentID(context.Background())

	require.NoError(testingInstance, err)
	assert.Equal(testingInstance, generalPractice.ID, departmentID)
	assert.NotEqual(testingInstance, cardiology.ID, departmentID)
}

func TestGetDefaultDepartmentID_NoActiveDepartment_ReturnsAppError(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	department := mockRepository.SeedDepartment("Clínica Geral")
	mockRepository.Departments[department.ID].IsActive = false
	staffService := staff.NewService(mockRepository, nil)

	result, errNotFound := staffService.GetDefaultDepartmentID(context.Background())

	assert.Equal(testingInstance, uuid.Nil, result)
	var appError apperrors.AppError
	require.True(testingInstance, errors.As(errNotFound, &appError))
	assert.Equal(testingInstance, "department not found", appError.Message)
}

func TestGetDefaultDepartmentID_RepositoryFailure_ReturnsError(testingInstance *testing.T) {
	expectedErr := errors.New("database unavailable")
	mockRepository := mocks.NewMockStaffRepository()
	mockRepository.Err = expectedErr
	staffService := staff.NewService(mockRepository, nil)

	result, err := staffService.GetDefaultDepartmentID(context.Background())

	assert.Equal(testingInstance, uuid.Nil, result)
	assert.ErrorIs(testingInstance, err, expectedErr)
}

func TestListEmployees_IncludesInactiveEmployees(testingInstance *testing.T) {
	mockRepository := mocks.NewMockStaffRepository()
	department := mockRepository.SeedDepartment("Urologia")
	staffService := staff.NewService(mockRepository, nil)
	contextParam := context.Background()

	firstEmployee, _ := staffService.CreateEmployee(contextParam, staff.CreateEmployeeInput{
		CreatedBy:    uuid.New().String(),
		FullName:     "Dr. A",
		Email:        "a@clinic.com",
		Role:         string(role.RoleDoctor),
		CRMNumber:    "CRM-1",
		DepartmentID: department.ID.String(),
	})
	secondEmployee, _ := staffService.CreateEmployee(contextParam, staff.CreateEmployeeInput{
		CreatedBy:    uuid.New().String(),
		FullName:     "Dr. B",
		Email:        "b@clinic.com",
		Role:         string(role.RoleDoctor),
		CRMNumber:    "CRM-2",
		DepartmentID: department.ID.String(),
	})

	_, disableErr := staffService.SetEmployeeActive(contextParam, staff.SetEmployeeActiveInput{
		EmployeeID: secondEmployee.ID.String(),
		IsActive:   false,
	})
	assert.NoError(testingInstance, disableErr)

	employees, err := staffService.ListEmployees(contextParam, "", "")

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, employees, 2)
	assert.Equal(testingInstance, firstEmployee.ID, employees[0].ID)
}
