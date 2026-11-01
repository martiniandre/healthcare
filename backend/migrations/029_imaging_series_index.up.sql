CREATE TABLE IF NOT EXISTS imaging_series (
    id UUID PRIMARY KEY,
    study_id UUID NOT NULL REFERENCES imaging_studies(id) ON DELETE CASCADE,
    series_instance_uid VARCHAR(128) NOT NULL,
    series_number INTEGER,
    series_description VARCHAR(255),
    modality VARCHAR(32),
    instance_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_imaging_series_study_series_uid UNIQUE (study_id, series_instance_uid)
);

CREATE INDEX IF NOT EXISTS idx_imaging_series_study_id
ON imaging_series(study_id);

CREATE INDEX IF NOT EXISTS idx_imaging_series_modality
ON imaging_series(modality);
