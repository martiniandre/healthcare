package medication

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
	"github.com/healthcare/backend/internal/shared/role"
)

type Service interface {
	CreateMedicationRequest(ctx context.Context, input CreateMedicationInput) (*CreateMedicationResult, error)
	GetMedicationRequestsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Medication, error)
	CheckSafety(ctx context.Context, input CheckSafetyInput) ([]SafetyAlert, error)
}

type SafetyDependencies struct {
	Allergies AllergyReader
	Catalog   SafetyCatalog
	Overrides OverrideRecorder
	Audit     AuditRecorder
}

type service struct {
	repo   Repository
	safety SafetyDependencies
}

func NewService(repo Repository, safety SafetyDependencies) Service {
	return &service{repo: repo, safety: safety}
}

func (medicationService *service) CreateMedicationRequest(ctx context.Context, input CreateMedicationInput) (*CreateMedicationResult, error) {
	if fieldViolations := validateMedicationFields(input); len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid medication request input", fieldViolations)
	}

	alerts, checkErr := medicationService.CheckSafety(ctx, CheckSafetyInput{
		PatientFHIRID:  input.PatientFHIRID,
		MedicationCode: input.MedicationCode,
		MedicationName: input.MedicationName,
	})
	if checkErr != nil {
		return nil, checkErr
	}

	blockingAlerts := filterBlockingAlerts(alerts)

	var pendingOverride *SafetyOverride
	if len(blockingAlerts) > 0 {
		justification := strings.TrimSpace(input.OverrideJustification)
		if utf8.RuneCountInString(justification) < MinimumOverrideJustificationLength {
			return nil, &BlockedByAllergyError{Alerts: alerts}
		}
		if !callerMayOverride(ctx) {
			return nil, apperrors.ErrPermissionDenied
		}
		preparedOverride, prepareErr := prepareOverride(ctx, input, blockingAlerts, justification)
		if prepareErr != nil {
			return nil, prepareErr
		}
		if recordErr := medicationService.safety.Overrides.RecordOverride(ctx, *preparedOverride); recordErr != nil {
			return nil, errors.Join(apperrors.ErrServiceUnavailable, recordErr)
		}
		pendingOverride = preparedOverride
	}

	newMedication := &Medication{
		EncounterFHIRID:    input.EncounterFHIRID,
		PatientFHIRID:      input.PatientFHIRID,
		PractitionerFHIRID: input.PractitionerFHIRID,
		MedicationCode:     input.MedicationCode,
		MedicationName:     input.MedicationName,
		DosageInstructions: input.DosageInstructions,
		Status:             "active",
	}

	createdMedication, createErr := medicationService.repo.CreateMedicationRequest(ctx, newMedication)
	if createErr != nil {
		if pendingOverride != nil {
			medicationService.discardOverride(ctx, pendingOverride.ID)
		}
		return nil, createErr
	}

	if pendingOverride != nil {
		medicationService.finalizeOverride(ctx, pendingOverride, createdMedication.FHIRResourceID, blockingAlerts)
	}

	return &CreateMedicationResult{Medication: createdMedication, Warnings: alerts}, nil
}

func (medicationService *service) GetMedicationRequestsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Medication, error) {
	return medicationService.repo.GetMedicationRequestsByEncounter(ctx, encounterFHIRID)
}

func (medicationService *service) CheckSafety(ctx context.Context, input CheckSafetyInput) ([]SafetyAlert, error) {
	if fieldViolations := validateSafetyFields(input); len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid medication safety check input", fieldViolations)
	}

	snapshot, snapshotErr := medicationService.safety.Catalog.LoadSnapshot(ctx)
	if snapshotErr != nil {
		return nil, errors.Join(apperrors.ErrServiceUnavailable, snapshotErr)
	}
	allergies, allergiesErr := medicationService.safety.Allergies.ListPatientAllergies(ctx, input.PatientFHIRID)
	if allergiesErr != nil {
		return nil, errors.Join(apperrors.ErrServiceUnavailable, allergiesErr)
	}
	activeMedications, medicationsErr := medicationService.repo.GetActiveMedicationRequestsByPatient(ctx, input.PatientFHIRID)
	if medicationsErr != nil {
		return nil, errors.Join(apperrors.ErrServiceUnavailable, medicationsErr)
	}

	requested := snapshot.ResolveMedication(input.MedicationCode, input.MedicationName)
	alerts := buildAllergyAlerts(snapshot, requested, allergies)
	alerts = append(alerts, buildInteractionAlerts(snapshot, requested, activeMedications)...)
	sortAlerts(alerts)
	return alerts, nil
}

