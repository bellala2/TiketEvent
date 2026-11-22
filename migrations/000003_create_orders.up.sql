CREATE TABLE orders (
    id           UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID          NOT NULL REFERENCES users(id),
    event_id     UUID          NOT NULL REFERENCES events(id),
    quantity     INTEGER       NOT NULL CHECK (quantity BETWEEN 1 AND 5),
    unit_price   NUMERIC(15,2) NOT NULL CHECK (unit_price >= 0),
    total_price  NUMERIC(15,2) NOT NULL,
    status       VARCHAR(20)   NOT NULL DEFAULT 'pending'
                 CHECK (status IN ('pending', 'paid', 'cancelled', 'expired')),
    expires_at   TIMESTAMPTZ   NOT NULL,
    paid_at      TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_orders_total CHECK (total_price = quantity * unit_price)
);

-- Dipakai worker expiry: cari order pending yang sudah lewat waktu.
CREATE INDEX idx_orders_status_expires_at ON orders (status, expires_at);
-- Dipakai GET /orders milik customer.
CREATE INDEX idx_orders_user_id ON orders (user_id);
-- Dipakai menghitung total tiket customer per event (batas 5).
CREATE INDEX idx_orders_event_user ON orders (event_id, user_id);

CREATE TRIGGER trg_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW EXECUTE FUNCTION set_updated_at();
