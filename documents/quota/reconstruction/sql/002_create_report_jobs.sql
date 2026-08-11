-- TARGET SOURCE ONLY.
-- Chưa chạy migration ở giai đoạn đọc/hoàn thiện source.
-- PostgreSQL được dùng như durable job queue để không cần thêm Kafka/RabbitMQ cho report hiện tại.

CREATE TABLE IF NOT EXISTS report_jobs (
    id VARCHAR(80) PRIMARY KEY CHECK (id <> ''),
    report_id BIGINT NOT NULL UNIQUE CHECK (report_id > 0) REFERENCES user_reported(id),
    user_id BIGINT NOT NULL CHECK (user_id > 0),

    operation VARCHAR(128) NOT NULL CHECK (operation <> ''),
    operation_id VARCHAR(128) NOT NULL CHECK (operation_id <> ''),
    command_key VARCHAR(128) NOT NULL CHECK (command_key <> ''),

    status VARCHAR(20) NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 0 CHECK (attempts >= 0),

    available_at TIMESTAMPTZ NOT NULL,
    locked_at TIMESTAMPTZ NULL,
    locked_by VARCHAR(100) NOT NULL DEFAULT '',
    last_error TEXT NOT NULL DEFAULT '',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_report_jobs_status
        CHECK (status IN ('pending', 'running', 'completed', 'failed')),
    CONSTRAINT chk_report_jobs_lock_state
        CHECK (
            (status = 'running' AND locked_at IS NOT NULL AND locked_by <> '')
            OR
            (status <> 'running' AND locked_at IS NULL AND locked_by = '')
        )
);

CREATE INDEX IF NOT EXISTS idx_report_jobs_pending
    ON report_jobs(status, available_at, created_at);

CREATE INDEX IF NOT EXISTS idx_report_jobs_operation_id
    ON report_jobs(operation_id);

CREATE INDEX IF NOT EXISTS idx_report_jobs_command_key
    ON report_jobs(command_key);

-- Current user_reported.client_request_id is varchar(100).
-- Target command key follows Operation-ID / Idempotency-Key max length, so widen it first.
ALTER TABLE user_reported
    ALTER COLUMN client_request_id TYPE VARCHAR(128);

ALTER TABLE user_reported
    ADD COLUMN IF NOT EXISTS request_hash VARCHAR(64) NOT NULL DEFAULT '';

CREATE UNIQUE INDEX IF NOT EXISTS uidx_user_report_client_request
    ON user_reported(user_id, client_request_id)
    WHERE client_request_id <> '' AND deleted_at IS NULL;
