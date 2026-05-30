package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidEmail(testingInstance *testing.T) {
	assert.True(testingInstance, IsValidEmail("test@example.com"))
	assert.True(testingInstance, IsValidEmail("user.name+tag@domain.co.uk"))
	assert.False(testingInstance, IsValidEmail("invalid-email"))
	assert.False(testingInstance, IsValidEmail("test@"))
	assert.False(testingInstance, IsValidEmail("@example.com"))
}

func TestIsValidPhone(testingInstance *testing.T) {
	assert.True(testingInstance, IsValidPhone("(11) 98765-4321"))
	assert.True(testingInstance, IsValidPhone("+5511987654321"))
	assert.False(testingInstance, IsValidPhone("11987654321"))
	assert.False(testingInstance, IsValidPhone("phone"))
}

func TestIsValidICD10(testingInstance *testing.T) {
	assert.True(testingInstance, IsValidICD10("I10"))
	assert.True(testingInstance, IsValidICD10("E11.9"))
	assert.True(testingInstance, IsValidICD10("J45.01"))
	assert.False(testingInstance, IsValidICD10("invalid"))
	assert.False(testingInstance, IsValidICD10("123"))
}

func TestIsValidLOINC(testingInstance *testing.T) {
	assert.True(testingInstance, IsValidLOINC("8867-4"))
	assert.True(testingInstance, IsValidLOINC("8310-5"))
	assert.False(testingInstance, IsValidLOINC("88674"))
	assert.False(testingInstance, IsValidLOINC("invalid"))
}

func TestIsValidDICOMModality(testingInstance *testing.T) {
	assert.True(testingInstance, IsValidDICOMModality("CT"))
	assert.True(testingInstance, IsValidDICOMModality("MR"))
	assert.True(testingInstance, IsValidDICOMModality("us"))
	assert.False(testingInstance, IsValidDICOMModality("INVALID"))
}

func TestIsValidClinicalStatus(testingInstance *testing.T) {
	assert.True(testingInstance, IsValidClinicalStatus("active"))
	assert.True(testingInstance, IsValidClinicalStatus("resolved"))
	assert.False(testingInstance, IsValidClinicalStatus("unknown"))
}

func TestIsValidGender(testingInstance *testing.T) {
	assert.True(testingInstance, IsValidGender("male"))
	assert.True(testingInstance, IsValidGender("female"))
	assert.False(testingInstance, IsValidGender("invalid"))
}

func TestIsValidCRMNumber(testingInstance *testing.T) {
	assert.True(testingInstance, IsValidCRMNumber("CRM-SP 12345"))
	assert.True(testingInstance, IsValidCRMNumber("COREN-RJ 54321"))
	assert.True(testingInstance, IsValidCRMNumber("CRM-12345"))
	assert.True(testingInstance, IsValidCRMNumber(""))
	assert.False(testingInstance, IsValidCRMNumber("INVALID"))
}

func TestIsValidCPF(testingInstance *testing.T) {
	assert.True(testingInstance, IsValidCPF("123.456.789-09"))
	assert.True(testingInstance, IsValidCPF("12345678909"))
	assert.False(testingInstance, IsValidCPF("111.111.111-11"))
	assert.False(testingInstance, IsValidCPF("123.456.789-00"))
	assert.False(testingInstance, IsValidCPF("invalid"))
}

func TestIsValidObservationRange(testingInstance *testing.T) {
	assert.True(testingInstance, IsValidObservationRange("8867-4", 75))
	assert.False(testingInstance, IsValidObservationRange("8867-4", 350))
	assert.True(testingInstance, IsValidObservationRange("8310-5", 36.5))
	assert.False(testingInstance, IsValidObservationRange("8310-5", 25))
}
