package clinical_test

import (
	"context"
	"errors"
	"testing"

	"github.com/healthcare/backend/internal/modules/clinical"
	"github.com/stretchr/testify/assert"
)

var errRepositoryFailure = errors.New("repository failure")

type mockRepository struct {
	createEncounterFn            func(ctx context.Context, encounter *clinical.Encounter) (*clinical.Encounter, error)
	getEncountersByPatientFn     func(ctx context.Context, patientFHIRID string) ([]*clinical.Encounter, error)
	createObservationFn          func(ctx context.Context, observation *clinical.Observation) (*clinical.Observation, error)
	getObservationsByEncounterFn func(ctx context.Context, encounterFHIRID string) ([]*clinical.Observation, error)
	getObservationsByPatientFn   func(ctx context.Context, patientFHIRID string) ([]*clinical.Observation, error)
	createConditionFn            func(ctx context.Context, condition *clinical.Condition) (*clinical.Condition, error)
	getConditionsByPatientFn     func(ctx context.Context, patientFHIRID string) ([]*clinical.Condition, error)
	createAllergyIntoleranceFn          func(ctx context.Context, allergy *clinical.AllergyIntolerance) (*clinical.AllergyIntolerance, error)
	getAllergyIntolerancesByPatientFn    func(ctx context.Context, patientFHIRID string) ([]*clinical.AllergyIntolerance, error)
	createMedicationRequestFn           func(ctx context.Context, medicationRequest *clinical.MedicationRequest) (*clinical.MedicationRequest, error)
	getMedicationRequestsByEncounterFn  func(ctx context.Context, encounterFHIRID string) ([]*clinical.MedicationRequest, error)
	createDiagnosticReportFn            func(ctx context.Context, report *clinical.DiagnosticReport) (*clinical.DiagnosticReport, error)
	getDiagnosticReportsByEncounterFn   func(ctx context.Context, encounterFHIRID string) ([]*clinical.DiagnosticReport, error)
}

func (repo *mockRepository) CreateEncounter(ctx context.Context, encounter *clinical.Encounter) (*clinical.Encounter, error) {
	if repo.createEncounterFn != nil {
		return repo.createEncounterFn(ctx, encounter)
	}
	return encounter, nil
}

func (repo *mockRepository) GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*clinical.Encounter, error) {
	if repo.getEncountersByPatientFn != nil {
		return repo.getEncountersByPatientFn(ctx, patientFHIRID)
	}
	return []*clinical.Encounter{}, nil
}

func (repo *mockRepository) CreateObservation(ctx context.Context, observation *clinical.Observation) (*clinical.Observation, error) {
	if repo.createObservationFn != nil {
		return repo.createObservationFn(ctx, observation)
	}
	return observation, nil
}

func (repo *mockRepository) GetObservationsByEncounter(ctx context.Context, encounterFHIRID string) ([]*clinical.Observation, error) {
	if repo.getObservationsByEncounterFn != nil {
		return repo.getObservationsByEncounterFn(ctx, encounterFHIRID)
	}
	return []*clinical.Observation{}, nil
}

func (repo *mockRepository) GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]*clinical.Observation, error) {
	if repo.getObservationsByPatientFn != nil {
		return repo.getObservationsByPatientFn(ctx, patientFHIRID)
	}
	return []*clinical.Observation{}, nil
}

func (repo *mockRepository) CreateCondition(ctx context.Context, condition *clinical.Condition) (*clinical.Condition, error) {
	if repo.createConditionFn != nil {
		return repo.createConditionFn(ctx, condition)
	}
	return condition, nil
}

func (repo *mockRepository) GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]*clinical.Condition, error) {
	if repo.getConditionsByPatientFn != nil {
		return repo.getConditionsByPatientFn(ctx, patientFHIRID)
	}
	return []*clinical.Condition{}, nil
}

func (repo *mockRepository) CreateAllergyIntolerance(ctx context.Context, allergy *clinical.AllergyIntolerance) (*clinical.AllergyIntolerance, error) {
	if repo.createAllergyIntoleranceFn != nil {
		return repo.createAllergyIntoleranceFn(ctx, allergy)
	}
	return allergy, nil
}

func (repo *mockRepository) GetAllergyIntolerancesByPatient(ctx context.Context, patientFHIRID string) ([]*clinical.AllergyIntolerance, error) {
	if repo.getAllergyIntolerancesByPatientFn != nil {
		return repo.getAllergyIntolerancesByPatientFn(ctx, patientFHIRID)
	}
	return []*clinical.AllergyIntolerance{}, nil
}

func (repo *mockRepository) CreateMedicationRequest(ctx context.Context, medicationRequest *clinical.MedicationRequest) (*clinical.MedicationRequest, error) {
	if repo.createMedicationRequestFn != nil {
		return repo.createMedicationRequestFn(ctx, medicationRequest)
	}
	return medicationRequest, nil
}

func (repo *mockRepository) GetMedicationRequestsByEncounter(ctx context.Context, encounterFHIRID string) ([]*clinical.MedicationRequest, error) {
	if repo.getMedicationRequestsByEncounterFn != nil {
		return repo.getMedicationRequestsByEncounterFn(ctx, encounterFHIRID)
	}
	return []*clinical.MedicationRequest{}, nil
}

