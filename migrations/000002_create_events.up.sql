CREATE TABLE events (
    id          UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(150)  NOT NULL CHECK (char_length(name) BETWEEN 3 AND 150),
    description TEXT,
    venue       VARCHAR(150)  NOT NULL,
    start_at    TIMESTAMPTZ   NOT NULL,
    end_at      TIMESTAMPTZ   NOT NULL,
    price       NUMERIC(15,2) NOT NULL CHECK (price >= 0),
    quota       INTEGER       NOT NULL CHECK (quota BETWEEN 1 AND 10000),
    available   INTEGER       NOT NULL,
    status      VARCHAR(20)   NOT NULL DEFAULT 'draft'
                CHECK (status IN ('draft', 'published', 'cancelled')),
    version     INTEGER       NOT NULL DEFAULT 1,
    created_by  UUID          NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_events_end_after_start CHECK (end_at > start_at),
    -- Pengaman terakhir: stok tidak boleh minus dan tidak boleh melebihi kuota.
    CONSTRAINT chk_events_available CHECK (available >= 0 AND available <= quota)
);

CREATE INDEX idx_events_status_start_at ON events (status, start_at);

CREATE TRIGGER trg_events_updated_at
BEFORE UPDATE ON events
FOR EACH ROW EXECUTE FUNCTION set_updated_at();
