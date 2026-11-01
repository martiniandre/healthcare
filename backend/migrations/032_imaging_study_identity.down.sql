-- Rollback of migration 032.
-- Restores the non-partial unique index from migration 028. This fails if
-- duplicate study_instance_uid values (including several empty strings) exist.
-- Counter values recomputed by 032 are kept; they match the real rows.

DROP INDEX IF EXISTS uq_imaging_studies_study_instance_uid;

CREATE UNIQUE INDEX uq_imaging_studies_study_instance_uid
ON imaging_studies(study_instance_uid);

ALTER TABLE imaging_instances
DROP COLUMN IF EXISTS object_key;

ALTER TABLE imaging_instances
DROP COLUMN IF EXISTS file_name;

ALTER TABLE imaging_studies
DROP COLUMN IF EXISTS fhir_resource_id;
