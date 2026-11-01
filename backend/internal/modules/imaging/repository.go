package imaging

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	uniqueViolationCode = "23505"

	studyColumns = `id, patient_fhir_id, COALESCE(title, ''), COALESCE(modality, ''), file_name, study_instance_uid,
		study_datetime, COALESCE(accession_number, ''), COALESCE(study_description, ''), COALESCE(fhir_resource_id, ''),
		status, created_at, updated_at`

	instanceColumns = `instance.id, instance.study_id, series.series_instance_uid, COALESCE(series.series_number, 0),
		COALESCE(series.series_description, ''), COALESCE(series.modality, ''), instance.sop_instance_uid,
		instance.sop_class_uid, instance.transfer_syntax_uid, COALESCE(instance.instance_number, 0), instance.frame_count,
		COALESCE(instance.rows_count, 0), COALESCE(instance.columns_count, 0), COALESCE(instance.pixel_spacing_row, 0),
		COALESCE(instance.pixel_spacing_column, 0), COALESCE(instance.file_name, ''), COALESCE(instance.object_key, ''),
		instance.created_at`
)

type Repository interface {
	CreateImagingStudy(ctx context.Context, study *ImagingStudy, instances []*InstanceIndex) error
	AppendInstances(ctx context.Context, studyID uuid.UUID, instances []*InstanceIndex, updatedAt time.Time) error
	GetImagingStudy(ctx context.Context, id uuid.UUID) (*ImagingStudy, error)
	FindImagingStudyByStudyInstanceUID(ctx context.Context, studyInstanceUID string) (*ImagingStudy, error)
	HasInstance(ctx context.Context, studyID uuid.UUID, sopInstanceUID string) (bool, error)
	ListInstances(ctx context.Context, studyID uuid.UUID) ([]*InstanceIndex, error)
	ListImagingStudiesByPatient(ctx context.Context, patientFhirID string) ([]*ImagingStudy, error)
	UpdateImagingStudy(ctx context.Context, study *ImagingStudy) error
}

type repository struct {
	dbPool *pgxpool.Pool
}

func NewRepository(dbPool *pgxpool.Pool) Repository {
	return &repository{dbPool: dbPool}
}

