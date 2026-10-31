package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/healthcare/backend/internal/modules/medication"
	"github.com/healthcare/backend/internal/modules/medication/mocks"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prescriptionInput(medicationName string, justification string) medication.CreateMedicationInput {
	return medication.CreateMedicationInput{
		EncounterFHIRID:       "encounter-1",
		PatientFHIRID:         "patient-1",
		PractitionerFHIRID:    "practitioner-1",
		MedicationName:        medicationName,
		DosageInstructions:    "1 tablet every 8 hours",
		OverrideJustification: justification,
	}
}

func penicillinAllergy() medication.PatientAllergy {
	return medication.PatientAllergy{FHIRResourceID: "allergy-1", AllergenCode: "91936005", AllergenDisplay: "Alergia a amoxicilina", ClinicalStatus: "active"}
}

func countingRepository(createdCounter *int) *mocks.MockMedicationRepository {
	return &mocks.MockMedicationRepository{
		CreateMedicationRequestFn: func(ctx context.Context, entity *medication.Medication) (*medication.Medication, error) {
			*createdCounter++
			entity.FHIRResourceID = "medication-request-1"
			return entity, nil
		},
	}
}

func TestCheckSafety_AllergyByIngredient_ReturnsBlockingAlert(testingInstance *testing.T) {
	fixture := newSafetyDependencies([]medication.PatientAllergy{penicillinAllergy()}, nil)
	medicationService := medication.NewService(&mocks.MockMedicationRepository{}, fixture.dependencies)

	alerts, err := medicationService.CheckSafety(context.Background(), medication.CheckSafetyInput{PatientFHIRID: "patient-1", MedicationName: "Amoxicilina 500mg"})

	require.NoError(testingInstance, err)
	require.Len(testingInstance, alerts, 1)
	assert.Equal(testingInstance, medication.AlertKindAllergy, alerts[0].Kind)
	assert.Equal(testingInstance, medication.SeverityContraindicated, alerts[0].Severity)
	assert.True(testingInstance, alerts[0].Blocking)
	assert.Equal(testingInstance, medication.AlertMessageKeyAllergyIngredient, alerts[0].MessageKey)
	assert.Equal(testingInstance, "allergy-1", alerts[0].RelatedResourceID)
	assert.Equal(testingInstance, "Alergia a amoxicilina", alerts[0].Params["allergen"])
}

func TestCheckSafety_CrossAllergyByClass_ReturnsBlockingAlertWithDrugClass(testingInstance *testing.T) {
	fixture := newSafetyDependencies([]medication.PatientAllergy{penicillinAllergy()}, nil)
	medicationService := medication.NewService(&mocks.MockMedicationRepository{}, fixture.dependencies)

	alerts, err := medicationService.CheckSafety(context.Background(), medication.CheckSafetyInput{PatientFHIRID: "patient-1", MedicationName: "Ampicilina 500mg"})

	require.NoError(testingInstance, err)
	require.Len(testingInstance, alerts, 1)
	assert.Equal(testingInstance, medication.AlertMessageKeyAllergyClass, alerts[0].MessageKey)
	assert.Equal(testingInstance, "Penicillins", alerts[0].Params["drugClass"])
	assert.True(testingInstance, alerts[0].Blocking)
}

func TestCheckSafety_UnrelatedClass_ReturnsNoAlerts(testingInstance *testing.T) {
	fixture := newSafetyDependencies([]medication.PatientAllergy{penicillinAllergy()}, nil)
	medicationService := medication.NewService(&mocks.MockMedicationRepository{}, fixture.dependencies)

	alerts, err := medicationService.CheckSafety(context.Background(), medication.CheckSafetyInput{PatientFHIRID: "patient-1", MedicationName: "Dipirona 500mg"})

	require.NoError(testingInstance, err)
	assert.Empty(testingInstance, alerts)
}

