-- TARGET SOURCE ONLY.
-- Chưa chạy migration ở giai đoạn đọc/hoàn thiện source.
-- File này mô tả table mà usageevent Postgres adapter cần sau khi bắt đầu migration thật.

CREATE TABLE IF NOT EXISTS quota_usage_events (
    id BIGSERIAL PRIMARY KEY,
    usage_key CHAR(64) NOT NULL UNIQUE,

    subject_type VARCHAR(32) NOT NULL
        CHECK (subject_type IN ('profile', 'organization')),
    subject_id VARCHAR(128) NOT NULL CHECK (subject_id <> ''),
    operation VARCHAR(128) NOT NULL CHECK (operation <> ''),

    operation_id VARCHAR(128) NOT NULL CHECK (operation_id <> ''),
    idempotency_key VARCHAR(128) NOT NULL DEFAULT '',
    command_key VARCHAR(128) NOT NULL CHECK (command_key <> ''),
    reservation_id VARCHAR(128) NOT NULL UNIQUE CHECK (reservation_id <> ''),

    amount BIGINT NOT NULL CHECK (amount > 0),

    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_quota_usage_period
        CHECK (period_end > period_start)
);

CREATE INDEX IF NOT EXISTS idx_quota_usage_subject_period
    ON quota_usage_events(
        subject_type,
        subject_id,
        operation,
        period_start,
        period_end
    );

CREATE INDEX IF NOT EXISTS idx_quota_usage_period_end
    ON quota_usage_events(period_end);
