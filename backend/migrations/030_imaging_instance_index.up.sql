CREATE TABLE IF NOT EXISTS imaging_instances (
    id UUID PRIMARY KEY,
    study_id UUID NOT NULL REFERENCES imaging_studies(id) ON DELETE CASCADE,
    series_id UUID NOT NULL REFERENCES imaging_series(id) ON DELETE CASCADE,
    sop_instance_uid VARCHAR(128) NOT NULL,
    sop_class_uid VARCHAR(128) NOT NULL,
    transfer_syntax_uid VARCHAR(128) NOT NULL,
    instance_number INTEGER,
    frame_count INTEGER NOT NULL DEFAULT 1,
    rows_count INTEGER,
    columns_count INTEGER,
    pixel_spacing_row DOUBLE PRECISION,
    pixel_spacing_column DOUBLE PRECISION,
    ingest_status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
    ingest_error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_imaging_instances_study_sop_uid UNIQUE (study_id, sop_instance_uid),
    CONSTRAINT imaging_instances_ingest_status_check
        CHECK (ingest_status IN ('PENDING', 'STORED', 'ALREADY_EXISTS', 'FAILED')),
    CONSTRAINT imaging_instances_geometry_check
        CHECK (
            (rows_count IS NULL OR rows_count > 0)
            AND (columns_count IS NULL OR columns_count > 0)
            AND (frame_count IS NULL OR frame_count > 0)
        )
);

CREATE INDEX IF NOT EXISTS idx_imaging_instances_series_id
ON imaging_instances(series_id);

CREATE INDEX IF NOT EXISTS idx_imaging_instances_study_id
ON imaging_instances(study_id);

CREATE INDEX IF NOT EXISTS idx_imaging_instances_ingest_status
ON imaging_instances(ingest_status);
