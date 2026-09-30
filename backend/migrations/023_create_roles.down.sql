ALTER TABLE employees ADD COLUMN role VARCHAR(50);

UPDATE employees AS employee
SET role = roles.code
FROM roles
WHERE employee.role_id = roles.id;

ALTER TABLE employees ALTER COLUMN role SET NOT NULL;

CREATE INDEX idx_employees_role ON employees(role);

ALTER TABLE employees DROP COLUMN role_id;

DROP TABLE roles;