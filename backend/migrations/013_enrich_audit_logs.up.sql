ALTER TABLE audit_logs
    ADD COLUMN resource_type VARCHAR(100),
    ADD COLUMN resource_id   VARCHAR(255),
    ADD COLUMN action        VARCHAR(100),
    ADD COLUMN payload_diff  JSONB;

CREATE INDEX idx_audit_logs_resource_type ON audit_logs(resource_type);
CREATE INDEX idx_audit_logs_resource_id ON audit_logs(resource_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
