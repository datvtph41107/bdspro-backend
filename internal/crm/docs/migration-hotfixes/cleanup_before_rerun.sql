-- Run this only if the database is currently stuck because legacy seo_domain rows contain NULLs.
-- It cleans legacy rows before rerunning migrate up.

BEGIN;

UPDATE seo_domain
SET ref_type = 0
WHERE ref_type IS NULL;

UPDATE seo_domain
SET origin_url = '/seo/' || id::text
WHERE origin_url IS NULL OR BTRIM(origin_url) = '';

UPDATE seo_domain
SET canonical_url = origin_url
WHERE canonical_url IS NULL OR BTRIM(canonical_url) = '';

UPDATE seo_domain
SET slug = CONCAT('seo-domain-', id)
WHERE slug IS NULL OR BTRIM(slug) = '';

UPDATE seo_domain
SET scope = 'public'
WHERE scope IS NULL OR BTRIM(scope) = '';

UPDATE seo_domain
SET classify = 50
WHERE classify IS NULL OR classify NOT IN (10,20,30,40,50);

UPDATE seo_domain
SET published = FALSE
WHERE published IS NULL;

UPDATE seo_domain
SET is_site_map = FALSE
WHERE is_site_map IS NULL;

UPDATE seo_domain
SET is_index = FALSE
WHERE is_index IS NULL;

UPDATE seo_domain
SET is_robot = FALSE
WHERE is_robot IS NULL;

UPDATE seo_domain
SET sitemap_priority = 0.5
WHERE sitemap_priority IS NULL OR sitemap_priority < 0 OR sitemap_priority > 1;

UPDATE seo_domain
SET sitemap_change_freq = 'daily'
WHERE sitemap_change_freq IS NULL
   OR BTRIM(sitemap_change_freq) = ''
   OR sitemap_change_freq NOT IN ('always', 'hourly', 'daily', 'weekly', 'monthly', 'yearly', 'never');

UPDATE seo_domain
SET need_generate = TRUE
WHERE need_generate IS NULL;

UPDATE seo_domain
SET metadata = '{}'::jsonb
WHERE metadata IS NULL;

UPDATE seo_domain
SET created_at = NOW()
WHERE created_at IS NULL;

UPDATE seo_domain
SET updated_at = NOW()
WHERE updated_at IS NULL;

COMMIT;