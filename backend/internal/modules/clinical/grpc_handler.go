package clinical

import (
	"context"
	"errors"

	pb "github.com/healthcare/backend/internal/modules/clinical/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	service Service
	pb.UnimplementedClinicalServiceServer
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func mapClinicalError(err error) error {
	switch {
	case errors.Is(err, ErrEncounterNotFound):
		return apperrors.ErrEncounterNotFound.ToGRPC()
	case errors.Is(err, ErrObservationNotFound):
		return apperrors.ErrObservationNotFound.ToGRPC()
	case errors.Is(err, ErrConditionNotFound):
		return apperrors.ErrConditionNotFound.ToGRPC()
	case errors.Is(err, ErrAllergyNotFound):
		return apperrors.ErrAllergyIntoleranceNotFound.ToGRPC()
	case errors.Is(err, ErrMedicationRequestNotFound):
		return apperrors.ErrMedicationRequestNotFound.ToGRPC()
	case errors.Is(err, ErrDiagnosticReportNotFound):
		return apperrors.ErrDiagnosticReportNotFound.ToGRPC()
	default:
		return apperrors.ToGRPCStatus(err)
	}
}

func (handler *GRPCHandler) CreateEncounter(ctx context.Context, req *pb.CreateEncounterRequest) (*pb.CreateEncounterResponse, error) {
	violations := make(map[string]string)
	if req.PatientFhirId == "" {
		violations["patient_fhir_id"] = "is required"
	}
	if req.PractitionerId == "" {
		violations["practitioner_id"] = "is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.ErrBadRequest.WithFields(violations)
	}

	encounter := &Encounter{
		PatientFHIRID:  req.PatientFhirId,
		PractitionerID: req.PractitionerId,
		ReasonCode:     req.ReasonCode,
		ReasonDisplay:  req.ReasonDisplay,
		Status:         "in-progress",
	}

	createdEncounter, err := handler.service.CreateEncounter(ctx, encounter)
	if err != nil {
		return nil, mapClinicalError(err)
	}

	return &pb.CreateEncounterResponse{EncounterFhirId: createdEncounter.FHIRResourceID}, nil
}

func (handler *GRPCHandler) GetEncounters(ctx context.Context, req *pb.GetEncountersRequest) (*pb.GetEncountersResponse, error) {
	if req.PatientFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"patient_fhir_id": "is required"})
	}

	encounters, err := handler.service.GetEncountersByPatient(ctx, req.PatientFhirId)
	if err != nil {
		return nil, mapClinicalError(err)
	}

	pbEncounters := make([]*pb.Encounter, 0, len(encounters))
	for _, encounter := range encounters {
		pbEncounters = append(pbEncounters, &pb.Encounter{
			FhirId:        encounter.FHIRResourceID,
			PatientFhirId: encounter.PatientFHIRID,
			Status:        encounter.Status,
			ReasonDisplay: encounter.ReasonDisplay,
		})
	}

	return &pb.GetEncountersResponse{Encounters: pbEncounters}, nil
}

func (handler *GRPCHandler) CreateObservation(ctx context.Context, req *pb.CreateObservationRequest) (*pb.CreateObservationResponse, error) {
	violations := make(map[string]string)
	if req.EncounterFhirId == "" {
		violations["encounter_fhir_id"] = "is required"
	}
	if req.PatientFhirId == "" {
		violations["patient_fhir_id"] = "is required"
	}
	if req.LoincCode == "" {
		violations["loinc_code"] = "is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.ErrBadRequest.WithFields(violations)
	}

	observation := &Observation{
		EncounterFHIRID: req.EncounterFhirId,
		PatientFHIRID:   req.PatientFhirId,
		LoincCode:       req.LoincCode,
		CodeDisplay:     req.CodeDisplay,
		ValueQuantity:   req.ValueQuantity,
		ValueUnit:       req.ValueUnit,
	}

	createdObservation, err := handler.service.CreateObservation(ctx, observation)
	if err != nil {
		return nil, mapClinicalError(err)
	}

	return &pb.CreateObservationResponse{ObservationFhirId: createdObservation.FHIRResourceID}, nil
}

