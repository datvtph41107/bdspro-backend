-- Request/audit/search stats tables used by hub request-stats APIs.
-- Idempotent for environments that already created these via legacy GORM AutoMigrate.

CREATE TABLE IF NOT EXISTS click_urls (
    id          BIGSERIAL PRIMARY KEY,
    event_time  BIGINT      NOT NULL,
    domain_id   BIGINT      NOT NULL,
    domain_type INTEGER     NOT NULL,
    profile_id  BIGINT
);
CREATE INDEX IF NOT EXISTS idx_click_urls_event_time ON click_urls (event_time);
CREATE INDEX IF NOT EXISTS idx_click_urls_domain_id ON click_urls (domain_id);
CREATE INDEX IF NOT EXISTS idx_click_urls_domain_type ON click_urls (domain_type);
CREATE INDEX IF NOT EXISTS idx_click_urls_profile_id ON click_urls (profile_id);
CREATE INDEX IF NOT EXISTS idx_click_urls_domain_type_id ON click_urls (domain_type, domain_id);

CREATE TABLE IF NOT EXISTS request_stats_log (
    id          BIGSERIAL PRIMARY KEY,
    method      VARCHAR(16)  NOT NULL,
    path        VARCHAR(512) NOT NULL,
    service     VARCHAR(64)  NOT NULL,
    status_code INTEGER      NOT NULL,
    count       BIGINT       NOT NULL,
    timestamp   TIMESTAMPTZ  NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_request_stats_log_method ON request_stats_log (method);
CREATE INDEX IF NOT EXISTS idx_request_stats_log_path ON request_stats_log (path);
CREATE INDEX IF NOT EXISTS idx_request_stats_log_service ON request_stats_log (service);
CREATE INDEX IF NOT EXISTS idx_request_stats_log_status_code ON request_stats_log (status_code);
CREATE INDEX IF NOT EXISTS idx_request_stats_log_timestamp ON request_stats_log (timestamp);

CREATE TABLE IF NOT EXISTS user_interaction_log (
    id              BIGSERIAL PRIMARY KEY,
    profile_id      BIGINT      NOT NULL,
    organization_id BIGINT,
    count           BIGINT      NOT NULL,
    timestamp       TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_user_interaction_log_profile_id ON user_interaction_log (profile_id);
CREATE INDEX IF NOT EXISTS idx_user_interaction_log_organization_id ON user_interaction_log (organization_id);
CREATE INDEX IF NOT EXISTS idx_user_interaction_log_timestamp ON user_interaction_log (timestamp);

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

CREATE TABLE IF NOT EXISTS keyword_search_log (
    id        BIGSERIAL PRIMARY KEY,
    keyword   VARCHAR(255) NOT NULL,
    source    VARCHAR(64)  NOT NULL,
    count     BIGINT       NOT NULL,
    timestamp TIMESTAMPTZ  NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_keyword_search_log_keyword ON keyword_search_log (keyword);
CREATE INDEX IF NOT EXISTS idx_keyword_search_log_source ON keyword_search_log (source);
CREATE INDEX IF NOT EXISTS idx_keyword_search_log_timestamp ON keyword_search_log (timestamp);

CREATE TABLE IF NOT EXISTS search_latency_log (
    id              BIGSERIAL PRIMARY KEY,
    method          VARCHAR(16)  NOT NULL,
    path            VARCHAR(512) NOT NULL,
    service         VARCHAR(64)  NOT NULL,
    count           BIGINT       NOT NULL,
    sum_duration_ms BIGINT       NOT NULL,
    min_duration_ms BIGINT       NOT NULL,
    max_duration_ms BIGINT       NOT NULL,
    p50_duration_ms BIGINT       NOT NULL DEFAULT 0,
    p95_duration_ms BIGINT       NOT NULL DEFAULT 0,
    p99_duration_ms BIGINT       NOT NULL DEFAULT 0,
    timestamp       TIMESTAMPTZ  NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_search_latency_log_method ON search_latency_log (method);
CREATE INDEX IF NOT EXISTS idx_search_latency_log_path ON search_latency_log (path);
CREATE INDEX IF NOT EXISTS idx_search_latency_log_service ON search_latency_log (service);
CREATE INDEX IF NOT EXISTS idx_search_latency_log_timestamp ON search_latency_log (timestamp);
