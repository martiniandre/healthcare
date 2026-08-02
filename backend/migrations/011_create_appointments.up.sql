CREATE TABLE appointments (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    patient_fhir_id VARCHAR(255) NOT NULL,
    staff_id        UUID NOT NULL REFERENCES employees(id),
    starts_at       TIMESTAMPTZ NOT NULL,
    ends_at         TIMESTAMPTZ NOT NULL,
    status          VARCHAR(30)  NOT NULL DEFAULT 'scheduled',
    reason          VARCHAR(500),
    version         INTEGER      NOT NULL DEFAULT 1,
    created_by      UUID,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_appointments_staff_starts ON appointments(staff_id, starts_at);
CREATE INDEX idx_appointments_patient ON appointments(patient_fhir_id);

CREATE UNIQUE INDEX uq_appointments_occupied_slot
    ON appointments(staff_id, starts_at)
    WHERE status IN ('scheduled', 'confirmed');

CREATE TABLE idempotency_keys (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    idempotency_key VARCHAR(255) NOT NULL UNIQUE,
    request_hash    VARCHAR(255) NOT NULL,
    response_status INTEGER      NOT NULL,
    response_body   JSONB        NOT NULL,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
