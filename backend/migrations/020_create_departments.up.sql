CREATE TABLE departments (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL UNIQUE,
    sort_order  INTEGER     NOT NULL DEFAULT 0,
    is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_departments_sort_order ON departments(sort_order);

INSERT INTO departments (name, sort_order) VALUES
    ('Clínica Geral', 1),
    ('Cardiologia', 2),
    ('Dermatologia', 3),
    ('Pediatria', 4),
    ('Neurologia', 5),
    ('Ortopedia e Traumatologia', 6),
    ('Gastroenterologia', 7),
    ('Pneumologia', 8),
    ('Urologia', 9),
    ('Oftalmologia', 10),
    ('Otorrinolaringologia', 11),
    ('Ginecologia e Obstetrícia', 12),
    ('Mastologia', 13),
    ('Endocrinologia e Metabologia', 14),
    ('Hematologia', 15),
    ('Reumatologia', 16),
    ('Infectologia', 17),
    ('Nefrologia', 18),
    ('Oncologia', 19),
    ('Radiologia e Diagnóstico por Imagem', 20);

ALTER TABLE employees ADD COLUMN department_id UUID REFERENCES departments(id);

UPDATE employees
SET department_id = (SELECT id FROM departments WHERE name = 'Clínica Geral')
WHERE department_id IS NULL;

ALTER TABLE employees ALTER COLUMN department_id SET NOT NULL;

CREATE INDEX idx_employees_department_id ON employees(department_id);
