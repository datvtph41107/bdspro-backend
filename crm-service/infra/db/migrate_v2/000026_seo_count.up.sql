-- Batch flush log: each row is one in-memory cache window.
CREATE TABLE IF NOT EXISTS seo_count_log (
    seo_id BIGINT NOT NULL,
    clicked BIGINT NOT NULL DEFAULT 0,
    "timestamp" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_seo_count_log_seo_id_timestamp
    ON seo_count_log(seo_id, "timestamp" DESC);

COMMENT ON TABLE seo_count_log IS
    'Batch click flush log; each row is one cache window from the public click API.';

-- Aggregate click counter per SEO page (source of truth for admin list).
CREATE TABLE IF NOT EXISTS seo_count (
    seo_id BIGINT PRIMARY KEY,
    count BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_seo_count_updated_at
    ON seo_count(updated_at DESC);

COMMENT ON TABLE seo_count IS
    'Aggregated SEO page click totals; updated when seo_count_log is flushed.';