func buildAllergyAlerts(snapshot *CatalogSnapshot, requested DrugProfile, allergies []PatientAllergy) []SafetyAlert {
	alerts := make([]SafetyAlert, 0)
	for _, patientAllergy := range allergies {
		if !isActiveAllergy(patientAllergy.ClinicalStatus) {
			continue
		}
		matchKind := snapshot.MatchAllergy(patientAllergy.AllergenCode, patientAllergy.AllergenDisplay, requested)
		if matchKind == "" {
			continue
		}
		allergenLabel := patientAllergy.AllergenDisplay
		if allergenLabel == "" {
			allergenLabel = patientAllergy.AllergenCode
		}
		alert := SafetyAlert{
			Kind:              AlertKindAllergy,
			Severity:          SeverityContraindicated,
			Blocking:          true,
			MessageKey:        AlertMessageKeyAllergyIngredient,
			Params:            map[string]string{"allergen": allergenLabel, "medication": medicationLabel(requested)},
			RelatedResourceID: patientAllergy.FHIRResourceID,
		}
		if matchKind == AllergyMatchClass {
			alert.MessageKey = AlertMessageKeyAllergyClass
			alert.Params["drugClass"] = requested.ClassName
		}
		alerts = append(alerts, alert)
	}
	return alerts
}

func buildInteractionAlerts(snapshot *CatalogSnapshot, requested DrugProfile, activeMedications []*Medication) []SafetyAlert {
	alerts := make([]SafetyAlert, 0)
	seenPairs := make(map[string]bool)
	for _, activeMedication := range activeMedications {
		activeProfile := snapshot.ResolveMedication(activeMedication.MedicationCode, activeMedication.MedicationName)
		for _, interaction := range snapshot.FindInteractions(requested, activeProfile) {
			pairKey := activeMedication.FHIRResourceID + "|" + interaction.DescriptionKey
			if seenPairs[pairKey] {
				continue
			}
			seenPairs[pairKey] = true
			alerts = append(alerts, SafetyAlert{
				Kind:       AlertKindInteraction,
				Severity:   interaction.Severity,
				Blocking:   false,
				MessageKey: interaction.DescriptionKey,
				Params: map[string]string{
					"medication":      medicationLabel(requested),
					"otherMedication": activeMedication.MedicationName,
				},
				RelatedResourceID: activeMedication.FHIRResourceID,
			})
		}
	}
	return alerts
}

func medicationLabel(profile DrugProfile) string {
	if profile.Name != "" {
		return profile.Name
	}
	return profile.Code
}

func isActiveAllergy(clinicalStatus string) bool {
	normalizedStatus := strings.ToLower(strings.TrimSpace(clinicalStatus))
	return normalizedStatus != "inactive" && normalizedStatus != "resolved"
}

func filterBlockingAlerts(alerts []SafetyAlert) []SafetyAlert {
	blockingAlerts := make([]SafetyAlert, 0)
	for _, alert := range alerts {
		if alert.Blocking {
			blockingAlerts = append(blockingAlerts, alert)
		}
	}
	return blockingAlerts
}

func callerMayOverride(ctx context.Context) bool {
	callerRole := role.Role(contextString(ctx, ctxkeys.RoleKey))
	return callerRole == role.RoleDoctor || callerRole == role.RoleAdmin
}

