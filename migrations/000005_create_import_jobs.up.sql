CREATE TABLE import_jobs (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    file_name    VARCHAR(255) NOT NULL,
    status       VARCHAR(20)  NOT NULL DEFAULT 'queued'
                 CHECK (status IN ('queued', 'processing', 'completed', 'failed')),
    total_rows   INTEGER      NOT NULL DEFAULT 0,
    success_rows INTEGER      NOT NULL DEFAULT 0,
    failed_rows  INTEGER      NOT NULL DEFAULT 0,
    errors       JSONB        NOT NULL DEFAULT '[]'::jsonb,
    created_by   UUID         NOT NULL REFERENCES users(id),
    started_at   TIMESTAMPTZ,
    finished_at  TIMESTAMPTZ,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Dipakai saat aplikasi start: tandai job yang menggantung sebagai failed.
CREATE INDEX idx_import_jobs_status ON import_jobs (status);
