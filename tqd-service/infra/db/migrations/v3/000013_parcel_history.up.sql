-- 000013_parcel_history.up.sql
-- Lịch sử thửa đất — cùng shape với qh_planning_events (Đồ án quy hoạch).

CREATE TABLE IF NOT EXISTS parcel_history (
    id            BIGSERIAL PRIMARY KEY,
    parcel_id     BIGINT NOT NULL,
    event_type    VARCHAR(80),
    event_name    VARCHAR(255) NOT NULL,
    description   TEXT,
    event_date    TIMESTAMPTZ,
    document_id   BIGINT NULL,
    source_url    TEXT,
    metadata      JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ NULL,
    created_by    BIGINT,
    updated_by    BIGINT
);

CREATE INDEX IF NOT EXISTS idx_parcel_history_parcel_timeline
    ON parcel_history (parcel_id, event_date DESC, id DESC)
    WHERE deleted_at IS NULL AND parcel_id IS NOT NULL;
