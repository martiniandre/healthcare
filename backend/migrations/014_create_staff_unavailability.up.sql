CREATE TABLE staff_unavailability (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    staff_id    UUID NOT NULL REFERENCES employees(id),
    starts_at   TIMESTAMPTZ NOT NULL,
    ends_at     TIMESTAMPTZ NOT NULL,
    reason      VARCHAR(500),
    created_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_staff_unavailability_staff_starts ON staff_unavailability(staff_id, starts_at);
