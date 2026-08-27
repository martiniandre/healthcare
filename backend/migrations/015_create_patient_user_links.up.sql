CREATE TABLE patient_user_links (
    user_id         UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    patient_fhir_id VARCHAR(64) NOT NULL UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_patient_user_links_fhir_id ON patient_user_links(patient_fhir_id);

INSERT INTO patient_user_links (user_id, patient_fhir_id)
SELECT id, id::text FROM users WHERE role = 'PATIENT'
ON CONFLICT (user_id) DO NOTHING;
