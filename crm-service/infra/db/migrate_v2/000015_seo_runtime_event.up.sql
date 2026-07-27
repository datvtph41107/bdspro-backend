-- SEO Operating System V2.U runtime interaction evidence.
-- This table stores sanitized first-party events emitted by Next/RN. It is
-- deliberately separate from Search Console, Analytics and crawl connectors.
CREATE TABLE IF NOT EXISTS seo_runtime_event (
    id BIGSERIAL PRIMARY KEY,
    event_key VARCHAR(64) NOT NULL,
    schema_version SMALLINT NOT NULL DEFAULT 1,
    source VARCHAR(32) NOT NULL,
    event_name VARCHAR(80) NOT NULL,
    session_id VARCHAR(120) NOT NULL,
    canonical_url TEXT NOT NULL,
    canonical_path TEXT NOT NULL,
    module_id VARCHAR(100) NOT NULL,
    environment VARCHAR(40) NOT NULL DEFAULT '',
    indexable BOOLEAN NULL,
    action VARCHAR(100) NOT NULL DEFAULT '',
    target TEXT NOT NULL DEFAULT '',
    occurred_at TIMESTAMPTZ NOT NULL,
    received_at TIMESTAMPTZ NOT NULL,
    metrics JSONB NOT NULL DEFAULT '{}'::jsonb,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    user_agent VARCHAR(300) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_seo_runtime_event_event_key
    ON seo_runtime_event(event_key);
CREATE INDEX IF NOT EXISTS idx_seo_runtime_event_canonical_time
    ON seo_runtime_event(canonical_path, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_seo_runtime_event_module_event_time
    ON seo_runtime_event(module_id, event_name, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_seo_runtime_event_session_time
    ON seo_runtime_event(session_id, occurred_at DESC);
