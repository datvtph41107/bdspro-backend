BEGIN;

CREATE TABLE IF NOT EXISTS payment_stats (
    id BIGSERIAL PRIMARY KEY,
    calculate_time TIMESTAMPTZ NOT NULL,
    count BIGINT NOT NULL DEFAULT 0 CHECK (count >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_payment_stats_calculate_time UNIQUE (calculate_time)
);

COMMIT;
