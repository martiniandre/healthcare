CREATE TABLE diagnostic_report_versions (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id  VARCHAR(255) NOT NULL,
    version    VARCHAR(50)  NOT NULL,
    snapshot   JSONB        NOT NULL,
    changed_by UUID,
    changed_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_diagnostic_report_versions_report_version UNIQUE (report_id, version)
);

CREATE INDEX idx_diagnostic_report_versions_report_id ON diagnostic_report_versions(report_id);
