-- ============================================================
-- 000007_seo_lifecycle_render_fields.up.sql
-- Purpose:
--   Add explicit lifecycle/render fields for SEO command architecture.
--   Existing booleans remain for backward compatibility.
-- ============================================================

BEGIN;

ALTER TABLE seo_domain
    ADD COLUMN IF NOT EXISTS page_status VARCHAR(30),
    ADD COLUMN IF NOT EXISTS render_status VARCHAR(30),
    ADD COLUMN IF NOT EXISTS rendered_html TEXT,
    ADD COLUMN IF NOT EXISTS last_render_error TEXT,
    ADD COLUMN IF NOT EXISTS template_key VARCHAR(120),
    ADD COLUMN IF NOT EXISTS template_version VARCHAR(80);

UPDATE seo_domain
SET page_status = CASE
    WHEN deleted_at IS NOT NULL THEN 'archived'
    WHEN published = TRUE THEN 'published'
    ELSE 'draft'
END
WHERE page_status IS NULL OR BTRIM(page_status) = '';

UPDATE seo_domain
SET render_status = CASE
    WHEN need_generate = TRUE THEN 'pending'
    WHEN static_html_hash IS NOT NULL AND BTRIM(static_html_hash) <> '' THEN 'success'
    WHEN generated_at IS NOT NULL THEN 'success'
    ELSE 'none'
END
WHERE render_status IS NULL OR BTRIM(render_status) = '';

UPDATE seo_domain
SET rendered_html = ''
WHERE rendered_html IS NULL;

UPDATE seo_domain
SET last_render_error = ''
WHERE last_render_error IS NULL;

UPDATE seo_domain
SET template_key = ''
WHERE template_key IS NULL;

UPDATE seo_domain
SET template_version = ''
WHERE template_version IS NULL;

ALTER TABLE seo_domain
    ALTER COLUMN page_status SET DEFAULT 'draft',
    ALTER COLUMN page_status SET NOT NULL,
    ALTER COLUMN render_status SET DEFAULT 'none',
    ALTER COLUMN render_status SET NOT NULL;

ALTER TABLE seo_domain
    ADD CONSTRAINT chk_seo_domain_page_status
    CHECK (page_status IN ('draft', 'published', 'archived')) NOT VALID;

ALTER TABLE seo_domain
    ADD CONSTRAINT chk_seo_domain_render_status
    CHECK (render_status IN ('none', 'pending', 'rendering', 'success', 'failed')) NOT VALID;

ALTER TABLE seo_domain VALIDATE CONSTRAINT chk_seo_domain_page_status;
ALTER TABLE seo_domain VALIDATE CONSTRAINT chk_seo_domain_render_status;

CREATE INDEX IF NOT EXISTS idx_seo_domain_page_status
ON seo_domain(page_status)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_render_status
ON seo_domain(render_status)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_template_key
ON seo_domain(template_key)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_page_render_status
ON seo_domain(page_status, render_status)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_render_queue
ON seo_domain(render_status, source_updated_at ASC NULLS FIRST, updated_at ASC)
WHERE deleted_at IS NULL
  AND published = TRUE
  AND need_generate = TRUE;

-- Extend generation log table if it already exists from earlier migrations.
CREATE TABLE IF NOT EXISTS seo_generation_log (
    id BIGSERIAL PRIMARY KEY,
    seo_domain_id BIGINT NOT NULL,
    trigger_type VARCHAR(40) NOT NULL DEFAULT 'manual',
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    static_html_path TEXT,
    static_html_hash VARCHAR(64),
    error_message TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

ALTER TABLE seo_generation_log
    ADD COLUMN IF NOT EXISTS trigger_type VARCHAR(40),
    ADD COLUMN IF NOT EXISTS static_html_path TEXT,
    ADD COLUMN IF NOT EXISTS static_html_hash VARCHAR(64),
    ADD COLUMN IF NOT EXISTS metadata JSONB;

UPDATE seo_generation_log
SET trigger_type = 'manual'
WHERE trigger_type IS NULL OR BTRIM(trigger_type) = '';

UPDATE seo_generation_log
SET metadata = '{}'::jsonb
WHERE metadata IS NULL;

ALTER TABLE seo_generation_log
    ALTER COLUMN trigger_type SET DEFAULT 'manual',
    ALTER COLUMN trigger_type SET NOT NULL,
    ALTER COLUMN metadata SET DEFAULT '{}'::jsonb,
    ALTER COLUMN metadata SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_seo_generation_log_domain_id
ON seo_generation_log(seo_domain_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_generation_log_status
ON seo_generation_log(status)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_generation_log_trigger_type
ON seo_generation_log(trigger_type)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_generation_log_static_hash
ON seo_generation_log(static_html_hash)
WHERE deleted_at IS NULL;

COMMIT;