func TestCheckSafety_InactiveAllergyIsIgnored(testingInstance *testing.T) {
	inactiveAllergy := penicillinAllergy()
	inactiveAllergy.ClinicalStatus = "inactive"
	resolvedAllergy := penicillinAllergy()
	resolvedAllergy.FHIRResourceID = "allergy-2"
	resolvedAllergy.ClinicalStatus = "resolved"
	fixture := newSafetyDependencies([]medication.PatientAllergy{inactiveAllergy, resolvedAllergy}, nil)
	medicationService := medication.NewService(&mocks.MockMedicationRepository{}, fixture.dependencies)

	alerts, err := medicationService.CheckSafety(context.Background(), medication.CheckSafetyInput{PatientFHIRID: "patient-1", MedicationName: "Amoxicilina"})

	require.NoError(testingInstance, err)
	assert.Empty(testingInstance, alerts)
}

func TestCheckSafety_InteractionWithActiveMedication_WarnsWithoutBlocking(testingInstance *testing.T) {
	fixture := newSafetyDependencies(nil, nil)
	repository := &mocks.MockMedicationRepository{
		GetActiveMedicationRequestsByPatientFn: func(ctx context.Context, patientFHIRID string) ([]*medication.Medication, error) {
			return []*medication.Medication{{FHIRResourceID: "request-warfarin", MedicationName: "Varfarina 5mg", Status: "active"}}, nil
		},
	}
	medicationService := medication.NewService(repository, fixture.dependencies)

	alerts, err := medicationService.CheckSafety(context.Background(), medication.CheckSafetyInput{PatientFHIRID: "patient-1", MedicationName: "Ibuprofeno 400mg"})

	require.NoError(testingInstance, err)
	require.Len(testingInstance, alerts, 1)
	assert.Equal(testingInstance, medication.AlertKindInteraction, alerts[0].Kind)
	assert.Equal(testingInstance, medication.SeverityMajor, alerts[0].Severity)
	assert.False(testingInstance, alerts[0].Blocking)
	assert.Equal(testingInstance, "medicationSafety.interaction.nsaidWarfarin", alerts[0].MessageKey)
	assert.Equal(testingInstance, "request-warfarin", alerts[0].RelatedResourceID)
	assert.Equal(testingInstance, "Varfarina 5mg", alerts[0].Params["otherMedication"])
}

func TestCheckSafety_InteractionIsSymmetric(testingInstance *testing.T) {
	fixture := newSafetyDependencies(nil, nil)
	repository := &mocks.MockMedicationRepository{
		GetActiveMedicationRequestsByPatientFn: func(ctx context.Context, patientFHIRID string) ([]*medication.Medication, error) {
			return []*medication.Medication{{FHIRResourceID: "request-diclofenac", MedicationName: "Diclofenaco", Status: "active"}}, nil
		},
	}
	medicationService := medication.NewService(repository, fixture.dependencies)

	alerts, err := medicationService.CheckSafety(context.Background(), medication.CheckSafetyInput{PatientFHIRID: "patient-1", MedicationName: "Warfarin"})

	require.NoError(testingInstance, err)
	require.Len(testingInstance, alerts, 1)
	assert.False(testingInstance, alerts[0].Blocking)
}

func TestCheckSafety_BlockingAlertsComeFirst(testingInstance *testing.T) {
	fixture := newSafetyDependencies([]medication.PatientAllergy{{FHIRResourceID: "allergy-nsaid", AllergenDisplay: "Ibuprofeno", ClinicalStatus: "active"}}, nil)
	repository := &mocks.MockMedicationRepository{
		GetActiveMedicationRequestsByPatientFn: func(ctx context.Context, patientFHIRID string) ([]*medication.Medication, error) {
			return []*medication.Medication{{FHIRResourceID: "request-warfarin", MedicationName: "Warfarin", Status: "active"}}, nil
		},
	}
	medicationService := medication.NewService(repository, fixture.dependencies)

	alerts, err := medicationService.CheckSafety(context.Background(), medication.CheckSafetyInput{PatientFHIRID: "patient-1", MedicationName: "Diclofenac"})

	require.NoError(testingInstance, err)
	require.Len(testingInstance, alerts, 2)
	assert.True(testingInstance, alerts[0].Blocking)
	assert.False(testingInstance, alerts[1].Blocking)
}

