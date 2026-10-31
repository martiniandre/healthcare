CREATE TABLE outbox_events (
    id           UUID PRIMARY KEY,
    name         TEXT NOT NULL,
    payload      JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    available_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    attempts     INT NOT NULL DEFAULT 0,
    last_error   TEXT,
    processed_at TIMESTAMPTZ,
    dead_at      TIMESTAMPTZ
);

CREATE INDEX idx_outbox_events_pending ON outbox_events (available_at)
    WHERE processed_at IS NULL AND dead_at IS NULL;
