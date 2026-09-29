-- ============================================================
-- 000002_seo_links.up.sql
-- Purpose:
--   Manage rendered internal links between SEO pages.
--   This is the actual link graph that appears in HTML and distributes
--   navigation/authority between guide/detail/search pages.
-- ============================================================

CREATE TABLE IF NOT EXISTS seo_links (
    id BIGSERIAL PRIMARY KEY,

    parent_seo_id BIGINT NOT NULL,
    child_seo_id BIGINT NOT NULL,

    -- Actual URL rendered into HTML. Use child canonical URL if application receives empty link.
    link TEXT NOT NULL DEFAULT '',

    -- Anchor text. Required for SEO operations.
    title TEXT NOT NULL DEFAULT '',

    -- navigation, breadcrumb, related, contextual, footer
    link_type VARCHAR(50) NOT NULL DEFAULT 'navigation',

    -- Higher priority link can be rendered first or preferred by generator.
    priority INTEGER NOT NULL DEFAULT 5,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

ALTER TABLE seo_links
    ADD COLUMN IF NOT EXISTS parent_seo_id BIGINT,
    ADD COLUMN IF NOT EXISTS child_seo_id BIGINT,
    ADD COLUMN IF NOT EXISTS link TEXT,
    ADD COLUMN IF NOT EXISTS title TEXT,
    ADD COLUMN IF NOT EXISTS link_type VARCHAR(50) NOT NULL DEFAULT 'navigation',
    ADD COLUMN IF NOT EXISTS priority INTEGER NOT NULL DEFAULT 5,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

UPDATE seo_links SET link = '' WHERE link IS NULL;
UPDATE seo_links SET title = '' WHERE title IS NULL;
UPDATE seo_links SET link_type = 'navigation' WHERE link_type IS NULL OR trim(link_type) = '';
UPDATE seo_links SET priority = 5 WHERE priority IS NULL;

ALTER TABLE seo_links
    ALTER COLUMN parent_seo_id SET NOT NULL,
    ALTER COLUMN child_seo_id SET NOT NULL,
    ALTER COLUMN link SET DEFAULT '',
    ALTER COLUMN link SET NOT NULL,
    ALTER COLUMN title SET DEFAULT '',
    ALTER COLUMN title SET NOT NULL,
    ALTER COLUMN link_type SET DEFAULT 'navigation',
    ALTER COLUMN link_type SET NOT NULL,
    ALTER COLUMN priority SET DEFAULT 5,
    ALTER COLUMN priority SET NOT NULL,
    ALTER COLUMN created_at SET DEFAULT NOW(),
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN updated_at SET DEFAULT NOW(),
    ALTER COLUMN updated_at SET NOT NULL;

ALTER TABLE seo_links
    ADD CONSTRAINT chk_seo_links_not_self
    CHECK (parent_seo_id <> child_seo_id) NOT VALID;

ALTER TABLE seo_links
    ADD CONSTRAINT chk_seo_links_link_type
    CHECK (link_type IN ('navigation', 'breadcrumb', 'related', 'contextual', 'footer')) NOT VALID;

ALTER TABLE seo_links
    ADD CONSTRAINT chk_seo_links_priority
    CHECK (priority >= 1 AND priority <= 10) NOT VALID;

ALTER TABLE seo_links VALIDATE CONSTRAINT chk_seo_links_not_self;
ALTER TABLE seo_links VALIDATE CONSTRAINT chk_seo_links_link_type;
ALTER TABLE seo_links VALIDATE CONSTRAINT chk_seo_links_priority;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_seo_links_parent_seo'
    ) THEN
        ALTER TABLE seo_links
        ADD CONSTRAINT fk_seo_links_parent_seo
        FOREIGN KEY (parent_seo_id) REFERENCES seo_domain(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_seo_links_child_seo'
    ) THEN
        ALTER TABLE seo_links
        ADD CONSTRAINT fk_seo_links_child_seo
        FOREIGN KEY (child_seo_id) REFERENCES seo_domain(id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_seo_links_parent_active
ON seo_links(parent_seo_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_links_child_active
ON seo_links(child_seo_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_links_type_active
ON seo_links(link_type)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_links_priority_active
ON seo_links(priority DESC)
WHERE deleted_at IS NULL;

-- One active link of the same type between the same SEO pages.
CREATE UNIQUE INDEX IF NOT EXISTS idx_seo_links_pair_type_active
ON seo_links(parent_seo_id, child_seo_id, link_type)
WHERE deleted_at IS NULL;