func TestCheckSafety_DependencyFailure_FailsClosed(testingInstance *testing.T) {
	failure := errors.New("fhir unavailable")
	testCases := map[string]func() medication.Service{
		"allergy reader": func() medication.Service {
			return medication.NewService(&mocks.MockMedicationRepository{}, newSafetyDependencies(nil, failure).dependencies)
		},
		"active medications": func() medication.Service {
			repository := &mocks.MockMedicationRepository{
				GetActiveMedicationRequestsByPatientFn: func(ctx context.Context, patientFHIRID string) ([]*medication.Medication, error) {
					return nil, failure
				},
			}
			return medication.NewService(repository, newSafetyDependencies(nil, nil).dependencies)
		},
		"catalog": func() medication.Service {
			fixture := newSafetyDependencies(nil, nil)
			fixture.dependencies.Catalog = &mocks.MockSafetyCatalog{Err: failure}
			return medication.NewService(&mocks.MockMedicationRepository{}, fixture.dependencies)
		},
	}
	for name, buildService := range testCases {
		testingInstance.Run(name, func(subTest *testing.T) {
			alerts, err := buildService().CheckSafety(context.Background(), medication.CheckSafetyInput{PatientFHIRID: "patient-1", MedicationName: "Amoxicilina"})

			assert.Nil(subTest, alerts)
			assert.ErrorIs(subTest, err, failure)
			var appError apperrors.AppError
			require.True(subTest, errors.As(err, &appError))
			assert.Equal(subTest, apperrors.ErrServiceUnavailable.HTTPCode, appError.HTTPCode)
		})
	}
}

func TestCheckSafety_MissingFields_ReturnsInvalidArgument(testingInstance *testing.T) {
	medicationService := medication.NewService(&mocks.MockMedicationRepository{}, newSafetyDependencies(nil, nil).dependencies)

	_, err := medicationService.CheckSafety(context.Background(), medication.CheckSafetyInput{})

	var appError apperrors.AppError
	require.True(testingInstance, errors.As(err, &appError))
	assert.Equal(testingInstance, apperrors.ErrBadRequest.HTTPCode, appError.HTTPCode)
}

func TestCreateMedicationRequest_NoAlerts_CreatesWithoutOverride(testingInstance *testing.T) {
	fixture := newSafetyDependencies(nil, nil)
	createdCounter := 0
	medicationService := medication.NewService(countingRepository(&createdCounter), fixture.dependencies)

	result, err := medicationService.CreateMedicationRequest(contextWithCaller("DOCTOR", testDoctorUserID), prescriptionInput("Amoxicilina", ""))

	require.NoError(testingInstance, err)
	assert.Equal(testingInstance, 1, createdCounter)
	assert.Empty(testingInstance, result.Warnings)
	assert.Empty(testingInstance, fixture.overrides.Recorded)
	assert.Empty(testingInstance, fixture.audit.Events)
}

func TestCreateMedicationRequest_InteractionWarning_DoesNotBlock(testingInstance *testing.T) {
	fixture := newSafetyDependencies(nil, nil)
	createdCounter := 0
	repository := countingRepository(&createdCounter)
	repository.GetActiveMedicationRequestsByPatientFn = func(ctx context.Context, patientFHIRID string) ([]*medication.Medication, error) {
		return []*medication.Medication{{FHIRResourceID: "request-warfarin", MedicationName: "Warfarin", Status: "active"}}, nil
	}
	medicationService := medication.NewService(repository, fixture.dependencies)

	result, err := medicationService.CreateMedicationRequest(contextWithCaller("DOCTOR", testDoctorUserID), prescriptionInput("Ibuprofeno", ""))

	require.NoError(testingInstance, err)
	assert.Equal(testingInstance, 1, createdCounter)
	require.Len(testingInstance, result.Warnings, 1)
	assert.False(testingInstance, result.Warnings[0].Blocking)
	assert.Empty(testingInstance, fixture.overrides.Recorded)
}

