package medication

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SafetyRepository struct {
	dbPool *pgxpool.Pool
}

func NewSafetyRepository(dbPool *pgxpool.Pool) *SafetyRepository {
	return &SafetyRepository{dbPool: dbPool}
}

func (safetyRepository *SafetyRepository) LoadSnapshot(ctx context.Context) (*CatalogSnapshot, error) {
	classes, classesErr := safetyRepository.loadClasses(ctx)
	if classesErr != nil {
		return nil, classesErr
	}
	entries, entriesErr := safetyRepository.loadEntries(ctx)
	if entriesErr != nil {
		return nil, entriesErr
	}
	interactions, interactionsErr := safetyRepository.loadInteractions(ctx)
	if interactionsErr != nil {
		return nil, interactionsErr
	}
	return NewCatalogSnapshot(classes, entries, interactions), nil
}

func (safetyRepository *SafetyRepository) loadClasses(ctx context.Context) ([]CatalogClass, error) {
	rows, queryErr := safetyRepository.dbPool.Query(ctx, `SELECT code, name, search_terms FROM medication_classes ORDER BY code`)
	if queryErr != nil {
		return nil, fmt.Errorf("failed to load medication classes: %w", queryErr)
	}
	defer rows.Close()
	var classes []CatalogClass
	for rows.Next() {
		var catalogClass CatalogClass
		if scanErr := rows.Scan(&catalogClass.Code, &catalogClass.Name, &catalogClass.SearchTerms); scanErr != nil {
			return nil, fmt.Errorf("failed to scan medication class: %w", scanErr)
		}
		classes = append(classes, catalogClass)
	}
	return classes, rows.Err()
}

func (safetyRepository *SafetyRepository) loadEntries(ctx context.Context) ([]CatalogEntry, error) {
	rows, queryErr := safetyRepository.dbPool.Query(ctx, `SELECT code, name, active_ingredient, COALESCE(class_code, ''), search_terms FROM medication_catalog ORDER BY code`)
	if queryErr != nil {
		return nil, fmt.Errorf("failed to load medication catalog: %w", queryErr)
	}
	defer rows.Close()
	var entries []CatalogEntry
	for rows.Next() {
		var entry CatalogEntry
		if scanErr := rows.Scan(&entry.Code, &entry.Name, &entry.ActiveIngredient, &entry.ClassCode, &entry.SearchTerms); scanErr != nil {
			return nil, fmt.Errorf("failed to scan medication catalog entry: %w", scanErr)
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func (safetyRepository *SafetyRepository) loadInteractions(ctx context.Context) ([]CatalogInteraction, error) {
	rows, queryErr := safetyRepository.dbPool.Query(ctx, `SELECT subject_a_type, subject_a, subject_b_type, subject_b, severity, description_key FROM medication_interactions ORDER BY created_at, id`)
	if queryErr != nil {
		return nil, fmt.Errorf("failed to load medication interactions: %w", queryErr)
	}
	defer rows.Close()
	var interactions []CatalogInteraction
	for rows.Next() {
		var interaction CatalogInteraction
		if scanErr := rows.Scan(&interaction.SubjectAType, &interaction.SubjectA, &interaction.SubjectBType, &interaction.SubjectB, &interaction.Severity, &interaction.DescriptionKey); scanErr != nil {
			return nil, fmt.Errorf("failed to scan medication interaction: %w", scanErr)
		}
		interactions = append(interactions, interaction)
	}
	return interactions, rows.Err()
}

func (safetyRepository *SafetyRepository) RecordOverride(ctx context.Context, override SafetyOverride) error {
	_, execErr := safetyRepository.dbPool.Exec(ctx,
		`INSERT INTO medication_safety_overrides (id, patient_fhir_id, medication_code, practitioner_user_id, alert_kinds, justification, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		override.ID, override.PatientFHIRID, override.MedicationCode, override.PractitionerUserID, override.AlertKinds, override.Justification, override.CreatedAt,
	)
	if execErr != nil {
		return fmt.Errorf("failed to record medication safety override: %w", execErr)
	}
	return nil
}

func (safetyRepository *SafetyRepository) LinkMedicationRequest(ctx context.Context, overrideID uuid.UUID, medicationRequestFHIRID string) error {
	_, execErr := safetyRepository.dbPool.Exec(ctx,
		`UPDATE medication_safety_overrides SET medication_request_fhir_id = $2 WHERE id = $1`,
		overrideID, medicationRequestFHIRID,
	)
	if execErr != nil {
		return fmt.Errorf("failed to link medication safety override: %w", execErr)
	}
	return nil
}

func (safetyRepository *SafetyRepository) DiscardOverride(ctx context.Context, overrideID uuid.UUID) error {
	_, execErr := safetyRepository.dbPool.Exec(ctx, `DELETE FROM medication_safety_overrides WHERE id = $1 AND medication_request_fhir_id IS NULL`, overrideID)
	if execErr != nil {
		return fmt.Errorf("failed to discard medication safety override: %w", execErr)
	}
	return nil
}
