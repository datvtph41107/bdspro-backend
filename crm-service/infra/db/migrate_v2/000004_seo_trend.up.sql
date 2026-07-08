-- ============================================================
-- 000004_seo_trend.up.sql
-- Purpose:
--   Store SEO performance metrics over time.
--   This table is for operations and strategy: impressions, clicks,
--   CTR, average position, and source of measurement.
-- ============================================================

CREATE TABLE IF NOT EXISTS seo_trend (
    id BIGSERIAL PRIMARY KEY,

    seo_domain_id BIGINT NOT NULL,
    date DATE NOT NULL,

    impressions BIGINT NOT NULL DEFAULT 0,
    clicks BIGINT NOT NULL DEFAULT 0,
    ctr DOUBLE PRECISION NOT NULL DEFAULT 0,
    avg_position DOUBLE PRECISION NOT NULL DEFAULT 0,

    -- manual, gsc, analytics, import, etc.
    source VARCHAR(50) NOT NULL DEFAULT 'manual',

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

ALTER TABLE seo_trend
    ADD COLUMN IF NOT EXISTS seo_domain_id BIGINT,
    ADD COLUMN IF NOT EXISTS date DATE,
    ADD COLUMN IF NOT EXISTS impressions BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS clicks BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ctr DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS avg_position DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS source VARCHAR(50) NOT NULL DEFAULT 'manual',
    ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

UPDATE seo_trend SET source = 'manual' WHERE source IS NULL OR trim(source) = '';
UPDATE seo_trend SET metadata = '{}'::jsonb WHERE metadata IS NULL;

ALTER TABLE seo_trend
    ALTER COLUMN seo_domain_id SET NOT NULL,
    ALTER COLUMN date SET NOT NULL,
    ALTER COLUMN impressions SET DEFAULT 0,
    ALTER COLUMN impressions SET NOT NULL,
    ALTER COLUMN clicks SET DEFAULT 0,
    ALTER COLUMN clicks SET NOT NULL,
    ALTER COLUMN ctr SET DEFAULT 0,
    ALTER COLUMN ctr SET NOT NULL,
    ALTER COLUMN avg_position SET DEFAULT 0,
    ALTER COLUMN avg_position SET NOT NULL,
    ALTER COLUMN source SET DEFAULT 'manual',
    ALTER COLUMN source SET NOT NULL,
    ALTER COLUMN metadata SET DEFAULT '{}'::jsonb,
    ALTER COLUMN metadata SET NOT NULL,
    ALTER COLUMN created_at SET DEFAULT NOW(),
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN updated_at SET DEFAULT NOW(),
    ALTER COLUMN updated_at SET NOT NULL;

ALTER TABLE seo_trend
    ADD CONSTRAINT chk_seo_trend_non_negative
    CHECK (impressions >= 0 AND clicks >= 0 AND ctr >= 0 AND avg_position >= 0) NOT VALID;

ALTER TABLE seo_trend VALIDATE CONSTRAINT chk_seo_trend_non_negative;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_seo_trend_domain'
    ) THEN
        ALTER TABLE seo_trend
        ADD CONSTRAINT fk_seo_trend_domain
        FOREIGN KEY (seo_domain_id) REFERENCES seo_domain(id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_seo_trend_domain_active
ON seo_trend(seo_domain_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_trend_date_active
ON seo_trend(date DESC)
WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_seo_trend_domain_date_source_active
ON seo_trend(seo_domain_id, date, source)
WHERE deleted_at IS NULL;