func (repo *mockRepository) CreateDiagnosticReport(ctx context.Context, report *clinical.DiagnosticReport) (*clinical.DiagnosticReport, error) {
	if repo.createDiagnosticReportFn != nil {
		return repo.createDiagnosticReportFn(ctx, report)
	}
	return report, nil
}

func (repo *mockRepository) GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*clinical.DiagnosticReport, error) {
	if repo.getDiagnosticReportsByEncounterFn != nil {
		return repo.getDiagnosticReportsByEncounterFn(ctx, encounterFHIRID)
	}
	return []*clinical.DiagnosticReport{}, nil
}

func TestCreateEncounter_Success(t *testing.T) {
	svc := clinical.NewService(&mockRepository{})

	encounter := &clinical.Encounter{
		PatientFHIRID:  "patient-fhir-123",
		PractitionerID: "practitioner-456",
		ReasonCode:     "Z00.0",
		ReasonDisplay:  "Routine check-up",
	}

	result, err := svc.CreateEncounter(context.Background(), encounter)

	assert.NoError(t, err)
	assert.Equal(t, "patient-fhir-123", result.PatientFHIRID)
}

func TestCreateEncounter_MissingPatientFHIRID_ReturnsError(t *testing.T) {
	svc := clinical.NewService(&mockRepository{})

	encounter := &clinical.Encounter{PatientFHIRID: "", PractitionerID: "practitioner-456"}

	result, err := svc.CreateEncounter(context.Background(), encounter)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreateObservation_Success(t *testing.T) {
	svc := clinical.NewService(&mockRepository{})

	observation := &clinical.Observation{
		EncounterFHIRID: "encounter-123",
		PatientFHIRID:   "patient-456",
		LoincCode:       "55284-4",
		CodeDisplay:     "Blood pressure",
		ValueQuantity:   120,
		ValueUnit:       "mmHg",
	}

	result, err := svc.CreateObservation(context.Background(), observation)

	assert.NoError(t, err)
	assert.Equal(t, "55284-4", result.LoincCode)
}

func TestCreateObservation_MissingEncounterID_ReturnsError(t *testing.T) {
	svc := clinical.NewService(&mockRepository{})

	observation := &clinical.Observation{EncounterFHIRID: "", PatientFHIRID: "patient-456", LoincCode: "55284-4"}

	result, err := svc.CreateObservation(context.Background(), observation)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreateCondition_DefaultsClinicalStatusToActive(t *testing.T) {
	svc := clinical.NewService(&mockRepository{})

	condition := &clinical.Condition{
		PatientFHIRID:   "patient-123",
		EncounterFHIRID: "encounter-456",
		ICD10Code:       "I10",
		CodeDisplay:     "Essential hypertension",
		ClinicalStatus:  "",
	}

	result, err := svc.CreateCondition(context.Background(), condition)

	assert.NoError(t, err)
	assert.Equal(t, "active", result.ClinicalStatus)
}

func TestCreateCondition_MissingICD10Code_ReturnsError(t *testing.T) {
	svc := clinical.NewService(&mockRepository{})

	condition := &clinical.Condition{PatientFHIRID: "patient-123", ICD10Code: ""}

	result, err := svc.CreateCondition(context.Background(), condition)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreateAllergyIntolerance_DefaultsToActive(t *testing.T) {
	svc := clinical.NewService(&mockRepository{})

	allergy := &clinical.AllergyIntolerance{
		PatientFHIRID:   "patient-123",
		AllergenCode:    "7980",
		AllergenDisplay: "Penicillin",
		Reaction:        "Anaphylaxis",
	}

	result, err := svc.CreateAllergyIntolerance(context.Background(), allergy)

	assert.NoError(t, err)
	assert.Equal(t, "active", result.ClinicalStatus)
}

func TestCreateMedicationRequest_DefaultsToActive(t *testing.T) {
	svc := clinical.NewService(&mockRepository{})

	medication := &clinical.MedicationRequest{
		PatientFHIRID:      "patient-123",
		EncounterFHIRID:    "encounter-456",
		MedicationCode:     "10582",
		MedicationName:     "Amoxicillin",
		DosageInstructions: "500mg every 8 hours",
	}

	result, err := svc.CreateMedicationRequest(context.Background(), medication)

	assert.NoError(t, err)
	assert.Equal(t, "active", result.Status)
}

func TestCreateDiagnosticReport_DefaultsToFinal(t *testing.T) {
	svc := clinical.NewService(&mockRepository{})

	report := &clinical.DiagnosticReport{
		PatientFHIRID:   "patient-123",
		EncounterFHIRID: "encounter-456",
		ReportCode:      "24323-8",
		ReportDisplay:   "Complete blood count",
		Conclusion:      "Normal values",
	}

	result, err := svc.CreateDiagnosticReport(context.Background(), report)

	assert.NoError(t, err)
	assert.Equal(t, "final", result.Status)
}

func TestGetObservationsByEncounter_RepositoryFailure_ReturnsError(t *testing.T) {
	svc := clinical.NewService(&mockRepository{
		getObservationsByEncounterFn: func(_ context.Context, _ string) ([]*clinical.Observation, error) {
			return nil, errRepositoryFailure
		},
	})

	result, err := svc.GetObservationsByEncounter(context.Background(), "encounter-123")

	assert.Error(t, err)
	assert.Nil(t, result)
}
