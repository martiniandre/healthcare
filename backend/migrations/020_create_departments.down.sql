DROP INDEX IF EXISTS idx_employees_department_id;

ALTER TABLE employees DROP COLUMN IF EXISTS department_id;

DROP TABLE IF EXISTS departments;
