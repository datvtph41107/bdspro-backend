-- Hub remaining tables migrated from hub-service.
-- Idempotent: CREATE TABLE IF NOT EXISTS + CREATE INDEX IF NOT EXISTS.
-- Tables already created in 000030_hub_core.up.sql are NOT repeated here.

CREATE TABLE IF NOT EXISTS applink (
    id BIGSERIAL PRIMARY KEY,
    code BIGINT NOT NULL,
    ref_id BIGINT NOT NULL,
    action INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_applink_code ON applink (code);
CREATE INDEX IF NOT EXISTS idx_applink_ref_action ON applink (ref_id, action);

CREATE TABLE IF NOT EXISTS tb_event_queue (
    id BIGSERIAL PRIMARY KEY,
    event_type TEXT NOT NULL,
    event_name TEXT NOT NULL,
    event_data TEXT,
    metadata TEXT,
    status INTEGER NOT NULL DEFAULT 10,
    profile_id BIGINT,
    organization_id BIGINT,
    scheduled_at TIMESTAMPTZ,
    processed_at TIMESTAMPTZ,
    retry_count BIGINT NOT NULL DEFAULT 0,
    max_retries BIGINT NOT NULL DEFAULT 3,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT
);
CREATE INDEX IF NOT EXISTS idx_tb_event_queue_status ON tb_event_queue (status);
CREATE INDEX IF NOT EXISTS idx_tb_event_queue_profile_id ON tb_event_queue (profile_id);
CREATE INDEX IF NOT EXISTS idx_tb_event_queue_organization_id ON tb_event_queue (organization_id);
CREATE INDEX IF NOT EXISTS idx_tb_event_queue_deleted_at ON tb_event_queue (deleted_at);

CREATE TABLE IF NOT EXISTS system_config (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    key VARCHAR(100) NOT NULL,
    value TEXT,
    group_key INTEGER NOT NULL DEFAULT 10,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT
);
CREATE INDEX IF NOT EXISTS idx_system_config_key ON system_config (key);
CREATE INDEX IF NOT EXISTS idx_system_config_deleted_at ON system_config (deleted_at);

CREATE TABLE IF NOT EXISTS interactive_events (
    start_time BIGINT NOT NULL,
    end_time BIGINT,
    event VARCHAR(30) NOT NULL,
    ref_id BIGINT,
    screen VARCHAR(40),
    duration BIGINT,
    profile_id BIGINT
);
CREATE INDEX IF NOT EXISTS idx_interactive_events_start_time ON interactive_events (start_time);
CREATE INDEX IF NOT EXISTS idx_interactive_events_end_time ON interactive_events (end_time);
CREATE INDEX IF NOT EXISTS idx_interactive_events_event ON interactive_events (event);
CREATE INDEX IF NOT EXISTS idx_interactive_events_ref_id ON interactive_events (ref_id);
CREATE INDEX IF NOT EXISTS idx_interactive_events_profile_id ON interactive_events (profile_id);

CREATE TABLE IF NOT EXISTS error_logs (
    id BIGSERIAL PRIMARY KEY,
    message TEXT NOT NULL,
    stack TEXT,
    screen VARCHAR(200),
    user_id VARCHAR(100),
    device JSONB,
    extra JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_error_logs_user_id ON error_logs (user_id);
CREATE INDEX IF NOT EXISTS idx_error_logs_created_at ON error_logs (created_at);

CREATE TABLE IF NOT EXISTS update_data (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL,
    resource_type INTEGER NOT NULL,
    resource_id BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_update_data_owner_resource ON update_data (owner_id, resource_type);
CREATE INDEX IF NOT EXISTS idx_update_data_updated_at ON update_data (updated_at);

CREATE TABLE IF NOT EXISTS click_urls (
    id BIGSERIAL PRIMARY KEY,
    event_time BIGINT NOT NULL,
    domain_id BIGINT NOT NULL,
    domain_type INTEGER NOT NULL,
    profile_id BIGINT
);
CREATE INDEX IF NOT EXISTS idx_click_urls_event_time ON click_urls (event_time);
CREATE INDEX IF NOT EXISTS idx_click_urls_domain_id ON click_urls (domain_id);
CREATE INDEX IF NOT EXISTS idx_click_urls_domain_type ON click_urls (domain_type);
CREATE INDEX IF NOT EXISTS idx_click_urls_profile_id ON click_urls (profile_id);

CREATE TABLE IF NOT EXISTS keyword_search_log (
    id BIGSERIAL PRIMARY KEY,
    keyword VARCHAR(255) NOT NULL,
    source VARCHAR(64) NOT NULL,
    count BIGINT NOT NULL,
    "timestamp" TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_keyword_search_log_keyword ON keyword_search_log (keyword);
CREATE INDEX IF NOT EXISTS idx_keyword_search_log_source ON keyword_search_log (source);
CREATE INDEX IF NOT EXISTS idx_keyword_search_log_timestamp ON keyword_search_log ("timestamp");

CREATE TABLE IF NOT EXISTS search_latency_log (
    id BIGSERIAL PRIMARY KEY,
    method VARCHAR(16) NOT NULL,
    path VARCHAR(512) NOT NULL,
    service VARCHAR(64) NOT NULL,
    count BIGINT NOT NULL,
    sum_duration_ms BIGINT NOT NULL,
    min_duration_ms BIGINT NOT NULL,
    max_duration_ms BIGINT NOT NULL,
    p50_duration_ms BIGINT NOT NULL DEFAULT 0,
    p95_duration_ms BIGINT NOT NULL DEFAULT 0,
    p99_duration_ms BIGINT NOT NULL DEFAULT 0,
    "timestamp" TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_search_latency_log_method ON search_latency_log (method);
CREATE INDEX IF NOT EXISTS idx_search_latency_log_path ON search_latency_log (path);
CREATE INDEX IF NOT EXISTS idx_search_latency_log_service ON search_latency_log (service);
CREATE INDEX IF NOT EXISTS idx_search_latency_log_timestamp ON search_latency_log ("timestamp");