func (repo *repository) CreateImagingStudy(ctx context.Context, study *ImagingStudy, instances []*InstanceIndex) error {
	transaction, beginError := repo.dbPool.Begin(ctx)
	if beginError != nil {
		return beginError
	}
	defer transaction.Rollback(ctx)

	insertStudyStatement := `
		INSERT INTO imaging_studies (
			id, patient_fhir_id, title, modality, file_name, study_instance_uid,
			study_datetime, accession_number, study_description, status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
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
		study.Status,
		study.CreatedAt,
		study.UpdatedAt,
	)
	if insertStudyError != nil {
		return mapInsertError(insertStudyError)
	}

	if insertError := insertInstances(ctx, transaction, study.ID, instances); insertError != nil {
		return mapInsertError(insertError)
	}
	if counterError := refreshStudyCounters(ctx, transaction, study.ID); counterError != nil {
		return counterError
	}

	return transaction.Commit(ctx)
}

func (repo *repository) AppendInstances(ctx context.Context, studyID uuid.UUID, instances []*InstanceIndex, updatedAt time.Time) error {
	transaction, beginError := repo.dbPool.Begin(ctx)
	if beginError != nil {
		return beginError
	}
	defer transaction.Rollback(ctx)

	var lockedStudyID uuid.UUID
	lockError := transaction.QueryRow(ctx, `SELECT id FROM imaging_studies WHERE id = $1 FOR UPDATE`, studyID).Scan(&lockedStudyID)
	if lockError != nil {
		if errors.Is(lockError, pgx.ErrNoRows) {
			return fmt.Errorf("failed to append imaging instances: %w", apperrors.ErrImagingStudyNotFound)
		}
		return lockError
	}

	if insertError := insertInstances(ctx, transaction, studyID, instances); insertError != nil {
		return mapInsertError(insertError)
	}
	if counterError := refreshStudyCounters(ctx, transaction, studyID); counterError != nil {
		return counterError
	}

	_, statusError := transaction.Exec(
		ctx,
		`UPDATE imaging_studies SET status = $2::varchar, updated_at = $3 WHERE id = $1`,
		studyID,
		ImagingStudyStatusReceiving,
		updatedAt,
	)
	if statusError != nil {
		return statusError
	}

	return transaction.Commit(ctx)
}

func insertInstances(ctx context.Context, transaction pgx.Tx, studyID uuid.UUID, instances []*InstanceIndex) error {
	seriesIDsByUID := make(map[string]uuid.UUID)
	for _, instance := range instances {
		seriesID, knownSeries := seriesIDsByUID[instance.SeriesInstanceUID]
		if !knownSeries {
			upsertedSeriesID, seriesError := upsertSeries(ctx, transaction, studyID, instance)
			if seriesError != nil {
				return seriesError
			}
			seriesID = upsertedSeriesID
			seriesIDsByUID[instance.SeriesInstanceUID] = seriesID
		}

		insertInstanceStatement := `
			INSERT INTO imaging_instances (
				id, study_id, series_id, sop_instance_uid, sop_class_uid, transfer_syntax_uid,
				instance_number, frame_count, rows_count, columns_count,
				pixel_spacing_row, pixel_spacing_column, file_name, object_key, ingest_status
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, 'STORED')
		`
		_, instanceError := transaction.Exec(
			ctx,
			insertInstanceStatement,
			instance.ID,
			studyID,
			seriesID,
			instance.SOPInstanceUID,
			instance.SOPClassUID,
			instance.TransferSyntaxUID,
			instance.InstanceNumber,
			instance.FrameCount,
			nullablePositiveInt(instance.Rows),
			nullablePositiveInt(instance.Columns),
			nullablePositiveFloat(instance.PixelSpacingRow),
			nullablePositiveFloat(instance.PixelSpacingColumn),
			nullableText(instance.FileName),
			nullableText(instance.ObjectKey),
		)
		if instanceError != nil {
			return instanceError
		}
	}
	return nil
}

func upsertSeries(ctx context.Context, transaction pgx.Tx, studyID uuid.UUID, instance *InstanceIndex) (uuid.UUID, error) {
	upsertStatement := `
		INSERT INTO imaging_series (id, study_id, series_instance_uid, series_number, series_description, modality, instance_count)
		VALUES ($1, $2, $3, $4, $5, $6, 0)
		ON CONFLICT (study_id, series_instance_uid) DO UPDATE SET updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`
	var seriesID uuid.UUID
	scanError := transaction.QueryRow(
		ctx,
		upsertStatement,
		uuid.New(),
		studyID,
		instance.SeriesInstanceUID,
		instance.SeriesNumber,
		nullableText(instance.SeriesDescription),
		nullableText(instance.Modality),
	).Scan(&seriesID)
	return seriesID, scanError
}

func refreshStudyCounters(ctx context.Context, transaction pgx.Tx, studyID uuid.UUID) error {
	_, seriesCountError := transaction.Exec(ctx, `
		UPDATE imaging_series series
		SET instance_count = (SELECT COUNT(*) FROM imaging_instances instance WHERE instance.series_id = series.id),
			updated_at = CURRENT_TIMESTAMP
		WHERE series.study_id = $1
	`, studyID)
	if seriesCountError != nil {
		return seriesCountError
	}

	_, studyCountError := transaction.Exec(ctx, `
		UPDATE imaging_studies
		SET series_count = (SELECT COUNT(*) FROM imaging_series WHERE study_id = $1),
			total_instances = (SELECT COUNT(*) FROM imaging_instances WHERE study_id = $1),
			stored_instances = (SELECT COUNT(*) FROM imaging_instances WHERE study_id = $1 AND ingest_status = 'STORED')
		WHERE id = $1
	`, studyID)
	return studyCountError
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
		&study.FHIRResourceID,
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

func (repo *repository) HasInstance(ctx context.Context, studyID uuid.UUID, sopInstanceUID string) (bool, error) {
	var instanceExists bool
	queryError := repo.dbPool.QueryRow(
		ctx,
		`SELECT EXISTS (SELECT 1 FROM imaging_instances WHERE study_id = $1 AND sop_instance_uid = $2)`,
		studyID,
		sopInstanceUID,
	).Scan(&instanceExists)
	return instanceExists, queryError
}

func (repo *repository) ListInstances(ctx context.Context, studyID uuid.UUID) ([]*InstanceIndex, error) {
	queryStatement := `
		SELECT ` + instanceColumns + `
		FROM imaging_instances instance
		JOIN imaging_series series ON series.id = instance.series_id
		WHERE instance.study_id = $1
		ORDER BY instance.created_at, COALESCE(instance.instance_number, 0), instance.file_name, instance.id
	`
	rowsResult, queryError := repo.dbPool.Query(ctx, queryStatement, studyID)
	if queryError != nil {
		return nil, queryError
	}
	defer rowsResult.Close()

	var instances []*InstanceIndex
	for rowsResult.Next() {
		instance := &InstanceIndex{}
		scanError := rowsResult.Scan(
			&instance.ID,
			&instance.StudyID,
			&instance.SeriesInstanceUID,
			&instance.SeriesNumber,
			&instance.SeriesDescription,
			&instance.Modality,
			&instance.SOPInstanceUID,
			&instance.SOPClassUID,
			&instance.TransferSyntaxUID,
			&instance.InstanceNumber,
			&instance.FrameCount,
			&instance.Rows,
			&instance.Columns,
			&instance.PixelSpacingRow,
			&instance.PixelSpacingColumn,
			&instance.FileName,
			&instance.ObjectKey,
			&instance.CreatedAt,
		)
		if scanError != nil {
			return nil, scanError
		}
		instances = append(instances, instance)
	}
	if rowsError := rowsResult.Err(); rowsError != nil {
		return nil, rowsError
	}
	return instances, nil
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
			fhir_resource_id = $8,
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
		nullableText(study.FHIRResourceID),
	)
	return executionError
}
