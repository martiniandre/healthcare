package imaging

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	uniqueViolationCode = "23505"

	studyColumns = `id, patient_fhir_id, COALESCE(title, ''), COALESCE(modality, ''), file_name, study_instance_uid,
		study_datetime, COALESCE(accession_number, ''), COALESCE(study_description, ''), status, created_at, updated_at`
)

type Repository interface {
	CreateImagingStudy(ctx context.Context, study *ImagingStudy, instanceIndex *InstanceIndex) error
	GetImagingStudy(ctx context.Context, id uuid.UUID) (*ImagingStudy, error)
	FindImagingStudyByStudyInstanceUID(ctx context.Context, studyInstanceUID string) (*ImagingStudy, error)
	GetInstanceIndex(ctx context.Context, studyID uuid.UUID) (*InstanceIndex, error)
	ListImagingStudiesByPatient(ctx context.Context, patientFhirID string) ([]*ImagingStudy, error)
	UpdateImagingStudy(ctx context.Context, study *ImagingStudy) error
}

type repository struct {
	dbPool *pgxpool.Pool
}

func NewRepository(dbPool *pgxpool.Pool) Repository {
	return &repository{dbPool: dbPool}
}

func (repo *repository) CreateImagingStudy(ctx context.Context, study *ImagingStudy, instanceIndex *InstanceIndex) error {
	transaction, beginError := repo.dbPool.Begin(ctx)
	if beginError != nil {
		return beginError
	}
	defer transaction.Rollback(ctx)

	insertStudyStatement := `
		INSERT INTO imaging_studies (
			id, patient_fhir_id, title, modality, file_name, study_instance_uid,
			study_datetime, accession_number, study_description,
			series_count, total_instances, stored_instances,
			status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	seriesCount, instanceCount := 0, 0
	if instanceIndex != nil {
		seriesCount, instanceCount = 1, 1
	}
	_, insertStudyError := transaction.Exec(
		ctx,
		insertStudyStatement,
		study.ID,
		study.PatientFhirID,
		study.Title,
		study.Modality,
		study.FileName,
		study.StudyInstanceUID,
		study.StudyDateTime,
		nullableText(study.AccessionNumber),
		nullableText(study.StudyDescription),
		seriesCount,
		instanceCount,
		instanceCount,
		study.Status,
		study.CreatedAt,
		study.UpdatedAt,
	)
	if insertStudyError != nil {
		return mapInsertError(insertStudyError)
	}

	if instanceIndex != nil {
		if indexError := insertInstanceIndex(ctx, transaction, study, instanceIndex); indexError != nil {
			return mapInsertError(indexError)
		}
	}

	return transaction.Commit(ctx)
}

func insertInstanceIndex(ctx context.Context, transaction pgx.Tx, study *ImagingStudy, instanceIndex *InstanceIndex) error {
	seriesID := uuid.New()
	insertSeriesStatement := `
		INSERT INTO imaging_series (id, study_id, series_instance_uid, series_number, series_description, modality, instance_count)
		VALUES ($1, $2, $3, $4, $5, $6, 1)
	`
	_, seriesError := transaction.Exec(
		ctx,
		insertSeriesStatement,
		seriesID,
		study.ID,
		instanceIndex.SeriesInstanceUID,
		instanceIndex.SeriesNumber,
		nullableText(instanceIndex.SeriesDescription),
		nullableText(instanceIndex.Modality),
	)
	if seriesError != nil {
		return seriesError
	}

	insertInstanceStatement := `
		INSERT INTO imaging_instances (
			id, study_id, series_id, sop_instance_uid, sop_class_uid, transfer_syntax_uid,
			instance_number, frame_count, rows_count, columns_count,
			pixel_spacing_row, pixel_spacing_column, ingest_status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, 'STORED')
	`
	_, instanceError := transaction.Exec(
		ctx,
		insertInstanceStatement,
		uuid.New(),
		study.ID,
		seriesID,
		instanceIndex.SOPInstanceUID,
		instanceIndex.SOPClassUID,
		instanceIndex.TransferSyntaxUID,
		instanceIndex.InstanceNumber,
		instanceIndex.FrameCount,
		nullablePositiveInt(instanceIndex.Rows),
		nullablePositiveInt(instanceIndex.Columns),
		nullablePositiveFloat(instanceIndex.PixelSpacingRow),
		nullablePositiveFloat(instanceIndex.PixelSpacingColumn),
	)
	return instanceError
}

func mapInsertError(insertError error) error {
	var postgresError *pgconn.PgError
	if errors.As(insertError, &postgresError) && postgresError.Code == uniqueViolationCode {
		return apperrors.ErrDICOMDuplicateInstance
	}
	return insertError
}

func nullableText(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func nullablePositiveInt(value int) *int {
	if value <= 0 {
		return nil
	}
	return &value
}

func nullablePositiveFloat(value float64) *float64 {
	if value <= 0 {
		return nil
	}
	return &value
}

func scanStudy(rowResult pgx.Row) (*ImagingStudy, error) {
	study := &ImagingStudy{}
	scanError := rowResult.Scan(
		&study.ID,
		&study.PatientFhirID,
		&study.Title,
		&study.Modality,
		&study.FileName,
		&study.StudyInstanceUID,
		&study.StudyDateTime,
		&study.AccessionNumber,
		&study.StudyDescription,
		&study.Status,
		&study.CreatedAt,
		&study.UpdatedAt,
	)
	if scanError != nil {
		return nil, scanError
	}
	return study, nil
}

func (repo *repository) GetImagingStudy(ctx context.Context, id uuid.UUID) (*ImagingStudy, error) {
	queryStatement := `SELECT ` + studyColumns + ` FROM imaging_studies WHERE id = $1`
	study, scanError := scanStudy(repo.dbPool.QueryRow(ctx, queryStatement, id))
	if scanError != nil {
		if errors.Is(scanError, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to get imaging study: %w", apperrors.ErrImagingStudyNotFound)
		}
		return nil, scanError
	}
	return study, nil
}

func (repo *repository) FindImagingStudyByStudyInstanceUID(ctx context.Context, studyInstanceUID string) (*ImagingStudy, error) {
	queryStatement := `SELECT ` + studyColumns + ` FROM imaging_studies WHERE study_instance_uid = $1`
	study, scanError := scanStudy(repo.dbPool.QueryRow(ctx, queryStatement, studyInstanceUID))
	if scanError != nil {
		if errors.Is(scanError, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to find imaging study: %w", apperrors.ErrImagingStudyNotFound)
		}
		return nil, scanError
	}
	return study, nil
}

func (repo *repository) GetInstanceIndex(ctx context.Context, studyID uuid.UUID) (*InstanceIndex, error) {
	queryStatement := `
		SELECT series.series_instance_uid, COALESCE(series.series_number, 0), COALESCE(series.series_description, ''),
			COALESCE(series.modality, ''), instance.sop_instance_uid, instance.sop_class_uid, instance.transfer_syntax_uid,
			COALESCE(instance.instance_number, 0), instance.frame_count,
			COALESCE(instance.rows_count, 0), COALESCE(instance.columns_count, 0),
			COALESCE(instance.pixel_spacing_row, 0), COALESCE(instance.pixel_spacing_column, 0)
		FROM imaging_instances instance
		JOIN imaging_series series ON series.id = instance.series_id
		WHERE instance.study_id = $1
		ORDER BY instance.created_at
		LIMIT 1
	`
	instanceIndex := &InstanceIndex{}
	scanError := repo.dbPool.QueryRow(ctx, queryStatement, studyID).Scan(
		&instanceIndex.SeriesInstanceUID,
		&instanceIndex.SeriesNumber,
		&instanceIndex.SeriesDescription,
		&instanceIndex.Modality,
		&instanceIndex.SOPInstanceUID,
		&instanceIndex.SOPClassUID,
		&instanceIndex.TransferSyntaxUID,
		&instanceIndex.InstanceNumber,
		&instanceIndex.FrameCount,
		&instanceIndex.Rows,
		&instanceIndex.Columns,
		&instanceIndex.PixelSpacingRow,
		&instanceIndex.PixelSpacingColumn,
	)
	if scanError != nil {
		if errors.Is(scanError, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to get instance index: %w", apperrors.ErrImagingStudyNotFound)
		}
		return nil, scanError
	}
	return instanceIndex, nil
}

func (repo *repository) ListImagingStudiesByPatient(ctx context.Context, patientFhirID string) ([]*ImagingStudy, error) {
	queryStatement := `SELECT ` + studyColumns + ` FROM imaging_studies WHERE patient_fhir_id = $1 ORDER BY created_at DESC`
	rowsResult, queryError := repo.dbPool.Query(ctx, queryStatement, patientFhirID)
	if queryError != nil {
		return nil, queryError
	}
	defer rowsResult.Close()

	var studiesList []*ImagingStudy
	for rowsResult.Next() {
		study, scanError := scanStudy(rowsResult)
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
		SET patient_fhir_id = $2, title = $3, modality = $4, file_name = $5, status = $6::varchar, updated_at = $7,
			ready_at = CASE WHEN $6::varchar = 'READY' THEN $7 ELSE ready_at END
		WHERE id = $1
	`
	_, executionError := repo.dbPool.Exec(
		ctx,
		queryStatement,
		study.ID,
		study.PatientFhirID,
		study.Title,
		study.Modality,
		study.FileName,
		study.Status,
		study.UpdatedAt,
	)
	return executionError
}
