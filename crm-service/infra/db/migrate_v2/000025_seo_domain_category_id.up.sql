BEGIN;

ALTER TABLE seo_domain
    ADD COLUMN IF NOT EXISTS category_id BIGINT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_seo_domain_category'
    ) THEN
        ALTER TABLE seo_domain
            ADD CONSTRAINT fk_seo_domain_category
            FOREIGN KEY (category_id) REFERENCES seo_category(id)
            ON DELETE SET NULL;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_seo_domain_category_id
    ON seo_domain(category_id)
    WHERE deleted_at IS NULL AND category_id IS NOT NULL;

-- Backfill from primary category assignment.
UPDATE seo_domain sd
SET category_id = a.category_id
FROM seo_content_category_assignment a
WHERE sd.category_id IS NULL
  AND a.seo_domain_id = sd.id
  AND a.deleted_at IS NULL
  AND a.is_primary = TRUE;

COMMENT ON COLUMN seo_domain.category_id IS
    'Primary seo_category id denormalized for fast filter/response on planning news.';

COMMIT;
