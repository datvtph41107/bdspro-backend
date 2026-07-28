BEGIN;

DROP VIEW IF EXISTS public.seo_domain_category;
DROP VIEW IF EXISTS public.seo_content_category;

DO $$
BEGIN
    IF to_regclass('public.seo_category') IS NOT NULL
       AND to_regclass('public.seo_content_category') IS NULL THEN
        ALTER TABLE public.seo_category RENAME TO seo_content_category;
    END IF;
END $$;

DROP INDEX IF EXISTS idx_seo_domain_visible_on_public;
ALTER TABLE seo_domain DROP CONSTRAINT IF EXISTS chk_seo_domain_visible_on;
ALTER TABLE seo_domain DROP COLUMN IF EXISTS visible_on;

COMMIT;
