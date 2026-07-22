-- SEO-M8 production measurement control plane.
-- Additive only: existing M8 evidence tables remain the metric source of truth.

CREATE TABLE IF NOT EXISTS seo_measurement_import_run (
    id BIGSERIAL PRIMARY KEY,
    source VARCHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'running',
    requested_by VARCHAR(128) NOT NULL DEFAULT '',
    range_start TIMESTAMPTZ NULL,
    range_end TIMESTAMPTZ NULL,
    row_count BIGINT NOT NULL DEFAULT 0,
    accepted_count BIGINT NOT NULL DEFAULT 0,
    rejected_count BIGINT NOT NULL DEFAULT 0,
    checksum VARCHAR(128) NOT NULL DEFAULT '',
    error_message TEXT NOT NULL DEFAULT '',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_seo_measurement_import_run_source ON seo_measurement_import_run(source, started_at DESC);
CREATE INDEX IF NOT EXISTS idx_seo_measurement_import_run_status ON seo_measurement_import_run(status, started_at DESC);
CREATE UNIQUE INDEX IF NOT EXISTS uq_seo_measurement_import_run_active_checksum
    ON seo_measurement_import_run(source, checksum)
    WHERE checksum <> '' AND status IN ('running', 'succeeded', 'partial');

CREATE TABLE IF NOT EXISTS seo_measurement_source_state (
    source VARCHAR(64) PRIMARY KEY,
    status VARCHAR(32) NOT NULL DEFAULT 'never_imported',
    last_run_id BIGINT NULL,
    last_attempt_at TIMESTAMPTZ NULL,
    last_success_at TIMESTAMPTZ NULL,
    last_row_count BIGINT NOT NULL DEFAULT 0,
    stale_after_hours INT NOT NULL DEFAULT 48,
    consecutive_failures INT NOT NULL DEFAULT 0,
    last_error TEXT NOT NULL DEFAULT '',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS seo_page_score_snapshot (
    id BIGSERIAL PRIMARY KEY,
    seo_domain_id BIGINT NULL,
    canonical_url TEXT NOT NULL,
    range_key VARCHAR(32) NOT NULL DEFAULT '28d',
    measured_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    index_score DOUBLE PRECISION NOT NULL DEFAULT 0,
    visibility_score DOUBLE PRECISION NOT NULL DEFAULT 0,
    engagement_score DOUBLE PRECISION NOT NULL DEFAULT 0,
    conversion_score DOUBLE PRECISION NOT NULL DEFAULT 0,
    technical_score DOUBLE PRECISION NOT NULL DEFAULT 0,
    issue_penalty DOUBLE PRECISION NOT NULL DEFAULT 0,
    overall_score DOUBLE PRECISION NOT NULL DEFAULT 0,
    grade VARCHAR(8) NOT NULL DEFAULT 'F',
    blockers JSONB NOT NULL DEFAULT '[]'::jsonb,
    warnings JSONB NOT NULL DEFAULT '[]'::jsonb,
    source_freshness JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(canonical_url, range_key, measured_at)
);
CREATE INDEX IF NOT EXISTS idx_seo_page_score_url ON seo_page_score_snapshot(canonical_url, measured_at DESC);
CREATE INDEX IF NOT EXISTS idx_seo_page_score_grade ON seo_page_score_snapshot(grade, overall_score);

ALTER TABLE seo_search_performance_daily ADD COLUMN IF NOT EXISTS source_run_id BIGINT NULL;
ALTER TABLE seo_url_inspection_snapshot ADD COLUMN IF NOT EXISTS source_run_id BIGINT NULL;
ALTER TABLE seo_traffic_daily ADD COLUMN IF NOT EXISTS source_run_id BIGINT NULL;
ALTER TABLE seo_ranking_snapshot ADD COLUMN IF NOT EXISTS source_run_id BIGINT NULL;
ALTER TABLE seo_crawl_event ADD COLUMN IF NOT EXISTS source_run_id BIGINT NULL;
ALTER TABLE seo_crawl_event ADD COLUMN IF NOT EXISTS data_source VARCHAR(32) NOT NULL DEFAULT 'crawl_log';

CREATE INDEX IF NOT EXISTS idx_seo_spd_source_run ON seo_search_performance_daily(source_run_id);
CREATE INDEX IF NOT EXISTS idx_seo_inspection_source_run ON seo_url_inspection_snapshot(source_run_id);
CREATE INDEX IF NOT EXISTS idx_seo_traffic_source_run ON seo_traffic_daily(source_run_id);
CREATE INDEX IF NOT EXISTS idx_seo_ranking_source_run ON seo_ranking_snapshot(source_run_id);
CREATE INDEX IF NOT EXISTS idx_seo_crawl_source_run ON seo_crawl_event(source_run_id);
