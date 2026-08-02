package diagnostic_report

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VersionRepository interface {
	RecordVersion(ctx context.Context, reportID string, version string, snapshot json.RawMessage, changedBy *uuid.UUID) (*DiagnosticReportVersion, error)
	ListVersions(ctx context.Context, reportID string) ([]*DiagnosticReportVersion, error)
	GetVersion(ctx context.Context, reportID string, version string) (*DiagnosticReportVersion, error)
}

type versionRepository struct {
	db *pgxpool.Pool
}

func NewVersionRepository(db *pgxpool.Pool) VersionRepository {
	return &versionRepository{db: db}
}

func (versionRepo *versionRepository) RecordVersion(ctx context.Context, reportID string, version string, snapshot json.RawMessage, changedBy *uuid.UUID) (*DiagnosticReportVersion, error) {
	entry := &DiagnosticReportVersion{
		ID:        uuid.New(),
		ReportID:  reportID,
		Version:   version,
		Snapshot:  snapshot,
		ChangedBy: changedBy,
	}

	insertQuery := `INSERT INTO diagnostic_report_versions (id, report_id, version, snapshot, changed_by)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := versionRepo.db.Exec(ctx, insertQuery, entry.ID, entry.ReportID, entry.Version, entry.Snapshot, entry.ChangedBy)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (versionRepo *versionRepository) ListVersions(ctx context.Context, reportID string) ([]*DiagnosticReportVersion, error) {
	listQuery := `SELECT id, report_id, version, snapshot, changed_by, changed_at
		FROM diagnostic_report_versions
		WHERE report_id = $1
		ORDER BY changed_at DESC`

	rows, err := versionRepo.db.Query(ctx, listQuery, reportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	versions := make([]*DiagnosticReportVersion, 0)
	for rows.Next() {
		var entry DiagnosticReportVersion
		if scanErr := rows.Scan(&entry.ID, &entry.ReportID, &entry.Version, &entry.Snapshot, &entry.ChangedBy, &entry.ChangedAt); scanErr != nil {
			return nil, scanErr
		}
		versions = append(versions, &entry)
	}
	return versions, nil
}

func (versionRepo *versionRepository) GetVersion(ctx context.Context, reportID string, version string) (*DiagnosticReportVersion, error) {
	getQuery := `SELECT id, report_id, version, snapshot, changed_by, changed_at
		FROM diagnostic_report_versions
		WHERE report_id = $1 AND version = $2`

	var entry DiagnosticReportVersion
	err := versionRepo.db.QueryRow(ctx, getQuery, reportID, version).Scan(&entry.ID, &entry.ReportID, &entry.Version, &entry.Snapshot, &entry.ChangedBy, &entry.ChangedAt)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}
