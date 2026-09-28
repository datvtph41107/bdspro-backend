-- Request stats / audit tables owned by v4 (migrated from hub-service).
-- Final shape includes 5-minute bucket_hour unique keys for upsert aggregation.

CREATE TABLE IF NOT EXISTS request_stats_log (
    id          BIGSERIAL PRIMARY KEY,
    method      VARCHAR(16)  NOT NULL,
    path        VARCHAR(512) NOT NULL,
    service     VARCHAR(64)  NOT NULL,
    status_code INTEGER      NOT NULL,
    count       BIGINT       NOT NULL,
    timestamp   TIMESTAMPTZ  NOT NULL,
    bucket_hour TIMESTAMPTZ  NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_request_stats_log_method ON request_stats_log (method);
CREATE INDEX IF NOT EXISTS idx_request_stats_log_path ON request_stats_log (path);
CREATE INDEX IF NOT EXISTS idx_request_stats_log_service ON request_stats_log (service);
CREATE INDEX IF NOT EXISTS idx_request_stats_log_status_code ON request_stats_log (status_code);
CREATE INDEX IF NOT EXISTS idx_request_stats_log_timestamp ON request_stats_log (timestamp);
CREATE UNIQUE INDEX IF NOT EXISTS uq_request_stats_log_hour_bucket
    ON request_stats_log (method, path, service, status_code, bucket_hour);

CREATE TABLE IF NOT EXISTS user_interaction_log (
    id              BIGSERIAL PRIMARY KEY,
    profile_id      BIGINT      NOT NULL,
    organization_id BIGINT      NOT NULL DEFAULT 0,
    count           BIGINT      NOT NULL,
    timestamp       TIMESTAMPTZ NOT NULL,
    bucket_hour     TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_user_interaction_log_profile_id ON user_interaction_log (profile_id);
CREATE INDEX IF NOT EXISTS idx_user_interaction_log_organization_id ON user_interaction_log (organization_id);
CREATE INDEX IF NOT EXISTS idx_user_interaction_log_timestamp ON user_interaction_log (timestamp);
CREATE UNIQUE INDEX IF NOT EXISTS uq_user_interaction_log_hour_bucket
    ON user_interaction_log (profile_id, organization_id, bucket_hour);

CREATE TABLE IF NOT EXISTS request_audit_log (
    id          BIGSERIAL PRIMARY KEY,
    method      VARCHAR(16)  NOT NULL,
    path        VARCHAR(512) NOT NULL,
    service     VARCHAR(64)  NOT NULL,
    status_code INTEGER      NOT NULL,
    duration_ms BIGINT       NOT NULL,
    timestamp   TIMESTAMPTZ  NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_request_audit_log_method ON request_audit_log (method);
CREATE INDEX IF NOT EXISTS idx_request_audit_log_path ON request_audit_log (path);
CREATE INDEX IF NOT EXISTS idx_request_audit_log_service ON request_audit_log (service);
CREATE INDEX IF NOT EXISTS idx_request_audit_log_status_code ON request_audit_log (status_code);
CREATE INDEX IF NOT EXISTS idx_request_audit_log_timestamp ON request_audit_log (timestamp);
