-- ============================================================
-- 000003_seo_relative.up.sql
-- Purpose:
--   Manage semantic parent/child and related-page relationships.
--   seo_relative is the semantic map. seo_links is the rendered HTML link graph.
-- ============================================================

CREATE TABLE IF NOT EXISTS seo_relative (
    id BIGSERIAL PRIMARY KEY,

    parent_seo_id BIGINT NOT NULL,
    child_seo_id BIGINT NOT NULL,

    -- parent_child, related, same_region, same_topic
    relation_type VARCHAR(50) NOT NULL DEFAULT 'parent_child',
    priority INTEGER NOT NULL DEFAULT 5,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

ALTER TABLE seo_relative
    ADD COLUMN IF NOT EXISTS parent_seo_id BIGINT,
    ADD COLUMN IF NOT EXISTS child_seo_id BIGINT,
    ADD COLUMN IF NOT EXISTS relation_type VARCHAR(50) NOT NULL DEFAULT 'parent_child',
    ADD COLUMN IF NOT EXISTS priority INTEGER NOT NULL DEFAULT 5,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

UPDATE seo_relative SET relation_type = 'parent_child' WHERE relation_type IS NULL OR trim(relation_type) = '';
UPDATE seo_relative SET priority = 5 WHERE priority IS NULL;

ALTER TABLE seo_relative
    ALTER COLUMN parent_seo_id SET NOT NULL,
    ALTER COLUMN child_seo_id SET NOT NULL,
    ALTER COLUMN relation_type SET DEFAULT 'parent_child',
    ALTER COLUMN relation_type SET NOT NULL,
    ALTER COLUMN priority SET DEFAULT 5,
    ALTER COLUMN priority SET NOT NULL,
    ALTER COLUMN created_at SET DEFAULT NOW(),
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN updated_at SET DEFAULT NOW(),
    ALTER COLUMN updated_at SET NOT NULL;

ALTER TABLE seo_relative
    ADD CONSTRAINT chk_seo_relative_not_self
    CHECK (parent_seo_id <> child_seo_id) NOT VALID;

ALTER TABLE seo_relative
    ADD CONSTRAINT chk_seo_relative_type
    CHECK (relation_type IN ('parent_child', 'related', 'same_region', 'same_topic')) NOT VALID;

ALTER TABLE seo_relative
    ADD CONSTRAINT chk_seo_relative_priority
    CHECK (priority >= 1 AND priority <= 10) NOT VALID;

ALTER TABLE seo_relative VALIDATE CONSTRAINT chk_seo_relative_not_self;
ALTER TABLE seo_relative VALIDATE CONSTRAINT chk_seo_relative_type;
ALTER TABLE seo_relative VALIDATE CONSTRAINT chk_seo_relative_priority;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_seo_relative_parent_seo'
    ) THEN
        ALTER TABLE seo_relative
        ADD CONSTRAINT fk_seo_relative_parent_seo
        FOREIGN KEY (parent_seo_id) REFERENCES seo_domain(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_seo_relative_child_seo'
    ) THEN
        ALTER TABLE seo_relative
        ADD CONSTRAINT fk_seo_relative_child_seo
        FOREIGN KEY (child_seo_id) REFERENCES seo_domain(id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_seo_relative_parent_active
ON seo_relative(parent_seo_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_relative_child_active
ON seo_relative(child_seo_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_relative_type_active
ON seo_relative(relation_type)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_relative_priority_active
ON seo_relative(priority DESC)
WHERE deleted_at IS NULL;

-- One active semantic relationship of a type between the same pages.
CREATE UNIQUE INDEX IF NOT EXISTS idx_seo_relative_pair_type_active
ON seo_relative(parent_seo_id, child_seo_id, relation_type)
WHERE deleted_at IS NULL;