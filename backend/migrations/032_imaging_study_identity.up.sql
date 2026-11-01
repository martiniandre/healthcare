-- Migration 032: imaging study identity and per-instance object keys.
--
-- Prerequisite (production, manual check, BEFORE migration 028 runs there):
--   SELECT study_instance_uid, COUNT(*)
--   FROM imaging_studies
--   WHERE study_instance_uid IS NOT NULL AND study_instance_uid <> ''
--   GROUP BY study_instance_uid
--   HAVING COUNT(*) > 1;
-- The unique index created by 028 is not partial, so any duplicate non-empty
-- study_instance_uid makes 028 fail. Migration 028 is intentionally not edited.
--
-- Study identity: the unique index becomes partial so that legacy rows with an
-- empty study_instance_uid never collide. Non-empty UIDs stay globally unique.
--
-- Instances: each stored DICOM instance records its own object key (and the
-- display file name). Legacy single-file studies are backfilled with the key
-- they were stored under: pacs/{patient}/{studyId}{ext}.
--
-- Counters: series/instance counters are recomputed from the real rows instead
-- of the historical hard-coded value of 1.

DROP INDEX IF EXISTS uq_imaging_studies_study_instance_uid;

CREATE UNIQUE INDEX uq_imaging_studies_study_instance_uid
ON imaging_studies(study_instance_uid)
WHERE study_instance_uid <> '';

ALTER TABLE imaging_studies
ADD COLUMN IF NOT EXISTS fhir_resource_id VARCHAR(128);

ALTER TABLE imaging_instances
ADD COLUMN IF NOT EXISTS file_name VARCHAR(255);

ALTER TABLE imaging_instances
ADD COLUMN IF NOT EXISTS object_key VARCHAR(512);

UPDATE imaging_instances instance
SET object_key = 'pacs/'
        || regexp_replace(study.patient_fhir_id, '[^A-Za-z0-9.-]', '_', 'g')
        || '/' || study.id::text
        || CASE
            WHEN position('.' IN study.file_name) > 0 THEN substring(study.file_name FROM '\.[^.]*$')
            ELSE '.dcm'
        END,
    file_name = COALESCE(instance.file_name, study.file_name)
FROM imaging_studies study
WHERE study.id = instance.study_id
  AND instance.object_key IS NULL;

UPDATE imaging_series series
SET instance_count = (
    SELECT COUNT(*)
    FROM imaging_instances instance
    WHERE instance.series_id = series.id
);

UPDATE imaging_studies study
SET series_count = (
        SELECT COUNT(*)
        FROM imaging_series series
        WHERE series.study_id = study.id
    ),
    total_instances = (
        SELECT COUNT(*)
        FROM imaging_instances instance
        WHERE instance.study_id = study.id
    ),
    stored_instances = (
        SELECT COUNT(*)
        FROM imaging_instances instance
        WHERE instance.study_id = study.id
          AND instance.ingest_status = 'STORED'
    );
