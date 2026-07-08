-- Phase 3 renamed:
-- SEO Measurement & Feedback Loop
--
-- Purpose:
--   Store crawl/index/search/traffic/ranking/issue evidence so admin can see
--   whether SEO pages are working, not just whether they were created.
--
-- Mock-first principle:
--   data_source='mock' rows allow admin/frontend to run now.
--   Later connectors can write data_source='search_console', 'analytics',
--   'crawl_log', 'rank_tracker' into the same schema.

CREATE TABLE IF NOT EXISTS seo_search_performance_daily (
    id BIGSERIAL PRIMARY KEY,
    seo_domain_id BIGINT NULL,
    canonical_url TEXT NOT NULL,
    query TEXT NOT NULL DEFAULT '',
    device VARCHAR(32) NOT NULL DEFAULT '',
    country VARCHAR(16) NOT NULL DEFAULT '',
    search_type VARCHAR(32) NOT NULL DEFAULT 'web',
    date DATE NOT NULL,
    clicks DOUBLE PRECISION NOT NULL DEFAULT 0,
    impressions DOUBLE PRECISION NOT NULL DEFAULT 0,
    ctr DOUBLE PRECISION NOT NULL DEFAULT 0,
    position DOUBLE PRECISION NOT NULL DEFAULT 0,
    data_source VARCHAR(32) NOT NULL DEFAULT 'mock',
    imported_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(canonical_url, query, device, country, search_type, date)
);
CREATE INDEX IF NOT EXISTS idx_seo_spd_date ON seo_search_performance_daily(date);
CREATE INDEX IF NOT EXISTS idx_seo_spd_url ON seo_search_performance_daily(canonical_url);
CREATE INDEX IF NOT EXISTS idx_seo_spd_query ON seo_search_performance_daily(query);

CREATE TABLE IF NOT EXISTS seo_url_inspection_snapshot (
    id BIGSERIAL PRIMARY KEY,
    seo_domain_id BIGINT NULL,
    canonical_url TEXT NOT NULL,
    inspection_status VARCHAR(64) NOT NULL DEFAULT '',
    index_status VARCHAR(64) NOT NULL DEFAULT '',
    coverage_state TEXT NOT NULL DEFAULT '',
    robots_txt_state VARCHAR(64) NOT NULL DEFAULT '',
    indexing_state VARCHAR(64) NOT NULL DEFAULT '',
    google_canonical TEXT NOT NULL DEFAULT '',
    user_canonical TEXT NOT NULL DEFAULT '',
    last_crawl_time TIMESTAMPTZ NULL,
    page_fetch_state VARCHAR(64) NOT NULL DEFAULT '',
    verdict VARCHAR(64) NOT NULL DEFAULT '',
    raw_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    data_source VARCHAR(32) NOT NULL DEFAULT 'mock',
    inspected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(canonical_url, data_source)
);
CREATE INDEX IF NOT EXISTS idx_seo_inspection_url ON seo_url_inspection_snapshot(canonical_url);
CREATE INDEX IF NOT EXISTS idx_seo_inspection_index_status ON seo_url_inspection_snapshot(index_status);
CREATE INDEX IF NOT EXISTS idx_seo_inspection_inspected_at ON seo_url_inspection_snapshot(inspected_at);

CREATE TABLE IF NOT EXISTS seo_traffic_daily (
    id BIGSERIAL PRIMARY KEY,
    seo_domain_id BIGINT NULL,
    canonical_url TEXT NOT NULL,
    date DATE NOT NULL,
    sessions DOUBLE PRECISION NOT NULL DEFAULT 0,
    users DOUBLE PRECISION NOT NULL DEFAULT 0,
    active_users DOUBLE PRECISION NOT NULL DEFAULT 0,
    page_views DOUBLE PRECISION NOT NULL DEFAULT 0,
    engagement_seconds DOUBLE PRECISION NOT NULL DEFAULT 0,
    conversions DOUBLE PRECISION NOT NULL DEFAULT 0,
    source VARCHAR(128) NOT NULL DEFAULT '',
    medium VARCHAR(128) NOT NULL DEFAULT '',
    device VARCHAR(32) NOT NULL DEFAULT '',
    country VARCHAR(16) NOT NULL DEFAULT '',
    data_source VARCHAR(32) NOT NULL DEFAULT 'mock',
    imported_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(canonical_url, date, source, medium, device, country)
);
CREATE INDEX IF NOT EXISTS idx_seo_traffic_date ON seo_traffic_daily(date);
CREATE INDEX IF NOT EXISTS idx_seo_traffic_url ON seo_traffic_daily(canonical_url);

