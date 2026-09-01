CREATE INDEX idx_audit_logs_access_granted_created_at
    ON audit_logs (access_granted, created_at DESC);