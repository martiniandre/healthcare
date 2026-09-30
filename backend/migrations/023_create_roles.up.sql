CREATE TABLE roles (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code       VARCHAR(50) NOT NULL UNIQUE,
    name       VARCHAR(100) NOT NULL,
    sort_order INTEGER     NOT NULL DEFAULT 0,
    is_active  BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_roles_sort_order ON roles(sort_order);

INSERT INTO roles (code, name, sort_order) VALUES
    ('ADMIN', 'Administrador', 1),
    ('DOCTOR', 'Médico', 2),
    ('NURSE', 'Enfermeiro', 3),
    ('RECEPTION', 'Recepcionista', 4),
    ('PATIENT', 'Paciente', 5);

ALTER TABLE employees ADD COLUMN role_id UUID REFERENCES roles(id);

UPDATE employees AS employee
SET role_id = roles.id
FROM roles
WHERE employee.role_id IS NULL
  AND roles.code = REPLACE(UPPER(employee.role), 'ROLE', '');

ALTER TABLE employees ALTER COLUMN role_id SET NOT NULL;

CREATE INDEX idx_employees_role_id ON employees(role_id);

ALTER TABLE employees DROP COLUMN role;