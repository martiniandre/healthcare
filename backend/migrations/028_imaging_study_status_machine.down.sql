DROP INDEX IF EXISTS idx_imaging_studies_status_ready_at;

DROP INDEX IF EXISTS uq_imaging_studies_study_instance_uid;

ALTER TABLE imaging_studies
ALTER COLUMN study_instance_uid DROP NOT NULL;

ALTER TABLE imaging_studies
DROP CONSTRAINT IF EXISTS imaging_studies_counters_check;

ALTER TABLE imaging_studies
DROP COLUMN IF EXISTS ready_at;

ALTER TABLE imaging_studies
DROP COLUMN IF EXISTS ingest_error;

ALTER TABLE imaging_studies
DROP COLUMN IF EXISTS failed_instances;

ALTER TABLE imaging_studies
DROP COLUMN IF EXISTS skipped_instances;

ALTER TABLE imaging_studies
DROP COLUMN IF EXISTS stored_instances;

ALTER TABLE imaging_studies
DROP COLUMN IF EXISTS total_instances;

ALTER TABLE imaging_studies
DROP COLUMN IF EXISTS series_count;

ALTER TABLE imaging_studies
DROP COLUMN IF EXISTS study_description;

ALTER TABLE imaging_studies
DROP COLUMN IF EXISTS accession_number;

ALTER TABLE imaging_studies
DROP COLUMN IF EXISTS study_datetime;

ALTER TABLE imaging_studies
ALTER COLUMN modality SET NOT NULL;

ALTER TABLE imaging_studies
ALTER COLUMN title SET NOT NULL;

DELETE FROM imaging_studies
WHERE modality IS NULL OR title IS NULL;

UPDATE imaging_studies
SET status = 'PROCESSED'
WHERE status = 'READY';
UPDATE imaging_studies
SET status = 'PENDING'
WHERE status = 'RECEIVING';

UPDATE imaging_studies
SET status = 'FAILED'
WHERE status = 'PARTIAL';

ALTER TABLE imaging_studies
DROP CONSTRAINT IF EXISTS imaging_studies_status_check;

ALTER TABLE imaging_studies
ADD CONSTRAINT imaging_studies_status_check
CHECK (status IN ('PENDING', 'PROCESSING', 'PROCESSED', 'FAILED'));
