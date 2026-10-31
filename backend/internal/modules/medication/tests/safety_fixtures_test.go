package tests

import (
	"context"

	"github.com/healthcare/backend/internal/modules/medication"
	"github.com/healthcare/backend/internal/modules/medication/mocks"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
)

const (
	testDoctorUserID   = "11111111-1111-1111-1111-111111111111"
	testNurseUserID    = "22222222-2222-2222-2222-222222222222"
	validJustification = "Patient tolerated the drug before; benefit outweighs risk"
)

type safetyFixture struct {
	dependencies medication.SafetyDependencies
	allergies    *mocks.MockAllergyReader
	overrides    *mocks.MockOverrideRecorder
	audit        *mocks.MockAuditRecorder
}

func sampleSnapshot() *medication.CatalogSnapshot {
	classes := []medication.CatalogClass{
		{Code: "penicillins", Name: "Penicillins", SearchTerms: []string{"penicilina", "penicilinas", "penicillin"}},
		{Code: "cephalosporins", Name: "Cephalosporins", SearchTerms: []string{"cefalosporina", "cephalosporin"}},
		{Code: "nsaids", Name: "NSAIDs", SearchTerms: []string{"antiinflamatorio", "nsaid"}},
		{Code: "pyrazolones", Name: "Pyrazolones", SearchTerms: []string{"pirazolona"}},
		{Code: "coumarin_anticoagulants", Name: "Coumarin anticoagulants"},
	}
	entries := []medication.CatalogEntry{
		{Code: "amoxicillin", Name: "Amoxicillin", ActiveIngredient: "amoxicillin", ClassCode: "penicillins", SearchTerms: []string{"amoxicilina"}},
		{Code: "ampicillin", Name: "Ampicillin", ActiveIngredient: "ampicillin", ClassCode: "penicillins", SearchTerms: []string{"ampicilina"}},
		{Code: "penicillin_g", Name: "Penicillin G", ActiveIngredient: "benzylpenicillin", ClassCode: "penicillins", SearchTerms: []string{"penicilina g", "benzilpenicilina"}},
		{Code: "cephalexin", Name: "Cephalexin", ActiveIngredient: "cephalexin", ClassCode: "cephalosporins", SearchTerms: []string{"cefalexina"}},
		{Code: "ibuprofen", Name: "Ibuprofen", ActiveIngredient: "ibuprofen", ClassCode: "nsaids", SearchTerms: []string{"ibuprofeno"}},
		{Code: "diclofenac", Name: "Diclofenac", ActiveIngredient: "diclofenac", ClassCode: "nsaids", SearchTerms: []string{"diclofenaco"}},
		{Code: "dipyrone", Name: "Dipyrone", ActiveIngredient: "metamizole", ClassCode: "pyrazolones", SearchTerms: []string{"dipirona"}},
		{Code: "warfarin", Name: "Warfarin", ActiveIngredient: "warfarin", ClassCode: "coumarin_anticoagulants", SearchTerms: []string{"varfarina"}},
	}
	interactions := []medication.CatalogInteraction{
		{SubjectAType: "class", SubjectA: "nsaids", SubjectBType: "ingredient", SubjectB: "warfarin", Severity: "major", DescriptionKey: "medicationSafety.interaction.nsaidWarfarin"},
	}
	return medication.NewCatalogSnapshot(classes, entries, interactions)
}

func newSafetyDependencies(allergies []medication.PatientAllergy, readerErr error) *safetyFixture {
	allergyReader := &mocks.MockAllergyReader{Allergies: allergies, Err: readerErr}
	overrideRecorder := &mocks.MockOverrideRecorder{}
	auditRecorder := &mocks.MockAuditRecorder{}
	return &safetyFixture{
		dependencies: medication.SafetyDependencies{
			Allergies: allergyReader,
			Catalog:   &mocks.MockSafetyCatalog{Snapshot: sampleSnapshot()},
			Overrides: overrideRecorder,
			Audit:     auditRecorder,
		},
		allergies: allergyReader,
		overrides: overrideRecorder,
		audit:     auditRecorder,
	}
}

func contextWithCaller(roleName string, userID string) context.Context {
	callerContext := context.WithValue(context.Background(), ctxkeys.RoleKey, roleName)
	callerContext = context.WithValue(callerContext, ctxkeys.UserIDKey, userID)
	return context.WithValue(callerContext, ctxkeys.RequestIDKey, "request-1")
}
