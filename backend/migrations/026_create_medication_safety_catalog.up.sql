CREATE TABLE medication_classes (
    code         TEXT PRIMARY KEY,
    name         TEXT NOT NULL,
    search_terms TEXT[] NOT NULL DEFAULT '{}'
);

CREATE TABLE medication_catalog (
    code              TEXT PRIMARY KEY,
    name              TEXT NOT NULL,
    active_ingredient TEXT NOT NULL,
    class_code        TEXT REFERENCES medication_classes (code),
    search_terms      TEXT[] NOT NULL DEFAULT '{}',
    catalog_version   TEXT NOT NULL DEFAULT 'sample-1'
);

CREATE INDEX idx_medication_catalog_class_code ON medication_catalog (class_code);
CREATE INDEX idx_medication_catalog_active_ingredient ON medication_catalog (active_ingredient);

CREATE TABLE medication_interactions (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    subject_a_type  TEXT NOT NULL CHECK (subject_a_type IN ('ingredient', 'class')),
    subject_a       TEXT NOT NULL,
    subject_b_type  TEXT NOT NULL CHECK (subject_b_type IN ('ingredient', 'class')),
    subject_b       TEXT NOT NULL,
    severity        TEXT NOT NULL CHECK (severity IN ('moderate', 'major', 'contraindicated')),
    description_key TEXT NOT NULL,
    catalog_version TEXT NOT NULL DEFAULT 'sample-1',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_medication_interactions_unique_pair ON medication_interactions (
    LEAST(subject_a_type || ':' || subject_a, subject_b_type || ':' || subject_b),
    GREATEST(subject_a_type || ':' || subject_a, subject_b_type || ':' || subject_b)
);

CREATE TABLE medication_safety_overrides (
    id                           UUID PRIMARY KEY,
    patient_fhir_id              TEXT NOT NULL,
    medication_request_fhir_id   TEXT,
    medication_code              TEXT NOT NULL DEFAULT '',
    practitioner_user_id         UUID NOT NULL REFERENCES users (id),
    alert_kinds                  TEXT[] NOT NULL,
    justification                TEXT NOT NULL CHECK (char_length(btrim(justification)) >= 10),
    created_at                   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_medication_safety_overrides_patient ON medication_safety_overrides (patient_fhir_id, created_at DESC);
CREATE INDEX idx_medication_safety_overrides_request ON medication_safety_overrides (medication_request_fhir_id);
