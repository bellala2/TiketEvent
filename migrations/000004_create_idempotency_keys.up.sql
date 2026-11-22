CREATE TABLE idempotency_keys (
    key             VARCHAR(100) NOT NULL,
    user_id         UUID         NOT NULL REFERENCES users(id),
    request_hash    VARCHAR(64)  NOT NULL,
    response_status INTEGER,
    response_body   JSONB,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, key)
);

-- Untuk membersihkan data yang lebih dari 24 jam.
CREATE INDEX idx_idempotency_keys_created_at ON idempotency_keys (created_at);     