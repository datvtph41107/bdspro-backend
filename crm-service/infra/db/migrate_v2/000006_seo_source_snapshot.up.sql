-- ============================================================
-- 000006_seo_source_snapshot.up.sql
-- Purpose:
--   Extend seo_domain with source-aware binding state.
--   This keeps SEO pages attached to normalized source snapshots so Studio,
--   generation, touch-by-ref and stale detection can use the same contract.
-- ============================================================

ALTER TABLE seo_domain
    ADD COLUMN IF NOT EXISTS ref_source VARCHAR(80),
    ADD COLUMN IF NOT EXISTS ref_label TEXT,
    ADD COLUMN IF NOT EXISTS ref_url TEXT,
    ADD COLUMN IF NOT EXISTS source_status VARCHAR(40),
    ADD COLUMN IF NOT EXISTS ref_missing BOOLEAN,
    ADD COLUMN IF NOT EXISTS ref_snapshot_json JSONB,
    ADD COLUMN IF NOT EXISTS ref_hash VARCHAR(64),
    ADD COLUMN IF NOT EXISTS ref_last_synced_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS quality_score INTEGER;

UPDATE seo_domain
SET ref_source = CASE
    WHEN ref_type IS NULL OR ref_type = 0 THEN 'manual'
    WHEN ref_type = 100 THEN 'bdspro.region'
    WHEN ref_type = 200 THEN 'tqd.parcel'
    WHEN ref_type = 300 THEN 'bdspro.area_region'
    WHEN ref_type = 400 THEN 'bdspro.project'
    WHEN ref_type = 500 THEN 'tqd.legal_document'
    WHEN ref_type = 600 THEN 'tqd.qh_layer'
    WHEN ref_type = 700 THEN 'crm.report'
    WHEN ref_source IS NULL OR BTRIM(ref_source) = '' THEN ''
    ELSE ref_source
END
WHERE ref_source IS NULL OR BTRIM(ref_source) = '';

UPDATE seo_domain
SET source_status = CASE
    WHEN ref_type IS NULL OR ref_type = 0 THEN 'manual'
    WHEN source_status IS NULL OR BTRIM(source_status) = '' THEN 'linked'
    ELSE source_status
END
WHERE source_status IS NULL OR BTRIM(source_status) = '';

UPDATE seo_domain
SET ref_missing = FALSE
WHERE ref_missing IS NULL;

UPDATE seo_domain
SET ref_snapshot_json = '{}'::jsonb
WHERE ref_snapshot_json IS NULL;

UPDATE seo_domain
SET quality_score = 0
WHERE quality_score IS NULL OR quality_score < 0 OR quality_score > 100;

ALTER TABLE IF EXISTS seo_domain DROP CONSTRAINT IF EXISTS chk_seo_domain_source_status;
ALTER TABLE IF EXISTS seo_domain DROP CONSTRAINT IF EXISTS chk_seo_domain_quality_score;

ALTER TABLE seo_domain
    ALTER COLUMN ref_source SET DEFAULT 'manual',
    ALTER COLUMN ref_source SET NOT NULL,
    ALTER COLUMN source_status SET DEFAULT 'manual',
    ALTER COLUMN source_status SET NOT NULL,
    ALTER COLUMN ref_missing SET DEFAULT FALSE,
    ALTER COLUMN ref_missing SET NOT NULL,
    ALTER COLUMN ref_snapshot_json SET DEFAULT '{}'::jsonb,
    ALTER COLUMN ref_snapshot_json SET NOT NULL,
    ALTER COLUMN quality_score SET DEFAULT 0,
    ALTER COLUMN quality_score SET NOT NULL;

ALTER TABLE seo_domain
    ADD CONSTRAINT chk_seo_domain_source_status
    CHECK (source_status IN ('manual', 'linked', 'stale', 'missing', 'permission_denied')) NOT VALID;

ALTER TABLE seo_domain
    ADD CONSTRAINT chk_seo_domain_quality_score
    CHECK (quality_score >= 0 AND quality_score <= 100) NOT VALID;

ALTER TABLE seo_domain VALIDATE CONSTRAINT chk_seo_domain_source_status;
ALTER TABLE seo_domain VALIDATE CONSTRAINT chk_seo_domain_quality_score;

CREATE INDEX IF NOT EXISTS idx_seo_domain_ref_source_active
ON seo_domain(ref_type, ref_source, ref_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_source_status_active
ON seo_domain(source_status)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_ref_hash_active
ON seo_domain(ref_hash)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_ref_last_synced_active
ON seo_domain(ref_last_synced_at)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_quality_score_active
ON seo_domain(quality_score)
WHERE deleted_at IS NULL;

-- Optional strict unique rule for environments that enforce one SEO page per source.
-- Enable after deduplicating legacy rows:
-- CREATE UNIQUE INDEX IF NOT EXISTS uniq_seo_domain_active_ref_source
-- ON seo_domain(ref_type, ref_source, ref_id)
-- WHERE deleted_at IS NULL AND ref_id IS NOT NULL;