func prepareOverride(ctx context.Context, input CreateMedicationInput, blockingAlerts []SafetyAlert, justification string) (*SafetyOverride, error) {
	practitionerUserID, parseErr := uuid.Parse(contextString(ctx, ctxkeys.UserIDKey))
	if parseErr != nil {
		return nil, apperrors.ErrPermissionDenied
	}
	return &SafetyOverride{
		ID:                 uuid.New(),
		PatientFHIRID:      input.PatientFHIRID,
		MedicationCode:     input.MedicationCode,
		PractitionerUserID: practitionerUserID,
		AlertKinds:         distinctAlertKinds(blockingAlerts),
		Justification:      justification,
		CreatedAt:          time.Now(),
	}, nil
}

func distinctAlertKinds(alerts []SafetyAlert) []string {
	seenKinds := make(map[string]bool)
	kinds := make([]string, 0)
	for _, alert := range alerts {
		if !seenKinds[alert.Kind] {
			seenKinds[alert.Kind] = true
			kinds = append(kinds, alert.Kind)
		}
	}
	return kinds
}

func (medicationService *service) discardOverride(ctx context.Context, overrideID uuid.UUID) {
	if discardErr := medicationService.safety.Overrides.DiscardOverride(ctx, overrideID); discardErr != nil {
		slog.Error("failed to discard medication safety override after creation failure", "error", discardErr, "override_id", overrideID.String())
	}
}

func (medicationService *service) finalizeOverride(ctx context.Context, override *SafetyOverride, medicationRequestFHIRID string, blockingAlerts []SafetyAlert) {
	if linkErr := medicationService.safety.Overrides.LinkMedicationRequest(ctx, override.ID, medicationRequestFHIRID); linkErr != nil {
		slog.Error("failed to link medication safety override to medication request", "error", linkErr, "override_id", override.ID.String(), "medication_request_fhir_id", medicationRequestFHIRID)
	}
	relatedResourceIDs := make([]string, 0, len(blockingAlerts))
	for _, alert := range blockingAlerts {
		relatedResourceIDs = append(relatedResourceIDs, alert.RelatedResourceID)
	}
	auditErr := medicationService.safety.Audit.RecordOverrideAudit(ctx, OverrideAuditEvent{
		CorrelationID:           requestIDFromContext(ctx),
		CallerUserID:            contextString(ctx, ctxkeys.UserIDKey),
		CallerRole:              contextString(ctx, ctxkeys.RoleKey),
		PatientFHIRID:           override.PatientFHIRID,
		MedicationRequestFHIRID: medicationRequestFHIRID,
		MedicationCode:          override.MedicationCode,
		AlertKinds:              override.AlertKinds,
		RelatedResourceIDs:      relatedResourceIDs,
	})
	if auditErr != nil {
		slog.Error("failed to audit medication safety override", "error", auditErr, "override_id", override.ID.String(), "medication_request_fhir_id", medicationRequestFHIRID)
	}
}

func contextString(ctx context.Context, key ctxkeys.ContextKey) string {
	value, _ := ctx.Value(key).(string)
	return value
}

func requestIDFromContext(ctx context.Context) string {
	if requestID := contextString(ctx, ctxkeys.RequestIDKey); requestID != "" {
		return requestID
	}
	return contextString(ctx, ctxkeys.CorrelationIDKey)
}

func validateMedicationFields(input CreateMedicationInput) map[string]string {
	fieldViolations := make(map[string]string)
	if input.EncounterFHIRID == "" {
		fieldViolations["encounter_fhir_id"] = "is required"
	}
	if input.PatientFHIRID == "" {
		fieldViolations["patient_fhir_id"] = "is required"
	}
	if input.MedicationCode == "" && input.MedicationName == "" {
		fieldViolations["medication"] = "either name or code is required"
	}
	return fieldViolations
}

func validateSafetyFields(input CheckSafetyInput) map[string]string {
	fieldViolations := make(map[string]string)
	if input.PatientFHIRID == "" {
		fieldViolations["patient_fhir_id"] = "is required"
	}
	if input.MedicationCode == "" && input.MedicationName == "" {
		fieldViolations["medication"] = "either name or code is required"
	}
	return fieldViolations
}
