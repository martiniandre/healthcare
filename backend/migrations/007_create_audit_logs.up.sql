CREATE TABLE audit_logs (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    correlation_id VARCHAR(255) NOT NULL,
    caller_user_id VARCHAR(255) NOT NULL,
    caller_role    VARCHAR(50) NOT NULL,
    method         VARCHAR(255) NOT NULL,
    access_granted BOOLEAN NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_correlation_id ON audit_logs(correlation_id);
CREATE INDEX idx_audit_logs_caller_user_id ON audit_logs(caller_user_id);
CREATE INDEX idx_audit_logs_method ON audit_logs(method);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
