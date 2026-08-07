package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/patients"
	"github.com/healthcare/backend/internal/modules/patients/mocks"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPatientEventBus struct {
	PublishedEvents []eventbus.Event
}

func (mockBus *mockPatientEventBus) Publish(ctx context.Context, event eventbus.Event) error {
	mockBus.PublishedEvents = append(mockBus.PublishedEvents, event)
	return nil
}

func (mockBus *mockPatientEventBus) Subscribe(eventName string, handler eventbus.Handler) {}

func TestCreatePatient_ValidInput_PublishesEvent(testingInstance *testing.T) {
	eventBus := &mockPatientEventBus{}
	mockRepository := mocks.NewMockPatientRepository()
	patientService := patients.NewService(mockRepository, eventBus)
	contextParam := context.Background()

	input := patients.CreatePatientInput{
		FullName:    "Maria Oliveira",
		BirthDate:   "1988-03-15",
		DocumentID:  "111.444.777-35",
		PhoneNumber: "(21) 98888-0000",
	}

	patient, creationError := patientService.CreatePatient(contextParam, input)

	assert.NoError(testingInstance, creationError)
	assert.NotNil(testingInstance, patient)
	assert.Len(testingInstance, eventBus.PublishedEvents, 1)
	assert.Equal(testingInstance, "patient.created", eventBus.PublishedEvents[0].Name)
	assert.Equal(testingInstance, "Novo Paciente Cadastrado", eventBus.PublishedEvents[0].Data["title"])
	assert.Equal(testingInstance, "patient", eventBus.PublishedEvents[0].Data["resource_type"])
	assert.Contains(testingInstance, eventBus.PublishedEvents[0].Data["body"], "Maria Oliveira")
}

func TestCreatePatient_ValidInput_ReturnsCreatedPatient(testingInstance *testing.T) {
	mockRepository := mocks.NewMockPatientRepository()
	patientService := patients.NewService(mockRepository, nil)
	contextParam := context.Background()

	input := patients.CreatePatientInput{
		FullName:    "Pedro Alves",
		BirthDate:   "1990-05-20",
		DocumentID:  "123.456.789-09",
		PhoneNumber: "(11) 99999-0000",
	}

	patient, creationError := patientService.CreatePatient(contextParam, input)

	assert.NoError(testingInstance, creationError)
	assert.NotNil(testingInstance, patient)
	assert.Equal(testingInstance, "Pedro Alves", patient.FullName)
	assert.Equal(testingInstance, "123.456.789-09", patient.DocumentID)
}

func TestCreatePatient_DuplicateDocument_ReturnsAlreadyExists(testingInstance *testing.T) {
	mockRepository := mocks.NewMockPatientRepository()
	patientService := patients.NewService(mockRepository, nil)
	contextParam := context.Background()

	firstInput := patients.CreatePatientInput{
		FullName:    "Pedro Alves",
		BirthDate:   "1990-05-20",
		DocumentID:  "123.456.789-09",
		PhoneNumber: "(11) 99999-0000",
	}
	duplicateInput := patients.CreatePatientInput{
		FullName:    "Pedro Alves Duplicado",
		BirthDate:   "1990-05-20",
		DocumentID:  "123.456.789-09",
		PhoneNumber: "(11) 99999-0000",
	}

	_, firstError := patientService.CreatePatient(contextParam, firstInput)
	require.NoError(testingInstance, firstError)

	_, errDuplicate := patientService.CreatePatient(contextParam, duplicateInput)

	var appError apperrors.AppError
	require.True(testingInstance, errors.As(errDuplicate, &appError))
	assert.Equal(testingInstance, "patient already exists", appError.Message)
}

