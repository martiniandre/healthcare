ALTER TABLE exam_analyses
ADD COLUMN file_path VARCHAR(1024);

ALTER TABLE imaging_studies
ADD COLUMN gcs_path VARCHAR(1024);

UPDATE imaging_studies
SET gcs_path = file_name;
