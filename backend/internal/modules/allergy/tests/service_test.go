package tests

import (
	"context"
	"testing"

	"github.com/healthcare/backend/internal/modules/allergy"
	"github.com/healthcare/backend/internal/modules/allergy/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateAllergyIntolerance_DefaultsToActive(t *testing.T) {
	allergyService := allergy.NewService(&mocks.MockAllergyRepository{})

	entity := &allergy.Allergy{
		PatientFHIRID:   "patient-123",
		AllergenCode:    "7980",
		AllergenDisplay: "Penicillin",
		Reaction:        "Anaphylaxis",
	}

	result, err := allergyService.CreateAllergyIntolerance(context.Background(), entity)

	assert.NoError(t, err)
	assert.Equal(t, "active", result.ClinicalStatus)
}
