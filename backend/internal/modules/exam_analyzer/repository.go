package exam_analyzer

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAnalysisNotFound = errors.New("exam analysis not found")

type Repository interface {
	CreateAnalysis(ctx context.Context, analysis *ExamAnalysis) error
	GetAnalysis(ctx context.Context, id uuid.UUID) (*ExamAnalysis, error)
	ListAnalyses(ctx context.Context, patientFhirID *string) ([]*ExamAnalysis, error)
	UpdateAnalysis(ctx context.Context, analysis *ExamAnalysis) error
	DeleteAnalysis(ctx context.Context, id uuid.UUID) error
	CreateAuditLog(ctx context.Context, log *ExamAnalysisAuditLog) error
}

type repository struct {
	databasePool *pgxpool.Pool
}

func NewRepository(databasePool *pgxpool.Pool) Repository {
	return &repository{databasePool: databasePool}
}

func (repo *repository) CreateAnalysis(ctx context.Context, analysis *ExamAnalysis) error {
	queryStatement := `
		INSERT INTO exam_analyses (
			id, user_id, patient_fhir_id, exam_type, file_name, file_path, status, 
			analysis_response, consent_given, anonymized, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, executionError := repo.databasePool.Exec(
		ctx,
		queryStatement,
		analysis.ID,
		analysis.UserID,
		analysis.PatientFhirID,
		analysis.ExamType,
		analysis.FileName,
		analysis.FilePath,
		analysis.Status,
		analysis.AnalysisResponse,
		analysis.ConsentGiven,
		analysis.Anonymized,
		analysis.CreatedAt,
		analysis.UpdatedAt,
	)
	return executionError
}

func (repo *repository) GetAnalysis(ctx context.Context, id uuid.UUID) (*ExamAnalysis, error) {
	queryStatement := `
		SELECT 
			id, user_id, patient_fhir_id, exam_type, file_name, file_path, status, 
			analysis_response, consent_given, anonymized, created_at, updated_at
		FROM exam_analyses
		WHERE id = $1
	`
	rowResult := repo.databasePool.QueryRow(ctx, queryStatement, id)

	analysis := &ExamAnalysis{}
	scanError := rowResult.Scan(
		&analysis.ID,
		&analysis.UserID,
		&analysis.PatientFhirID,
		&analysis.ExamType,
		&analysis.FileName,
		&analysis.FilePath,
		&analysis.Status,
		&analysis.AnalysisResponse,
		&analysis.ConsentGiven,
		&analysis.Anonymized,
		&analysis.CreatedAt,
		&analysis.UpdatedAt,
	)
	if scanError != nil {
		if errors.Is(scanError, pgx.ErrNoRows) {
			return nil, ErrAnalysisNotFound
		}
		return nil, scanError
	}
	return analysis, nil
}

func (repo *repository) ListAnalyses(ctx context.Context, patientFhirID *string) ([]*ExamAnalysis, error) {
	var queryStatement string
	var rowsResult pgx.Rows
	var queryError error

	if patientFhirID != nil && *patientFhirID != "" {
		queryStatement = `
			SELECT 
				id, user_id, patient_fhir_id, exam_type, file_name, file_path, status, 
				analysis_response, consent_given, anonymized, created_at, updated_at
			FROM exam_analyses
			WHERE patient_fhir_id = $1
			ORDER BY created_at DESC
		`
		rowsResult, queryError = repo.databasePool.Query(ctx, queryStatement, *patientFhirID)
	} else {
		queryStatement = `
			SELECT 
				id, user_id, patient_fhir_id, exam_type, file_name, file_path, status, 
				analysis_response, consent_given, anonymized, created_at, updated_at
			FROM exam_analyses
			ORDER BY created_at DESC
		`
		rowsResult, queryError = repo.databasePool.Query(ctx, queryStatement)
	}

	if queryError != nil {
		return nil, queryError
	}
	defer rowsResult.Close()

	var analysesList []*ExamAnalysis
	for rowsResult.Next() {
		analysis := &ExamAnalysis{}
		scanError := rowsResult.Scan(
			&analysis.ID,
			&analysis.UserID,
			&analysis.PatientFhirID,
			&analysis.ExamType,
			&analysis.FileName,
			&analysis.FilePath,
			&analysis.Status,
			&analysis.AnalysisResponse,
			&analysis.ConsentGiven,
			&analysis.Anonymized,
			&analysis.CreatedAt,
			&analysis.UpdatedAt,
		)
		if scanError != nil {
			return nil, scanError
		}
		analysesList = append(analysesList, analysis)
	}

	if rowsError := rowsResult.Err(); rowsError != nil {
		return nil, rowsError
	}
	return analysesList, nil
}

func (repo *repository) UpdateAnalysis(ctx context.Context, analysis *ExamAnalysis) error {
	queryStatement := `
		UPDATE exam_analyses
		SET user_id = $2, patient_fhir_id = $3, exam_type = $4, file_name = $5, file_path = $6, 
			status = $7, analysis_response = $8, consent_given = $9, anonymized = $10, updated_at = $11
		WHERE id = $1
	`
	_, executionError := repo.databasePool.Exec(
		ctx,
		queryStatement,
		analysis.ID,
		analysis.UserID,
		analysis.PatientFhirID,
		analysis.ExamType,
		analysis.FileName,
		analysis.FilePath,
		analysis.Status,
		analysis.AnalysisResponse,
		analysis.ConsentGiven,
		analysis.Anonymized,
		analysis.UpdatedAt,
	)
	return executionError
}

func (repo *repository) DeleteAnalysis(ctx context.Context, id uuid.UUID) error {
	queryStatement := `
		DELETE FROM exam_analyses
		WHERE id = $1
	`
	_, executionError := repo.databasePool.Exec(ctx, queryStatement, id)
	return executionError
}

func (repo *repository) CreateAuditLog(ctx context.Context, log *ExamAnalysisAuditLog) error {
	queryStatement := `
		INSERT INTO exam_analyses_audit_logs (
			id, analysis_id, action_type, performed_by, ip_address, details, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, executionError := repo.databasePool.Exec(
		ctx,
		queryStatement,
		log.ID,
		log.AnalysisID,
		log.ActionType,
		log.PerformedBy,
		log.IPAddress,
		log.Details,
		log.CreatedAt,
	)
	return executionError
}