CREATE TABLE IF NOT EXISTS seo_keyword_target (
    id BIGSERIAL PRIMARY KEY,
    seo_domain_id BIGINT NOT NULL,
    keyword TEXT NOT NULL,
    intent VARCHAR(64) NOT NULL DEFAULT '',
    priority INT NOT NULL DEFAULT 0,
    target_url TEXT NOT NULL DEFAULT '',
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(seo_domain_id, keyword)
);
CREATE INDEX IF NOT EXISTS idx_seo_keyword_target_keyword ON seo_keyword_target(keyword);
CREATE INDEX IF NOT EXISTS idx_seo_keyword_target_url ON seo_keyword_target(target_url);

CREATE TABLE IF NOT EXISTS seo_ranking_snapshot (
    id BIGSERIAL PRIMARY KEY,
    keyword_target_id BIGINT NULL,
    seo_domain_id BIGINT NULL,
    keyword TEXT NOT NULL,
    canonical_url TEXT NOT NULL,
    position DOUBLE PRECISION NOT NULL DEFAULT 0,
    previous_position DOUBLE PRECISION NOT NULL DEFAULT 0,
    position_delta DOUBLE PRECISION NOT NULL DEFAULT 0,
    device VARCHAR(32) NOT NULL DEFAULT '',
    location VARCHAR(128) NOT NULL DEFAULT '',
    language VARCHAR(32) NOT NULL DEFAULT '',
    data_source VARCHAR(32) NOT NULL DEFAULT 'mock',
    checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    raw_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_seo_ranking_keyword ON seo_ranking_snapshot(keyword);
CREATE INDEX IF NOT EXISTS idx_seo_ranking_url ON seo_ranking_snapshot(canonical_url);
CREATE INDEX IF NOT EXISTS idx_seo_ranking_checked_at ON seo_ranking_snapshot(checked_at);

CREATE TABLE IF NOT EXISTS seo_crawl_event (
    id BIGSERIAL PRIMARY KEY,
    seo_domain_id BIGINT NULL,
    canonical_url TEXT NOT NULL,
    bot_name VARCHAR(64) NOT NULL DEFAULT '',
    user_agent TEXT NOT NULL DEFAULT '',
    status_code INT NOT NULL DEFAULT 0,
    method VARCHAR(16) NOT NULL DEFAULT 'GET',
    response_time_ms BIGINT NOT NULL DEFAULT 0,
    referer TEXT NOT NULL DEFAULT '',
    ip TEXT NOT NULL DEFAULT '',
    crawled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    raw_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_seo_crawl_url ON seo_crawl_event(canonical_url);
CREATE INDEX IF NOT EXISTS idx_seo_crawl_bot ON seo_crawl_event(bot_name);
CREATE INDEX IF NOT EXISTS idx_seo_crawl_status ON seo_crawl_event(status_code);
CREATE INDEX IF NOT EXISTS idx_seo_crawl_time ON seo_crawl_event(crawled_at);

CREATE TABLE IF NOT EXISTS seo_issue (
    id BIGSERIAL PRIMARY KEY,
    seo_domain_id BIGINT NULL,
    canonical_url TEXT NOT NULL,
    issue_type VARCHAR(64) NOT NULL,
    severity VARCHAR(32) NOT NULL DEFAULT 'medium',
    status VARCHAR(32) NOT NULL DEFAULT 'open',
    title TEXT NOT NULL DEFAULT '',
    message TEXT NOT NULL DEFAULT '',
    evidence JSONB NOT NULL DEFAULT '{}'::jsonb,
    recommendation TEXT NOT NULL DEFAULT '',
    detected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(canonical_url, issue_type, status)
);
CREATE INDEX IF NOT EXISTS idx_seo_issue_url ON seo_issue(canonical_url);
CREATE INDEX IF NOT EXISTS idx_seo_issue_type ON seo_issue(issue_type);
CREATE INDEX IF NOT EXISTS idx_seo_issue_status ON seo_issue(status);
CREATE INDEX IF NOT EXISTS idx_seo_issue_severity ON seo_issue(severity);
