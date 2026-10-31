CREATE TABLE refresh_tokens (
    id                  UUID PRIMARY KEY,
    family_id           UUID        NOT NULL,
    user_id             UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash          TEXT        NOT NULL UNIQUE,
    issued_at           TIMESTAMPTZ NOT NULL,
    expires_at          TIMESTAMPTZ NOT NULL,
    absolute_expires_at TIMESTAMPTZ NOT NULL,
    replaced_by         UUID,
    revoked_at          TIMESTAMPTZ,
    user_agent          TEXT,
    ip                  TEXT
);

CREATE INDEX idx_refresh_tokens_family_id ON refresh_tokens(family_id);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
