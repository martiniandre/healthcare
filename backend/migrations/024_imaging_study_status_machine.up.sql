UPDATE imaging_studies
SET status = 'READY'
WHERE status = 'PROCESSED';

UPDATE imaging_studies
SET status = 'RECEIVING'
WHERE status = 'PENDING';

ALTER TABLE imaging_studies
DROP CONSTRAINT IF EXISTS imaging_studies_status_check;

ALTER TABLE imaging_studies
ADD CONSTRAINT imaging_studies_status_check
CHECK (status IN ('RECEIVING', 'PROCESSING', 'READY', 'PARTIAL', 'FAILED'));

ALTER TABLE imaging_studies
ALTER COLUMN title DROP NOT NULL;

ALTER TABLE imaging_studies
ALTER COLUMN modality DROP NOT NULL;

ALTER TABLE imaging_studies
ADD COLUMN study_date DATE;

ALTER TABLE imaging_studies
ADD COLUMN study_time TIME;

ALTER TABLE imaging_studies
ADD COLUMN accession_number VARCHAR(64);

ALTER TABLE imaging_studies
ADD COLUMN study_description VARCHAR(255);

ALTER TABLE imaging_studies
ADD COLUMN series_count INTEGER NOT NULL DEFAULT 0;

ALTER TABLE imaging_studies
ADD COLUMN total_instances INTEGER NOT NULL DEFAULT 0;

ALTER TABLE imaging_studies
ADD COLUMN stored_instances INTEGER NOT NULL DEFAULT 0;

ALTER TABLE imaging_studies
ADD COLUMN skipped_instances INTEGER NOT NULL DEFAULT 0;

ALTER TABLE imaging_studies
ADD COLUMN failed_instances INTEGER NOT NULL DEFAULT 0;

ALTER TABLE imaging_studies
ADD COLUMN ingest_error TEXT;

ALTER TABLE imaging_studies
ADD COLUMN ready_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE imaging_studies
ADD CONSTRAINT imaging_studies_counters_check
CHECK (
    series_count >= 0
    AND total_instances >= 0
    AND stored_instances >= 0
    AND skipped_instances >= 0
    AND failed_instances >= 0
);

DELETE FROM imaging_studies
WHERE study_instance_uid IS NULL;

ALTER TABLE imaging_studies
ALTER COLUMN study_instance_uid SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_imaging_studies_study_instance_uid
ON imaging_studies(study_instance_uid);

CREATE INDEX IF NOT EXISTS idx_imaging_studies_status_ready_at
ON imaging_studies(status, ready_at);
