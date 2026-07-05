ALTER TABLE employees DROP CONSTRAINT IF EXISTS employees_user_id_fkey;
ALTER TABLE employees DROP COLUMN IF EXISTS user_id;
ALTER TABLE employees ADD COLUMN created_by UUID;
