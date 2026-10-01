-- Ensure NOT NULL jsonb columns accept inserts when app omits values.
UPDATE seo_domain
SET ref_snapshot_json = '{}'::jsonb
WHERE ref_snapshot_json IS NULL;

UPDATE seo_domain
SET metadata = '{}'::jsonb
WHERE metadata IS NULL;

ALTER TABLE seo_domain
    ALTER COLUMN ref_snapshot_json SET DEFAULT '{}'::jsonb,
    ALTER COLUMN metadata SET DEFAULT '{}'::jsonb;
