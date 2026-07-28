BEGIN;

-- Public Planning News feed classification owned by CRM.
ALTER TABLE seo_domain
    ADD COLUMN IF NOT EXISTS visible_on SMALLINT NOT NULL DEFAULT 30;

ALTER TABLE seo_domain
    DROP CONSTRAINT IF EXISTS chk_seo_domain_visible_on;
ALTER TABLE seo_domain
    ADD CONSTRAINT chk_seo_domain_visible_on
    CHECK (visible_on IN (10, 20, 30, 40, 50)) NOT VALID;
ALTER TABLE seo_domain VALIDATE CONSTRAINT chk_seo_domain_visible_on;

COMMENT ON COLUMN seo_domain.visible_on IS
'Planning News feed: 10=featured, 20=attention, 30=latest, 40=most_read, 50=breaking.';

CREATE INDEX IF NOT EXISTS idx_seo_domain_visible_on_public
    ON seo_domain(visible_on, published_at DESC, id DESC)
    WHERE deleted_at IS NULL AND published = TRUE AND scope = 'public';

-- Task compatibility: seo_category is the canonical physical taxonomy table.
-- Migration 000012 originally created seo_content_category. Rename the same
-- physical table instead of creating a duplicate source of truth. The former
-- name remains as an automatically-updatable compatibility view.
DO $$
DECLARE
    seo_category_kind "char";
BEGIN
    SELECT relkind INTO seo_category_kind
    FROM pg_class
    WHERE oid = to_regclass('public.seo_category');

    IF seo_category_kind = 'v' THEN
        DROP VIEW public.seo_category;
        seo_category_kind := NULL;
    END IF;

    IF to_regclass('public.seo_category') IS NULL
       AND to_regclass('public.seo_content_category') IS NOT NULL THEN
        ALTER TABLE public.seo_content_category RENAME TO seo_category;
    END IF;
END $$;

CREATE OR REPLACE VIEW public.seo_content_category AS
SELECT
    id,
    domain_id,
    code,
    slug,
    name,
    description,
    status,
    sort_order,
    metadata,
    created_at,
    updated_at,
    deleted_at
FROM public.seo_category;

COMMENT ON TABLE public.seo_category IS
'Canonical SEO/content topic table. seo_domain pages are linked through seo_content_category_assignment.';
COMMENT ON VIEW public.seo_content_category IS
'Compatibility name for the canonical seo_category table.';

CREATE OR REPLACE VIEW public.seo_domain_category AS
SELECT
    a.id,
    a.seo_domain_id,
    a.category_id,
    a.is_primary,
    a.source_system,
    a.source_type,
    a.source_id,
    a.created_at,
    a.deleted_at
FROM public.seo_content_category_assignment a;

COMMIT;
