package adapters

import (
	"context"

	"github.com/healthcare/backend/internal/modules/allergy"
	"github.com/healthcare/backend/internal/modules/audit_logs"
	"github.com/healthcare/backend/internal/modules/medication"
)

const (
	medicationOverrideAuditMethod       = "medication.SafetyOverride"
	medicationOverrideAuditAction       = "medication_safety_override"
	medicationOverrideAuditResourceType = "medication_request"
)

type MedicationAllergyReader struct {
	allergyService allergy.Service
}

func NewMedicationAllergyReader(allergyService allergy.Service) *MedicationAllergyReader {
	return &MedicationAllergyReader{allergyService: allergyService}
}

func (reader *MedicationAllergyReader) ListPatientAllergies(ctx context.Context, patientFHIRID string) ([]medication.PatientAllergy, error) {
	patientAllergies, listErr := reader.allergyService.GetAllergyIntolerancesByPatient(ctx, patientFHIRID)
	if listErr != nil {
		return nil, listErr
	}
	allergies := make([]medication.PatientAllergy, 0, len(patientAllergies))
	for _, patientAllergy := range patientAllergies {
		if patientAllergy.PatientFHIRID != "" && patientAllergy.PatientFHIRID != patientFHIRID {
			continue
		}
		allergies = append(allergies, medication.PatientAllergy{
			FHIRResourceID:  patientAllergy.FHIRResourceID,
			AllergenCode:    patientAllergy.AllergenCode,
			AllergenDisplay: patientAllergy.AllergenDisplay,
			ClinicalStatus:  patientAllergy.ClinicalStatus,
		})
	}
	return allergies, nil
}

type MedicationAuditRecorder struct {
	auditService audit_logs.Service
}

func NewMedicationAuditRecorder(auditService audit_logs.Service) *MedicationAuditRecorder {
	return &MedicationAuditRecorder{auditService: auditService}
}

func (recorder *MedicationAuditRecorder) RecordOverrideAudit(ctx context.Context, event medication.OverrideAuditEvent) error {
	_, createErr := recorder.auditService.CreateResourceAuditLog(ctx, audit_logs.ResourceAuditLog{
		CorrelationID: event.CorrelationID,
		CallerUserID:  event.CallerUserID,
		CallerRole:    event.CallerRole,
		Method:        medicationOverrideAuditMethod,
		AccessGranted: true,
		ResourceType:  medicationOverrideAuditResourceType,
		ResourceID:    event.MedicationRequestFHIRID,
		Action:        medicationOverrideAuditAction,
		PayloadDiff: map[string]any{
			"patient_fhir_id":      event.PatientFHIRID,
			"medication_code":      event.MedicationCode,
			"alert_kinds":          event.AlertKinds,
			"related_resource_ids": event.RelatedResourceIDs,
		},
	})
	return createErr
}