func TestCreateMedicationRequest_AllergyWithoutJustification_IsBlocked(testingInstance *testing.T) {
	fixture := newSafetyDependencies([]medication.PatientAllergy{penicillinAllergy()}, nil)
	createdCounter := 0
	medicationService := medication.NewService(countingRepository(&createdCounter), fixture.dependencies)

	result, err := medicationService.CreateMedicationRequest(contextWithCaller("DOCTOR", testDoctorUserID), prescriptionInput("Amoxicilina", ""))

	assert.Nil(testingInstance, result)
	var blockedError *medication.BlockedByAllergyError
	require.True(testingInstance, errors.As(err, &blockedError))
	require.Len(testingInstance, blockedError.Alerts, 1)
	assert.Equal(testingInstance, "allergy-1", blockedError.Alerts[0].RelatedResourceID)
	var appError apperrors.AppError
	require.True(testingInstance, errors.As(err, &appError))
	assert.Equal(testingInstance, "MEDICATION_ALLERGY_BLOCKED", appError.Code)
	assert.Equal(testingInstance, 409, appError.HTTPCode)
	assert.Equal(testingInstance, 0, createdCounter)
	assert.Empty(testingInstance, fixture.overrides.Recorded)
}

func TestCreateMedicationRequest_ShortJustification_IsRejected(testingInstance *testing.T) {
	fixture := newSafetyDependencies([]medication.PatientAllergy{penicillinAllergy()}, nil)
	createdCounter := 0
	medicationService := medication.NewService(countingRepository(&createdCounter), fixture.dependencies)

	for _, shortJustification := range []string{"short", "   too short  ", "123456789"} {
		_, err := medicationService.CreateMedicationRequest(contextWithCaller("DOCTOR", testDoctorUserID), prescriptionInput("Amoxicilina", shortJustification))

		var blockedError *medication.BlockedByAllergyError
		assert.True(testingInstance, errors.As(err, &blockedError), "justification %q must be rejected", shortJustification)
	}
	assert.Equal(testingInstance, 0, createdCounter)
}

func TestCreateMedicationRequest_OverrideByDoctor_CreatesRecordsAndAudits(testingInstance *testing.T) {
	fixture := newSafetyDependencies([]medication.PatientAllergy{penicillinAllergy()}, nil)
	createdCounter := 0
	medicationService := medication.NewService(countingRepository(&createdCounter), fixture.dependencies)

	result, err := medicationService.CreateMedicationRequest(contextWithCaller("DOCTOR", testDoctorUserID), prescriptionInput("Amoxicilina", validJustification))

	require.NoError(testingInstance, err)
	assert.Equal(testingInstance, 1, createdCounter)
	assert.Equal(testingInstance, "medication-request-1", result.FHIRResourceID)
	require.Len(testingInstance, result.Warnings, 1)
	assert.True(testingInstance, result.Warnings[0].Blocking)

	require.Len(testingInstance, fixture.overrides.Recorded, 1)
	recordedOverride := fixture.overrides.Recorded[0]
	assert.Equal(testingInstance, "patient-1", recordedOverride.PatientFHIRID)
	assert.Equal(testingInstance, testDoctorUserID, recordedOverride.PractitionerUserID.String())
	assert.Equal(testingInstance, validJustification, recordedOverride.Justification)
	assert.Equal(testingInstance, []string{medication.AlertKindAllergy}, recordedOverride.AlertKinds)
	assert.Equal(testingInstance, "medication-request-1", fixture.overrides.Linked[recordedOverride.ID])

	require.Len(testingInstance, fixture.audit.Events, 1)
	auditEvent := fixture.audit.Events[0]
	assert.Equal(testingInstance, "medication-request-1", auditEvent.MedicationRequestFHIRID)
	assert.Equal(testingInstance, "DOCTOR", auditEvent.CallerRole)
	assert.Equal(testingInstance, testDoctorUserID, auditEvent.CallerUserID)
	assert.Equal(testingInstance, "request-1", auditEvent.CorrelationID)
	assert.Equal(testingInstance, []string{"allergy-1"}, auditEvent.RelatedResourceIDs)
}

