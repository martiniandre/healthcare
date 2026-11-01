ALTER TABLE exam_analyses
    ADD COLUMN IF NOT EXISTS anonymization_report JSONB;
