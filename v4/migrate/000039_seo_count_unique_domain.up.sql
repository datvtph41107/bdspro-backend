-- seo_count may already exist from AutoMigrate without UNIQUE(seo_domain_id),
-- so CREATE TABLE IF NOT EXISTS in 000037 skipped the constraint. ON CONFLICT
-- in RecordClick requires it.

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conrelid = 'seo_count'::regclass
          AND contype IN ('u', 'p')
          AND pg_get_constraintdef(oid) ILIKE '%seo_domain_id%'
    ) AND NOT EXISTS (
        SELECT 1
        FROM pg_indexes
        WHERE tablename = 'seo_count'
          AND indexdef ILIKE '%UNIQUE%seo_domain_id%'
    ) THEN
        -- Collapse duplicates before adding UNIQUE.
        DELETE FROM seo_count a
        USING seo_count b
        WHERE a.seo_domain_id = b.seo_domain_id
          AND a.id > b.id;

        ALTER TABLE seo_count
            ADD CONSTRAINT seo_count_seo_domain_id_key UNIQUE (seo_domain_id);
    END IF;
END $$;
