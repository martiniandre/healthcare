package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/patients"
	"github.com/healthcare/backend/internal/modules/patients/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPatientService_CreatePatient(testingInstance *testing.T) {
	mockRepository := mocks.NewMockPatientRepository()
	patientService := patients.NewService(mockRepository)
	contextParam := context.Background()

	patient, err := patientService.CreatePatient(contextParam, "Pedro Alves", "1990-05-20", "123.456.789-00", "+55 11 99999-0000", "A+", "Nenhuma")

	assert.NoError(testingInstance, err)
	assert.NotNil(testingInstance, patient)
	assert.Equal(testingInstance, "Pedro Alves", patient.FullName)
	assert.Equal(testingInstance, "123.456.789-00", patient.DocumentID)

	_, errDuplicate := patientService.CreatePatient(contextParam, "Pedro Alves Duplicado", "1990-05-20", "123.456.789-00", "", "", "")
	assert.ErrorIs(testingInstance, errDuplicate, patients.ErrPatientAlreadyExists)
}

func TestPatientService_CreatePatient_InvalidDate(testingInstance *testing.T) {
	mockRepository := mocks.NewMockPatientRepository()
	patientService := patients.NewService(mockRepository)
	contextParam := context.Background()

	_, err := patientService.CreatePatient(contextParam, "Nome Teste", "20/05/1990", "999.999.999-99", "", "", "")

	assert.Error(testingInstance, err)
	assert.Contains(testingInstance, err.Error(), "invalid birth date format")
}

func TestPatientService_GetPatient(testingInstance *testing.T) {
	mockRepository := mocks.NewMockPatientRepository()
	patientService := patients.NewService(mockRepository)
	contextParam := context.Background()

	created, _ := patientService.CreatePatient(contextParam, "Ana Souza", "1985-10-15", "987.654.321-00", "", "O-", "Penicilina")

	found, err := patientService.GetPatient(contextParam, created.ID)
	assert.NoError(testingInstance, err)
	assert.Equal(testingInstance, created.ID, found.ID)

	_, errNotFound := patientService.GetPatient(contextParam, uuid.New())
	assert.ErrorIs(testingInstance, errNotFound, patients.ErrPatientNotFound)
}

func TestPatientService_GetPatientByDocument(testingInstance *testing.T) {
	mockRepository := mocks.NewMockPatientRepository()
	patientService := patients.NewService(mockRepository)
	contextParam := context.Background()

	patientService.CreatePatient(contextParam, "Carlos Melo", "2000-01-01", "111.222.333-44", "", "B+", "")

	found, err := patientService.GetPatientByDocument(contextParam, "111.222.333-44")
	assert.NoError(testingInstance, err)
	assert.Equal(testingInstance, "Carlos Melo", found.FullName)
}
