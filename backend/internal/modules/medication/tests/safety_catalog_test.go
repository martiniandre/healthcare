package tests

import (
	"testing"

	"github.com/healthcare/backend/internal/modules/medication"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeText(testingInstance *testing.T) {
	testCases := map[string]string{
		"Penicilina G":               "penicilinag",
		"  Ácido Acetil-Salicílico ": "acidoacetilsalicilico",
		"AMOXICILINA 500 mg":         "amoxicilina500mg",
		"Cefalosporína":              "cefalosporina",
		"Pão-de-açúcar ñ":            "paodeacucarn",
		"":                           "",
		"!!!":                        "",
	}
	for rawText, expected := range testCases {
		assert.Equal(testingInstance, expected, medication.NormalizeText(rawText), rawText)
	}
}

func TestResolveMedication(testingInstance *testing.T) {
	snapshot := sampleSnapshot()
	testCases := []struct {
		name               string
		code               string
		displayName        string
		expectedIngredient string
		expectedClass      string
		expectedCataloged  bool
	}{
		{name: "by catalog code", code: "AMOXICILLIN", displayName: "anything", expectedIngredient: "amoxicillin", expectedClass: "penicillins", expectedCataloged: true},
		{name: "by portuguese name with strength", displayName: "Amoxicilina 500mg", expectedIngredient: "amoxicillin", expectedClass: "penicillins", expectedCataloged: true},
		{name: "by multi word term", displayName: "Penicilina G benzatina", expectedIngredient: "benzylpenicillin", expectedClass: "penicillins", expectedCataloged: true},
		{name: "ingredient differs from name", displayName: "Dipirona", expectedIngredient: "metamizole", expectedClass: "pyrazolones", expectedCataloged: true},
		{name: "uncataloged drug of known class", displayName: "Penicilina V 500mg", expectedClass: "penicillins"},
		{name: "unknown drug", displayName: "Xyzzyfen", expectedCataloged: false},
	}
	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(subTest *testing.T) {
			profile := snapshot.ResolveMedication(testCase.code, testCase.displayName)

			assert.Equal(subTest, testCase.expectedIngredient, profile.Ingredient)
			assert.Equal(subTest, testCase.expectedClass, profile.ClassCode)
			assert.Equal(subTest, testCase.expectedCataloged, profile.Cataloged)
		})
	}
}

func TestMatchAllergy(testingInstance *testing.T) {
	snapshot := sampleSnapshot()
	testCases := []struct {
		name            string
		allergenCode    string
		allergenDisplay string
		drug            medication.DrugProfile
		expectedMatch   string
	}{
		{name: "same ingredient", allergenDisplay: "Amoxicilina", drug: snapshot.ResolveMedication("", "Amoxicillin 875mg"), expectedMatch: medication.AllergyMatchIngredient},
		{name: "class named in free text", allergenDisplay: "Alergia à penicilina", drug: snapshot.ResolveMedication("", "Ampicilina"), expectedMatch: medication.AllergyMatchClass},
		{name: "allergen identified by catalog code", allergenCode: "ibuprofen", drug: snapshot.ResolveMedication("", "Diclofenaco"), expectedMatch: medication.AllergyMatchClass},
		{name: "accents and case ignored", allergenDisplay: "ALERGIA A CEFALOSPORÍNAS", drug: snapshot.ResolveMedication("", "Cefalexina"), expectedMatch: medication.AllergyMatchClass},
		{name: "different class", allergenDisplay: "Penicilina", drug: snapshot.ResolveMedication("", "Cefalexina"), expectedMatch: ""},
		{name: "uncataloged identical code", allergenCode: "RX-777", drug: snapshot.ResolveMedication("rx777", "Mysterydrug"), expectedMatch: medication.AllergyMatchIngredient},
		{name: "uncataloged name contains allergen text", allergenDisplay: "Mysterydrug", drug: snapshot.ResolveMedication("", "Mysterydrug 10mg"), expectedMatch: medication.AllergyMatchIngredient},
		{name: "blank allergen", drug: snapshot.ResolveMedication("", "Amoxicilina"), expectedMatch: ""},
		{name: "short allergen text is not substring matched", allergenDisplay: "Ab", drug: snapshot.ResolveMedication("", "Abacavir"), expectedMatch: ""},
	}
	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(subTest *testing.T) {
			assert.Equal(subTest, testCase.expectedMatch, snapshot.MatchAllergy(testCase.allergenCode, testCase.allergenDisplay, testCase.drug))
		})
	}
}

func TestFindInteractions_ClassAndIngredientSubjects(testingInstance *testing.T) {
	snapshot := sampleSnapshot()
	warfarin := snapshot.ResolveMedication("", "Warfarin")
	ibuprofen := snapshot.ResolveMedication("", "Ibuprofeno")
	dipyrone := snapshot.ResolveMedication("", "Dipirona")

	assert.Len(testingInstance, snapshot.FindInteractions(ibuprofen, warfarin), 1)
	assert.Len(testingInstance, snapshot.FindInteractions(warfarin, ibuprofen), 1)
	assert.Empty(testingInstance, snapshot.FindInteractions(dipyrone, warfarin))
	assert.Empty(testingInstance, snapshot.FindInteractions(ibuprofen, ibuprofen))
}
