ALTER TABLE employees ADD COLUMN fhir_resource_id VARCHAR(255);
CREATE INDEX idx_employees_fhir_resource_id ON employees(fhir_resource_id);