func TestCreatePatient_InvalidInput_ReturnsErrorAndDoesNotCallRepository(testingInstance *testing.T) {
	testCases := []struct {
		name             string
		input            patients.CreatePatientInput
		expectedFieldKey string
	}{
		{
			name: "missing full name",
			input: patients.CreatePatientInput{
				BirthDate:   "1990-05-20",
				DocumentID:  "123.456.789-09",
				PhoneNumber: "(11) 99999-0000",
			},
			expectedFieldKey: "full_name",
		},
		{
			name: "invalid birth date format",
			input: patients.CreatePatientInput{
				FullName:    "Pedro Alves",
				BirthDate:   "20/05/1990",
				DocumentID:  "123.456.789-09",
				PhoneNumber: "(11) 99999-0000",
			},
			expectedFieldKey: "birth_date",
		},
		{
			name: "birth date in the future",
			input: patients.CreatePatientInput{
				FullName:    "Pedro Alves",
				BirthDate:   "2999-05-20",
				DocumentID:  "123.456.789-09",
				PhoneNumber: "(11) 99999-0000",
			},
			expectedFieldKey: "birth_date",
		},
		{
			name: "invalid cpf",
			input: patients.CreatePatientInput{
				FullName:    "Pedro Alves",
				BirthDate:   "1990-05-20",
				DocumentID:  "invalid-cpf",
				PhoneNumber: "(11) 99999-0000",
			},
			expectedFieldKey: "document_id",
		},
		{
			name: "invalid phone",
			input: patients.CreatePatientInput{
				FullName:    "Pedro Alves",
				BirthDate:   "1990-05-20",
				DocumentID:  "123.456.789-09",
				PhoneNumber: "not-a-phone",
			},
			expectedFieldKey: "phone_number",
		},
	}

	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(subTest *testing.T) {
			mockRepository := mocks.NewMockPatientRepository()
			patientService := patients.NewService(mockRepository, nil)

			result, err := patientService.CreatePatient(context.Background(), testCase.input)

			require.Error(subTest, err)
			assert.Nil(subTest, result)
			var appError apperrors.AppError
			require.True(subTest, errors.As(err, &appError))
			assert.Contains(subTest, appError.Message, testCase.expectedFieldKey)
		})
	}
}

func TestGetPatient_ReturnsPatient(testingInstance *testing.T) {
	mockRepository := mocks.NewMockPatientRepository()
	patientService := patients.NewService(mockRepository, nil)
	contextParam := context.Background()

	created, _ := patientService.CreatePatient(contextParam, patients.CreatePatientInput{
		FullName:    "Ana Souza",
		BirthDate:   "1985-10-15",
		DocumentID:  "987.654.321-00",
		PhoneNumber: "(21) 98888-0000",
	})

	found, getError := patientService.GetPatient(contextParam, created.FHIRResourceID)
	assert.NoError(testingInstance, getError)
	assert.Equal(testingInstance, created.FHIRResourceID, found.FHIRResourceID)
}

func TestGetPatient_NotFound_ReturnsAppError(testingInstance *testing.T) {
	mockRepository := mocks.NewMockPatientRepository()
	patientService := patients.NewService(mockRepository, nil)

	_, errNotFound := patientService.GetPatient(context.Background(), uuid.New().String())

	var appError apperrors.AppError
	require.True(testingInstance, errors.As(errNotFound, &appError))
	assert.Equal(testingInstance, "patient not found", appError.Message)
}

func TestGetPatientByDocument_ReturnsPatient(testingInstance *testing.T) {
	mockRepository := mocks.NewMockPatientRepository()
	patientService := patients.NewService(mockRepository, nil)
	contextParam := context.Background()

	patientService.CreatePatient(contextParam, patients.CreatePatientInput{
		FullName:    "Carlos Melo",
		BirthDate:   "2000-01-01",
		DocumentID:  "529.982.247-25",
		PhoneNumber: "(11) 99999-0000",
	})

	found, getError := patientService.GetPatientByDocument(contextParam, "529.982.247-25")
	assert.NoError(testingInstance, getError)
	assert.Equal(testingInstance, "Carlos Melo", found.FullName)
}
