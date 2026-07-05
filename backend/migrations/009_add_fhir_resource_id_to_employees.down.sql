DROP INDEX IF EXISTS idx_employees_fhir_resource_id;
ALTER TABLE employees DROP COLUMN IF EXISTS fhir_resource_id;
