-- ============================================================
-- 000005_seo_generation_log.up.sql
-- Purpose:
--   Track nightly/static generation jobs.
--   If a page fails generation, need_generate remains true and this table
--   explains why, making operations retryable and debuggable.
-- ============================================================

CREATE TABLE IF NOT EXISTS seo_generation_log (
    id BIGSERIAL PRIMARY KEY,

    seo_domain_id BIGINT NOT NULL,
    job_id VARCHAR(100),

    -- generate_html, rebuild_sitemap, revalidate, render_data, etc.
    action VARCHAR(50) NOT NULL DEFAULT 'generate_html',

    -- pending, processing, success, failed, skipped
    status VARCHAR(30) NOT NULL DEFAULT 'pending',

    error_message TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    duration_ms BIGINT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

ALTER TABLE seo_generation_log
    ADD COLUMN IF NOT EXISTS seo_domain_id BIGINT,
    ADD COLUMN IF NOT EXISTS job_id VARCHAR(100),
    ADD COLUMN IF NOT EXISTS action VARCHAR(50) NOT NULL DEFAULT 'generate_html',
    ADD COLUMN IF NOT EXISTS status VARCHAR(30) NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS error_message TEXT,
    ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS started_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS finished_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS duration_ms BIGINT,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

UPDATE seo_generation_log SET action = 'generate_html' WHERE action IS NULL OR trim(action) = '';
UPDATE seo_generation_log SET status = 'pending' WHERE status IS NULL OR trim(status) = '';
UPDATE seo_generation_log SET metadata = '{}'::jsonb WHERE metadata IS NULL;

ALTER TABLE seo_generation_log
    ALTER COLUMN seo_domain_id SET NOT NULL,
    ALTER COLUMN action SET DEFAULT 'generate_html',
    ALTER COLUMN action SET NOT NULL,
    ALTER COLUMN status SET DEFAULT 'pending',
    ALTER COLUMN status SET NOT NULL,
    ALTER COLUMN metadata SET DEFAULT '{}'::jsonb,
    ALTER COLUMN metadata SET NOT NULL,
    ALTER COLUMN created_at SET DEFAULT NOW(),
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN updated_at SET DEFAULT NOW(),
    ALTER COLUMN updated_at SET NOT NULL;

ALTER TABLE seo_generation_log
    ADD CONSTRAINT chk_seo_generation_log_status
    CHECK (status IN ('pending', 'processing', 'success', 'failed', 'skipped')) NOT VALID;

ALTER TABLE seo_generation_log
    ADD CONSTRAINT chk_seo_generation_log_duration
    CHECK (duration_ms IS NULL OR duration_ms >= 0) NOT VALID;

ALTER TABLE seo_generation_log VALIDATE CONSTRAINT chk_seo_generation_log_status;
ALTER TABLE seo_generation_log VALIDATE CONSTRAINT chk_seo_generation_log_duration;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_seo_generation_log_domain'
    ) THEN
        ALTER TABLE seo_generation_log
        ADD CONSTRAINT fk_seo_generation_log_domain
        FOREIGN KEY (seo_domain_id) REFERENCES seo_domain(id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_seo_generation_log_domain_active
ON seo_generation_log(seo_domain_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_generation_log_job_id_active
ON seo_generation_log(job_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_generation_log_status_active
ON seo_generation_log(status)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_generation_log_created_at_active
ON seo_generation_log(created_at DESC)
WHERE deleted_at IS NULL;