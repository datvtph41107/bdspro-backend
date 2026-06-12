-- ============================================================
-- 000001_seo_domain_registry.up.sql
-- Purpose:
--   Build seo_domain as the central SEO content registry.
--   This table is the source of truth for URL identity, publish state,
--   sitemap policy, static HTML generation state, source binding,
--   deeplink, metadata/schema, and long-term SEO operations.
-- Notes:
--   - Designed to be safe for projects that already have seo_domain.
--   - If duplicate canonical_url values exist among active rows, the
--     unique index at the end will fail. Clean duplicates before running.
-- ============================================================

CREATE TABLE IF NOT EXISTS seo_domain (
    id BIGSERIAL PRIMARY KEY,

    -- URL identity
    slug VARCHAR(255),
    origin_url TEXT,
    canonical_url TEXT,

    -- Source binding: parcel/layer/document/region/etc.
    ref_type SMALLINT NOT NULL DEFAULT 0,
    ref_id BIGINT,

    -- Operation scope: public/internal/app/admin/static/etc.
    scope VARCHAR(50) NOT NULL DEFAULT 'public',

    -- SEO classify: A=10, B=20, C=30, D=40, N=50
    classify SMALLINT NOT NULL DEFAULT 50,

    -- SEO content registry
    title TEXT,
    description TEXT,
    content TEXT,
    summary TEXT,

    -- Publish lifecycle
    published BOOLEAN NOT NULL DEFAULT FALSE,
    published_at TIMESTAMPTZ,

    -- Search engine policy
    is_site_map BOOLEAN NOT NULL DEFAULT FALSE,
    is_index BOOLEAN NOT NULL DEFAULT FALSE,
    is_robot BOOLEAN DEFAULT FALSE,

    -- Sitemap data
    site_map_lasted_at TIMESTAMPTZ,
    sitemap_priority REAL NOT NULL DEFAULT 0.5,
    sitemap_change_freq VARCHAR(20) NOT NULL DEFAULT 'daily',

    -- Static HTML generation pipeline
    need_generate BOOLEAN NOT NULL DEFAULT TRUE,
    generated_at TIMESTAMPTZ,
    source_updated_at TIMESTAMPTZ,
    static_html_path TEXT,
    static_html_hash VARCHAR(64),

    -- Mobile app deeplink
    deep_link TEXT,

    -- Flexible SEO data: schema JSON-LD, breadcrumb, OG, robots, custom data
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    note TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Existing-table safe ALTERs.
ALTER TABLE seo_domain
    ADD COLUMN IF NOT EXISTS slug VARCHAR(255),
    ADD COLUMN IF NOT EXISTS origin_url TEXT,
    ADD COLUMN IF NOT EXISTS canonical_url TEXT,
    ADD COLUMN IF NOT EXISTS ref_type SMALLINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ref_id BIGINT,
    ADD COLUMN IF NOT EXISTS scope VARCHAR(50) NOT NULL DEFAULT 'public',
    ADD COLUMN IF NOT EXISTS classify SMALLINT NOT NULL DEFAULT 50,
    ADD COLUMN IF NOT EXISTS title TEXT,
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS content TEXT,
    ADD COLUMN IF NOT EXISTS summary TEXT,
    ADD COLUMN IF NOT EXISTS published BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS published_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS is_site_map BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_index BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_robot BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS site_map_lasted_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS sitemap_priority REAL NOT NULL DEFAULT 0.5,
    ADD COLUMN IF NOT EXISTS sitemap_change_freq VARCHAR(20) NOT NULL DEFAULT 'daily',
    ADD COLUMN IF NOT EXISTS need_generate BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS generated_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS source_updated_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS static_html_path TEXT,
    ADD COLUMN IF NOT EXISTS static_html_hash VARCHAR(64),
    ADD COLUMN IF NOT EXISTS deep_link TEXT,
    ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS note TEXT,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- Backfill required identity fields for legacy rows.
UPDATE seo_domain
SET origin_url = '/seo/' || id::text
WHERE origin_url IS NULL OR trim(origin_url) = '';

UPDATE seo_domain
SET canonical_url = origin_url
WHERE canonical_url IS NULL OR trim(canonical_url) = '';

UPDATE seo_domain
SET slug = regexp_replace(
    regexp_replace(
        lower(trim(coalesce(slug, canonical_url, origin_url, 'seo-' || id::text))),
        '[^a-z0-9/_-]+',
        '-',
        'g'
    ),
    '-+',
    '-',
    'g'
)
WHERE slug IS NULL OR trim(slug) = '';

UPDATE seo_domain
SET scope = 'public'
WHERE scope IS NULL OR trim(scope) = '';

UPDATE seo_domain
SET sitemap_change_freq = 'daily'
WHERE sitemap_change_freq IS NULL OR trim(sitemap_change_freq) = '';

UPDATE seo_domain
SET sitemap_priority = 0.5
WHERE sitemap_priority IS NULL;

UPDATE seo_domain
SET metadata = '{}'::jsonb
WHERE metadata IS NULL;

-- Convert metadata to jsonb when legacy column is text/json compatible.
-- If legacy metadata contains invalid JSON, this migration should fail fast;
-- clean those rows before rerunning.
ALTER TABLE seo_domain
ALTER COLUMN metadata TYPE jsonb
USING
CASE
    WHEN metadata IS NULL THEN '{}'::jsonb
    ELSE metadata::jsonb
END;

-- Enforce core constraints after backfill.
ALTER TABLE seo_domain
    ALTER COLUMN slug SET NOT NULL,
    ALTER COLUMN origin_url SET NOT NULL,
    ALTER COLUMN canonical_url SET NOT NULL,
    ALTER COLUMN ref_type SET DEFAULT 0,
    ALTER COLUMN ref_type SET NOT NULL,
    ALTER COLUMN scope SET DEFAULT 'public',
    ALTER COLUMN scope SET NOT NULL,
    ALTER COLUMN classify SET DEFAULT 50,
    ALTER COLUMN classify SET NOT NULL,
    ALTER COLUMN published SET DEFAULT FALSE,
    ALTER COLUMN published SET NOT NULL,
    ALTER COLUMN is_site_map SET DEFAULT FALSE,
    ALTER COLUMN is_site_map SET NOT NULL,
    ALTER COLUMN is_index SET DEFAULT FALSE,
    ALTER COLUMN is_index SET NOT NULL,
    ALTER COLUMN sitemap_priority SET DEFAULT 0.5,
    ALTER COLUMN sitemap_priority SET NOT NULL,
    ALTER COLUMN sitemap_change_freq SET DEFAULT 'daily',
    ALTER COLUMN sitemap_change_freq SET NOT NULL,
    ALTER COLUMN need_generate SET DEFAULT TRUE,
    ALTER COLUMN need_generate SET NOT NULL,
    ALTER COLUMN metadata SET DEFAULT '{}'::jsonb,
    ALTER COLUMN metadata SET NOT NULL,
    ALTER COLUMN created_at SET DEFAULT NOW(),
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN updated_at SET DEFAULT NOW(),
    ALTER COLUMN updated_at SET NOT NULL;

-- Guardrails.
ALTER TABLE seo_domain
    ADD CONSTRAINT chk_seo_domain_classify
    CHECK (classify IN (10, 20, 30, 40, 50)) NOT VALID;

ALTER TABLE seo_domain
    ADD CONSTRAINT chk_seo_domain_sitemap_priority
    CHECK (sitemap_priority >= 0 AND sitemap_priority <= 1) NOT VALID;

ALTER TABLE seo_domain
    ADD CONSTRAINT chk_seo_domain_changefreq
    CHECK (sitemap_change_freq IN ('always', 'hourly', 'daily', 'weekly', 'monthly', 'yearly', 'never')) NOT VALID;

-- Validate constraints if data is already clean.
-- If this fails, fix legacy rows then rerun validation manually.
ALTER TABLE seo_domain VALIDATE CONSTRAINT chk_seo_domain_classify;
ALTER TABLE seo_domain VALIDATE CONSTRAINT chk_seo_domain_sitemap_priority;
ALTER TABLE seo_domain VALIDATE CONSTRAINT chk_seo_domain_changefreq;

-- Indexes.
-- Canonical is the official route. One active canonical must map to one active SEO page.
CREATE UNIQUE INDEX IF NOT EXISTS idx_seo_domain_canonical_url_active
ON seo_domain(canonical_url)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_slug_active
ON seo_domain(slug)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_ref_active
ON seo_domain(ref_type, ref_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_scope_active
ON seo_domain(scope)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_classify_active
ON seo_domain(classify)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_published_sitemap_active
ON seo_domain(published, is_site_map, is_index)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_need_generate_active
ON seo_domain(need_generate)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_source_updated_at_active
ON seo_domain(source_updated_at)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_generated_at_active
ON seo_domain(generated_at)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_updated_at_active
ON seo_domain(updated_at DESC)
WHERE deleted_at IS NULL;

-- Optional rule: enable only if one source entity has exactly one SEO page.
-- If one parcel/layer/document may have multiple landing pages, DO NOT enable this.
-- CREATE UNIQUE INDEX IF NOT EXISTS idx_seo_domain_ref_unique_active
-- ON seo_domain(ref_type, ref_id)
-- WHERE deleted_at IS NULL AND ref_id IS NOT NULL;
