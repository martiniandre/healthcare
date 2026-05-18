package clinical

import "context"

type Service interface {
	CreateEncounter(ctx context.Context, encounter *Encounter) (*Encounter, error)
	GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*Encounter, error)

	CreateObservation(ctx context.Context, observation *Observation) (*Observation, error)
	GetObservationsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Observation, error)
	GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]*Observation, error)

	CreateCondition(ctx context.Context, condition *Condition) (*Condition, error)
	GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]*Condition, error)

	CreateAllergyIntolerance(ctx context.Context, allergy *AllergyIntolerance) (*AllergyIntolerance, error)
	GetAllergyIntolerancesByPatient(ctx context.Context, patientFHIRID string) ([]*AllergyIntolerance, error)

	CreateMedicationRequest(ctx context.Context, medicationRequest *MedicationRequest) (*MedicationRequest, error)
	GetMedicationRequestsByEncounter(ctx context.Context, encounterFHIRID string) ([]*MedicationRequest, error)

	CreateDiagnosticReport(ctx context.Context, report *DiagnosticReport) (*DiagnosticReport, error)
	GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*DiagnosticReport, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (clinicalService *service) CreateEncounter(ctx context.Context, encounter *Encounter) (*Encounter, error) {
	if encounter.PatientFHIRID == "" {
		return nil, ErrEncounterNotFound
	}
	return clinicalService.repo.CreateEncounter(ctx, encounter)
}

func (clinicalService *service) GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*Encounter, error) {
	return clinicalService.repo.GetEncountersByPatient(ctx, patientFHIRID)
}

func (clinicalService *service) CreateObservation(ctx context.Context, observation *Observation) (*Observation, error) {
	if observation.EncounterFHIRID == "" || observation.PatientFHIRID == "" {
		return nil, ErrObservationNotFound
	}
	return clinicalService.repo.CreateObservation(ctx, observation)
}

func (clinicalService *service) GetObservationsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Observation, error) {
	return clinicalService.repo.GetObservationsByEncounter(ctx, encounterFHIRID)
}

func (clinicalService *service) GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]*Observation, error) {
	return clinicalService.repo.GetObservationsByPatient(ctx, patientFHIRID)
}

func (clinicalService *service) CreateCondition(ctx context.Context, condition *Condition) (*Condition, error) {
	if condition.PatientFHIRID == "" || condition.ICD10Code == "" {
		return nil, ErrConditionNotFound
	}
	if condition.ClinicalStatus == "" {
		condition.ClinicalStatus = "active"
	}
	return clinicalService.repo.CreateCondition(ctx, condition)
}

func (clinicalService *service) GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]*Condition, error) {
	return clinicalService.repo.GetConditionsByPatient(ctx, patientFHIRID)
}

func (clinicalService *service) CreateAllergyIntolerance(ctx context.Context, allergy *AllergyIntolerance) (*AllergyIntolerance, error) {
	if allergy.PatientFHIRID == "" || allergy.AllergenCode == "" {
		return nil, ErrAllergyNotFound
	}
	if allergy.ClinicalStatus == "" {
		allergy.ClinicalStatus = "active"
	}
	return clinicalService.repo.CreateAllergyIntolerance(ctx, allergy)
}

func (clinicalService *service) GetAllergyIntolerancesByPatient(ctx context.Context, patientFHIRID string) ([]*AllergyIntolerance, error) {
	return clinicalService.repo.GetAllergyIntolerancesByPatient(ctx, patientFHIRID)
}

func (clinicalService *service) CreateMedicationRequest(ctx context.Context, medicationRequest *MedicationRequest) (*MedicationRequest, error) {
	if medicationRequest.PatientFHIRID == "" || medicationRequest.MedicationCode == "" {
		return nil, ErrMedicationRequestNotFound
	}
	if medicationRequest.Status == "" {
		medicationRequest.Status = "active"
	}
	return clinicalService.repo.CreateMedicationRequest(ctx, medicationRequest)
}

func (clinicalService *service) GetMedicationRequestsByEncounter(ctx context.Context, encounterFHIRID string) ([]*MedicationRequest, error) {
	return clinicalService.repo.GetMedicationRequestsByEncounter(ctx, encounterFHIRID)
}

func (clinicalService *service) CreateDiagnosticReport(ctx context.Context, report *DiagnosticReport) (*DiagnosticReport, error) {
	if report.PatientFHIRID == "" || report.ReportCode == "" {
		return nil, ErrDiagnosticReportNotFound
	}
	if report.Status == "" {
		report.Status = "final"
	}
	return clinicalService.repo.CreateDiagnosticReport(ctx, report)
}

func (clinicalService *service) GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*DiagnosticReport, error) {
	return clinicalService.repo.GetDiagnosticReportsByEncounter(ctx, encounterFHIRID)
}
