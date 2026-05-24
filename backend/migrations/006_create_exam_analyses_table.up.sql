CREATE TABLE IF NOT EXISTS exam_analyses (
    id UUID PRIMARY KEY,
    user_id UUID,
    patient_fhir_id VARCHAR(255),
    exam_type VARCHAR(255),
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(1024) NOT NULL,
    status VARCHAR(50) NOT NULL,
    analysis_response JSONB,
    consent_given BOOLEAN NOT NULL DEFAULT FALSE,
    anonymized BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS exam_analyses_audit_logs (
    id UUID PRIMARY KEY,
    analysis_id UUID,
    action_type VARCHAR(50) NOT NULL,
    performed_by VARCHAR(255) NOT NULL,
    ip_address VARCHAR(50),
    details TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