func TestCreateMedicationRequest_OverrideByAdmin_IsAllowed(testingInstance *testing.T) {
	fixture := newSafetyDependencies([]medication.PatientAllergy{penicillinAllergy()}, nil)
	createdCounter := 0
	medicationService := medication.NewService(countingRepository(&createdCounter), fixture.dependencies)

	_, err := medicationService.CreateMedicationRequest(contextWithCaller("ADMIN", testDoctorUserID), prescriptionInput("Amoxicilina", validJustification))

	require.NoError(testingInstance, err)
	assert.Equal(testingInstance, 1, createdCounter)
}

func TestCreateMedicationRequest_OverrideByNurse_IsDenied(testingInstance *testing.T) {
	fixture := newSafetyDependencies([]medication.PatientAllergy{penicillinAllergy()}, nil)
	createdCounter := 0
	medicationService := medication.NewService(countingRepository(&createdCounter), fixture.dependencies)

	result, err := medicationService.CreateMedicationRequest(contextWithCaller("NURSE", testNurseUserID), prescriptionInput("Amoxicilina", validJustification))

	assert.Nil(testingInstance, result)
	assert.ErrorIs(testingInstance, err, apperrors.ErrPermissionDenied)
	assert.Equal(testingInstance, 0, createdCounter)
	assert.Empty(testingInstance, fixture.overrides.Recorded)
	assert.Empty(testingInstance, fixture.audit.Events)
}

func TestCreateMedicationRequest_OverrideWithoutCallerIdentity_IsDenied(testingInstance *testing.T) {
	fixture := newSafetyDependencies([]medication.PatientAllergy{penicillinAllergy()}, nil)
	medicationService := medication.NewService(&mocks.MockMedicationRepository{}, fixture.dependencies)

	_, err := medicationService.CreateMedicationRequest(contextWithCaller("DOCTOR", "not-a-uuid"), prescriptionInput("Amoxicilina", validJustification))

	assert.ErrorIs(testingInstance, err, apperrors.ErrPermissionDenied)
}

func TestCreateMedicationRequest_OverrideRecordFailure_AbortsBeforeCreating(testingInstance *testing.T) {
	fixture := newSafetyDependencies([]medication.PatientAllergy{penicillinAllergy()}, nil)
	fixture.overrides.RecordErr = errors.New("database down")
	createdCounter := 0
	medicationService := medication.NewService(countingRepository(&createdCounter), fixture.dependencies)

	_, err := medicationService.CreateMedicationRequest(contextWithCaller("DOCTOR", testDoctorUserID), prescriptionInput("Amoxicilina", validJustification))

	require.Error(testingInstance, err)
	assert.Equal(testingInstance, 0, createdCounter)
}

func TestCreateMedicationRequest_CreationFailureAfterOverride_DiscardsOverride(testingInstance *testing.T) {
	fixture := newSafetyDependencies([]medication.PatientAllergy{penicillinAllergy()}, nil)
	repository := &mocks.MockMedicationRepository{
		CreateMedicationRequestFn: func(ctx context.Context, entity *medication.Medication) (*medication.Medication, error) {
			return nil, errors.New("fhir unavailable")
		},
	}
	medicationService := medication.NewService(repository, fixture.dependencies)

	_, err := medicationService.CreateMedicationRequest(contextWithCaller("DOCTOR", testDoctorUserID), prescriptionInput("Amoxicilina", validJustification))

	require.Error(testingInstance, err)
	require.Len(testingInstance, fixture.overrides.Recorded, 1)
	assert.Equal(testingInstance, []string{fixture.overrides.Recorded[0].ID.String()}, discardedIDs(fixture.overrides))
	assert.Empty(testingInstance, fixture.audit.Events)
}

func discardedIDs(recorder *mocks.MockOverrideRecorder) []string {
	identifiers := make([]string, 0, len(recorder.Discarded))
	for _, discardedID := range recorder.Discarded {
		identifiers = append(identifiers, discardedID.String())
	}
	return identifiers
}
