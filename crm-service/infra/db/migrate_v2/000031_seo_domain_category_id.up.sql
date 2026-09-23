BEGIN;

-- Restore fields still used by the technical SEO page API after
-- 000024_simplify_seo_to_pages dropped strategy/taxonomy columns.
ALTER TABLE seo_domain
    ADD COLUMN IF NOT EXISTS visible_on SMALLINT NOT NULL DEFAULT 30;

ALTER TABLE seo_domain
    DROP CONSTRAINT IF EXISTS chk_seo_domain_visible_on;
ALTER TABLE seo_domain
    ADD CONSTRAINT chk_seo_domain_visible_on
    CHECK (visible_on IN (10, 20, 30, 40, 50)) NOT VALID;
ALTER TABLE seo_domain VALIDATE CONSTRAINT chk_seo_domain_visible_on;

CREATE INDEX IF NOT EXISTS idx_seo_domain_visible_on_public
    ON seo_domain(visible_on, published_at DESC, id DESC)
    WHERE deleted_at IS NULL AND published = TRUE AND scope = 'public';

COMMENT ON COLUMN seo_domain.visible_on IS
    'Surface visibility for planning news (10=web, 20=app, 30=both, ...).';

-- Denormalized primary chủ đề (seo_category.id) for list/filter.
ALTER TABLE seo_domain
    ADD COLUMN IF NOT EXISTS category_id BIGINT NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_category_id
    ON seo_domain(category_id)
    WHERE deleted_at IS NULL AND category_id IS NOT NULL;

COMMENT ON COLUMN seo_domain.category_id IS
    'Primary chủ đề (seo_category.id) denormalized for planning news list/filter.';

COMMIT;