func (handler *GRPCHandler) GetObservations(ctx context.Context, req *pb.GetObservationsRequest) (*pb.GetObservationsResponse, error) {
	if req.EncounterFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"encounter_fhir_id": "is required"})
	}

	observations, err := handler.service.GetObservationsByEncounter(ctx, req.EncounterFhirId)
	if err != nil {
		return nil, mapClinicalError(err)
	}

	pbObservations := make([]*pb.Observation, 0, len(observations))
	for _, observation := range observations {
		pbObservations = append(pbObservations, &pb.Observation{
			FhirId:        observation.FHIRResourceID,
			LoincCode:     observation.LoincCode,
			CodeDisplay:   observation.CodeDisplay,
			ValueQuantity: observation.ValueQuantity,
			ValueUnit:     observation.ValueUnit,
		})
	}

	return &pb.GetObservationsResponse{Observations: pbObservations}, nil
}

func (handler *GRPCHandler) CreateCondition(ctx context.Context, req *pb.CreateConditionRequest) (*pb.CreateConditionResponse, error) {
	violations := make(map[string]string)
	if req.PatientFhirId == "" {
		violations["patient_fhir_id"] = "is required"
	}
	if req.Icd10Code == "" {
		violations["icd10_code"] = "is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.ErrBadRequest.WithFields(violations)
	}

	condition := &Condition{
		EncounterFHIRID: req.EncounterFhirId,
		PatientFHIRID:   req.PatientFhirId,
		ICD10Code:       req.Icd10Code,
		CodeDisplay:     req.CodeDisplay,
		ClinicalStatus:  req.ClinicalStatus,
	}

	createdCondition, err := handler.service.CreateCondition(ctx, condition)
	if err != nil {
		return nil, mapClinicalError(err)
	}

	return &pb.CreateConditionResponse{ConditionFhirId: createdCondition.FHIRResourceID}, nil
}

func (handler *GRPCHandler) GetConditions(ctx context.Context, req *pb.GetConditionsRequest) (*pb.GetConditionsResponse, error) {
	if req.PatientFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"patient_fhir_id": "is required"})
	}

	conditions, err := handler.service.GetConditionsByPatient(ctx, req.PatientFhirId)
	if err != nil {
		return nil, mapClinicalError(err)
	}

	pbConditions := make([]*pb.Condition, 0, len(conditions))
	for _, condition := range conditions {
		pbConditions = append(pbConditions, &pb.Condition{
			FhirId:         condition.FHIRResourceID,
			Icd10Code:      condition.ICD10Code,
			CodeDisplay:    condition.CodeDisplay,
			ClinicalStatus: condition.ClinicalStatus,
		})
	}

	return &pb.GetConditionsResponse{Conditions: pbConditions}, nil
}

func (handler *GRPCHandler) CreateAllergyIntolerance(ctx context.Context, req *pb.CreateAllergyIntoleranceRequest) (*pb.CreateAllergyIntoleranceResponse, error) {
	violations := make(map[string]string)
	if req.PatientFhirId == "" {
		violations["patient_fhir_id"] = "is required"
	}
	if req.AllergenCode == "" {
		violations["allergen_code"] = "is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.ErrBadRequest.WithFields(violations)
	}

	allergy := &AllergyIntolerance{
		PatientFHIRID:   req.PatientFhirId,
		AllergenCode:    req.AllergenCode,
		AllergenDisplay: req.AllergenDisplay,
		ClinicalStatus:  req.ClinicalStatus,
		Reaction:        req.Reaction,
	}

	createdAllergy, err := handler.service.CreateAllergyIntolerance(ctx, allergy)
	if err != nil {
		return nil, mapClinicalError(err)
	}

	return &pb.CreateAllergyIntoleranceResponse{AllergyFhirId: createdAllergy.FHIRResourceID}, nil
}

func (handler *GRPCHandler) GetAllergyIntolerances(ctx context.Context, req *pb.GetAllergyIntolerancesRequest) (*pb.GetAllergyIntolerancesResponse, error) {
	if req.PatientFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"patient_fhir_id": "is required"})
	}

	allergies, err := handler.service.GetAllergyIntolerancesByPatient(ctx, req.PatientFhirId)
	if err != nil {
		return nil, mapClinicalError(err)
	}

	pbAllergies := make([]*pb.AllergyIntolerance, 0, len(allergies))
	for _, allergy := range allergies {
		pbAllergies = append(pbAllergies, &pb.AllergyIntolerance{
			FhirId:          allergy.FHIRResourceID,
			AllergenDisplay: allergy.AllergenDisplay,
			ClinicalStatus:  allergy.ClinicalStatus,
			Reaction:        allergy.Reaction,
		})
	}

	return &pb.GetAllergyIntolerancesResponse{Allergies: pbAllergies}, nil
}

