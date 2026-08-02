ALTER TABLE audit_logs
    DROP COLUMN IF EXISTS resource_type,
    DROP COLUMN IF EXISTS resource_id,
    DROP COLUMN IF EXISTS action,
    DROP COLUMN IF EXISTS payload_diff;

DROP INDEX IF EXISTS idx_audit_logs_resource_type;
DROP INDEX IF EXISTS idx_audit_logs_resource_id;
DROP INDEX IF EXISTS idx_audit_logs_action;
