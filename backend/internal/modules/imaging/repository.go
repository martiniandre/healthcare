package imaging

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateImagingStudy(ctx context.Context, study *ImagingStudy) error
	GetImagingStudy(ctx context.Context, id uuid.UUID) (*ImagingStudy, error)
	ListImagingStudiesByPatient(ctx context.Context, patientFhirID string) ([]*ImagingStudy, error)
	UpdateImagingStudy(ctx context.Context, study *ImagingStudy) error
}

type repository struct {
	dbPool *pgxpool.Pool
}

func NewRepository(dbPool *pgxpool.Pool) Repository {
	return &repository{dbPool: dbPool}
}

func (repo *repository) CreateImagingStudy(ctx context.Context, study *ImagingStudy) error {
	queryStatement := `
		INSERT INTO imaging_studies (id, patient_fhir_id, title, modality, gcs_path, study_instance_uid, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, executionError := repo.dbPool.Exec(
		ctx,
		queryStatement,
		study.ID,
		study.PatientFhirID,
		study.Title,
		study.Modality,
		study.GCSPath,
		study.StudyInstanceUID,
		study.Status,
		study.CreatedAt,
		study.UpdatedAt,
	)
	return executionError
}

func (repo *repository) GetImagingStudy(ctx context.Context, id uuid.UUID) (*ImagingStudy, error) {
	queryStatement := `
		SELECT id, patient_fhir_id, title, modality, gcs_path, study_instance_uid, status, created_at, updated_at
		FROM imaging_studies
		WHERE id = $1
	`
	rowResult := repo.dbPool.QueryRow(ctx, queryStatement, id)

	study := &ImagingStudy{}
	scanError := rowResult.Scan(
		&study.ID,
		&study.PatientFhirID,
		&study.Title,
		&study.Modality,
		&study.GCSPath,
		&study.StudyInstanceUID,
		&study.Status,
		&study.CreatedAt,
		&study.UpdatedAt,
	)
	if scanError != nil {
		if errors.Is(scanError, pgx.ErrNoRows) {
			return nil, ErrImagingStudyNotFound
		}
		return nil, scanError
	}
	return study, nil
}

func (repo *repository) ListImagingStudiesByPatient(ctx context.Context, patientFhirID string) ([]*ImagingStudy, error) {
	queryStatement := `
		SELECT id, patient_fhir_id, title, modality, gcs_path, study_instance_uid, status, created_at, updated_at
		FROM imaging_studies
		WHERE patient_fhir_id = $1
		ORDER BY created_at DESC
	`
	rowsResult, queryError := repo.dbPool.Query(ctx, queryStatement, patientFhirID)
	if queryError != nil {
		return nil, queryError
	}
	defer rowsResult.Close()

	var studiesList []*ImagingStudy
	for rowsResult.Next() {
		study := &ImagingStudy{}
		scanError := rowsResult.Scan(
			&study.ID,
			&study.PatientFhirID,
			&study.Title,
			&study.Modality,
			&study.GCSPath,
			&study.StudyInstanceUID,
			&study.Status,
			&study.CreatedAt,
			&study.UpdatedAt,
		)
		if scanError != nil {
			return nil, scanError
		}
		studiesList = append(studiesList, study)
	}
	if rowsError := rowsResult.Err(); rowsError != nil {
		return nil, rowsError
	}
	return studiesList, nil
}

func (repo *repository) UpdateImagingStudy(ctx context.Context, study *ImagingStudy) error {
	queryStatement := `
		UPDATE imaging_studies
		SET patient_fhir_id = $2, title = $3, modality = $4, gcs_path = $5, study_instance_uid = $6, status = $7, updated_at = $8
		WHERE id = $1
	`
	_, executionError := repo.dbPool.Exec(
		ctx,
		queryStatement,
		study.ID,
		study.PatientFhirID,
		study.Title,
		study.Modality,
		study.GCSPath,
		study.StudyInstanceUID,
		study.Status,
		study.UpdatedAt,
	)
	return executionError
}
