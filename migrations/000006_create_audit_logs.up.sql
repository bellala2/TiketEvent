CREATE TABLE audit_logs (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id    UUID        REFERENCES users(id),  -- NULL jika dilakukan sistem
    actor_type  VARCHAR(20) NOT NULL CHECK (actor_type IN ('user', 'system')),
    entity_type VARCHAR(30) NOT NULL CHECK (entity_type IN ('event', 'order')),
    entity_id   UUID        NOT NULL,
    action      VARCHAR(30) NOT NULL,
    before      JSONB,
    after       JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_entity ON audit_logs (entity_type, entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs (created_at DESC);

-- Audit log tidak boleh diubah atau dihapus, dijaga langsung oleh database.
CREATE OR REPLACE FUNCTION prevent_audit_log_change()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'audit_logs bersifat append-only';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_audit_logs_immutable
BEFORE UPDATE OR DELETE ON audit_logs
FOR EACH ROW EXECUTE FUNCTION prevent_audit_log_change();
