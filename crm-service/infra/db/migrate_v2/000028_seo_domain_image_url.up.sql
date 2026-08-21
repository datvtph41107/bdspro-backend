-- 000028_seo_domain_image_url.up.sql
-- Cover / thumbnail image for SEO pages (list + detail).

ALTER TABLE seo_domain
    ADD COLUMN IF NOT EXISTS image_url TEXT;

COMMENT ON COLUMN seo_domain.image_url IS
    'Public image/cover URL for SEO page list and detail responses';
