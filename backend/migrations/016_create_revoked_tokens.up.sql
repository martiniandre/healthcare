CREATE TABLE revoked_tokens (
    token_digest VARCHAR(64) PRIMARY KEY,
    expires_at   TIMESTAMPTZ NOT NULL,
    revoked_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_revoked_tokens_expires_at ON revoked_tokens(expires_at);
