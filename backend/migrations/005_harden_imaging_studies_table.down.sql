ALTER TABLE imaging_studies
DROP CONSTRAINT IF EXISTS imaging_studies_status_check;

DROP INDEX IF EXISTS idx_imaging_studies_status_created_at;

DROP INDEX IF EXISTS idx_imaging_studies_patient_created_at;
