ALTER TABLE imaging_studies
ADD COLUMN file_name VARCHAR(255);

UPDATE imaging_studies
SET file_name = split_part(gcs_path, '/', -1)
WHERE file_name IS NULL;

ALTER TABLE imaging_studies
ALTER COLUMN file_name SET NOT NULL;

ALTER TABLE imaging_studies
DROP COLUMN gcs_path;

ALTER TABLE exam_analyses
DROP COLUMN file_path;