func (handler *GRPCHandler) CreateMedicationRequest(ctx context.Context, req *pb.CreateMedicationRequestRequest) (*pb.CreateMedicationRequestResponse, error) {
	violations := make(map[string]string)
	if req.PatientFhirId == "" {
		violations["patient_fhir_id"] = "is required"
	}
	if req.MedicationCode == "" {
		violations["medication_code"] = "is required"
	}
	if req.EncounterFhirId == "" {
		violations["encounter_fhir_id"] = "is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.ErrBadRequest.WithFields(violations)
	}

	medicationRequest := &MedicationRequest{
		EncounterFHIRID:    req.EncounterFhirId,
		PatientFHIRID:      req.PatientFhirId,
		PractitionerFHIRID: req.PractitionerFhirId,
		MedicationCode:     req.MedicationCode,
		MedicationName:     req.MedicationName,
		DosageInstructions: req.DosageInstructions,
	}

	createdMedication, err := handler.service.CreateMedicationRequest(ctx, medicationRequest)
	if err != nil {
		return nil, mapClinicalError(err)
	}

	return &pb.CreateMedicationRequestResponse{MedicationRequestFhirId: createdMedication.FHIRResourceID}, nil
}

func (handler *GRPCHandler) GetMedicationRequests(ctx context.Context, req *pb.GetMedicationRequestsRequest) (*pb.GetMedicationRequestsResponse, error) {
	if req.EncounterFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"encounter_fhir_id": "is required"})
	}

	medications, err := handler.service.GetMedicationRequestsByEncounter(ctx, req.EncounterFhirId)
	if err != nil {
		return nil, mapClinicalError(err)
	}

	pbMedications := make([]*pb.MedicationRequest, 0, len(medications))
	for _, medication := range medications {
		pbMedications = append(pbMedications, &pb.MedicationRequest{
			FhirId:             medication.FHIRResourceID,
			MedicationName:     medication.MedicationName,
			DosageInstructions: medication.DosageInstructions,
			Status:             medication.Status,
		})
	}

	return &pb.GetMedicationRequestsResponse{MedicationRequests: pbMedications}, nil
}

func (handler *GRPCHandler) CreateDiagnosticReport(ctx context.Context, req *pb.CreateDiagnosticReportRequest) (*pb.CreateDiagnosticReportResponse, error) {
	violations := make(map[string]string)
	if req.PatientFhirId == "" {
		violations["patient_fhir_id"] = "is required"
	}
	if req.ReportCode == "" {
		violations["report_code"] = "is required"
	}
	if req.EncounterFhirId == "" {
		violations["encounter_fhir_id"] = "is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.ErrBadRequest.WithFields(violations)
	}

	report := &DiagnosticReport{
		EncounterFHIRID: req.EncounterFhirId,
		PatientFHIRID:   req.PatientFhirId,
		ReportCode:      req.ReportCode,
		ReportDisplay:   req.ReportDisplay,
		Conclusion:      req.Conclusion,
	}

	createdReport, err := handler.service.CreateDiagnosticReport(ctx, report)
	if err != nil {
		return nil, mapClinicalError(err)
	}

	return &pb.CreateDiagnosticReportResponse{DiagnosticReportFhirId: createdReport.FHIRResourceID}, nil
}

func (handler *GRPCHandler) GetDiagnosticReports(ctx context.Context, req *pb.GetDiagnosticReportsRequest) (*pb.GetDiagnosticReportsResponse, error) {
	if req.EncounterFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"encounter_fhir_id": "is required"})
	}

	reports, err := handler.service.GetDiagnosticReportsByEncounter(ctx, req.EncounterFhirId)
	if err != nil {
		return nil, mapClinicalError(err)
	}

	pbReports := make([]*pb.DiagnosticReport, 0, len(reports))
	for _, report := range reports {
		pbReports = append(pbReports, &pb.DiagnosticReport{
			FhirId:        report.FHIRResourceID,
			ReportDisplay: report.ReportDisplay,
			Status:        report.Status,
			Conclusion:    report.Conclusion,
		})
	}

	return &pb.GetDiagnosticReportsResponse{DiagnosticReports: pbReports}, nil
}
