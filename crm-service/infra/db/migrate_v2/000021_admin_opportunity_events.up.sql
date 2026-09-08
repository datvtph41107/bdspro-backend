-- Admin Sales opportunity activity log (IV.10.9 nhật ký)
CREATE TABLE IF NOT EXISTS admin_opportunity_events (
    id BIGSERIAL PRIMARY KEY,
    opportunity_id BIGINT NOT NULL,
    actor_id BIGINT NOT NULL,
    action VARCHAR(64) NOT NULL,
    before_json TEXT,
    after_json TEXT,
    note TEXT,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_admin_opportunity_events_opportunity_id
    ON admin_opportunity_events (opportunity_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_admin_opportunity_events_created_at
    ON admin_opportunity_events (created_at DESC)
    WHERE deleted_at IS NULL;
