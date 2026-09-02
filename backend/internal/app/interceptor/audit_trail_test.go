package interceptor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsClinicalOrCriticalMethod_ExcludesAuthLogin(testingInstance *testing.T) {
	assert.False(testingInstance, isClinicalOrCriticalMethod("/auth.v1.AuthService/Login"))
}

func TestIsClinicalOrCriticalMethod_IncludesAuthLogout(testingInstance *testing.T) {
	assert.True(testingInstance, isClinicalOrCriticalMethod("/auth.v1.AuthService/Logout"))
}

func TestIsClinicalOrCriticalMethod_IncludesClinicalMethods(testingInstance *testing.T) {
	assert.True(testingInstance, isClinicalOrCriticalMethod("/patients.v1.PatientService/GetPatient"))
	assert.True(testingInstance, isClinicalOrCriticalMethod("/encounters.v1.EncounterService/GetEncounter"))
	assert.True(testingInstance, isClinicalOrCriticalMethod("/observations.v1.ObservationService/GetObservations"))
}