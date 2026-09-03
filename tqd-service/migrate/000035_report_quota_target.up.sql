BEGIN;

-- Database cũ có bảng này từ GORM AutoMigrate, nhưng canonical SQL chưa từng
-- tạo nó trước khi quota/report job sử dụng. Khởi tạo đầy đủ layout legacy tại
-- dependency đầu tiên; version 000040 sẽ đổi client_request_id thành
-- command_key theo model hiện tại.
CREATE TABLE IF NOT EXISTS user_reported (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    report_type BIGINT NOT NULL,
    status BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    subtitle VARCHAR(255) DEFAULT '',
    address TEXT DEFAULT '',
    province VARCHAR(100) DEFAULT '',
    province_code VARCHAR(20) DEFAULT '',
    ward_code VARCHAR(20) DEFAULT '',
    parcel_id BIGINT,
    region_id BIGINT,
    min_lon DOUBLE PRECISION DEFAULT 0,
    min_lat DOUBLE PRECISION DEFAULT 0,
    max_lon DOUBLE PRECISION DEFAULT 0,
    max_lat DOUBLE PRECISION DEFAULT 0,
    center_lat DOUBLE PRECISION DEFAULT 0,
    center_lon DOUBLE PRECISION DEFAULT 0,
    thumbnail_url TEXT DEFAULT '',
    image_url TEXT DEFAULT '',
    pdf_url TEXT DEFAULT '',
    share_url TEXT DEFAULT '',
    file_size BIGINT DEFAULT 0,
    format VARCHAR(20) DEFAULT '',
    comparison JSONB DEFAULT '{}'::jsonb,
    metadata JSONB DEFAULT '{}'::jsonb,
    client_request_id VARCHAR(100) NOT NULL DEFAULT '',
    job_id VARCHAR(100) DEFAULT '',
    error_message TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

ALTER TABLE user_reported
    ALTER COLUMN client_request_id TYPE VARCHAR(128);

ALTER TABLE user_reported
    ADD COLUMN IF NOT EXISTS request_hash VARCHAR(64) NOT NULL DEFAULT '';

DO $$
DECLARE
    existing_definition TEXT;
BEGIN
    SELECT indexdef
      INTO existing_definition
      FROM pg_indexes
     WHERE schemaname = current_schema()
       AND tablename = 'user_reported'
       AND indexname = 'uidx_user_report_client_request';

    IF existing_definition IS NOT NULL
       AND (existing_definition NOT ILIKE '%UNIQUE INDEX%'
            OR existing_definition NOT ILIKE '%(user_id, client_request_id)%'
            OR existing_definition NOT ILIKE '%deleted_at IS NULL%') THEN
        RAISE EXCEPTION
            'uidx_user_report_client_request exists with an incompatible definition: %',
            existing_definition;
    END IF;
END $$;

CREATE UNIQUE INDEX IF NOT EXISTS uidx_user_report_client_request
    ON user_reported(user_id, client_request_id)
    WHERE client_request_id <> '' AND deleted_at IS NULL;

CREATE TABLE quota_usage_events (
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
    CONSTRAINT chk_quota_usage_period CHECK (period_end > period_start)
);

CREATE INDEX idx_quota_usage_subject_period
    ON quota_usage_events(
        subject_type,
        subject_id,
        operation,
        period_start,
        period_end
    );

CREATE INDEX idx_quota_usage_period_end
    ON quota_usage_events(period_end);

CREATE TABLE report_jobs (
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

CREATE INDEX idx_report_jobs_pending
    ON report_jobs(status, available_at, created_at);
CREATE INDEX idx_report_jobs_operation_id
    ON report_jobs(operation_id);
CREATE INDEX idx_report_jobs_command_key
    ON report_jobs(command_key);

COMMIT;
