CREATE INDEX IF NOT EXISTS idx_imaging_studies_patient_created_at
ON imaging_studies(patient_fhir_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_imaging_studies_status_created_at
ON imaging_studies(status, created_at);

ALTER TABLE imaging_studies
ADD CONSTRAINT imaging_studies_status_check
CHECK (status IN ('PENDING', 'PROCESSING', 'PROCESSED', 'FAILED'));
