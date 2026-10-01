-- seo_category may exist from AutoMigrate without PRIMARY KEY / UNIQUE(code),
-- so CREATE TABLE IF NOT EXISTS in 000037 skipped those constraints.
-- Without PK, GROUP BY c.id is invalid in PostgreSQL.

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conrelid = 'seo_category'::regclass AND contype = 'p'
    ) THEN
        ALTER TABLE seo_category ADD CONSTRAINT seo_category_pkey PRIMARY KEY (id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conrelid = 'seo_category'::regclass
          AND contype = 'u'
          AND pg_get_constraintdef(oid) ILIKE '%(code)%'
    ) AND NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE tablename = 'seo_category' AND indexdef ILIKE '%UNIQUE%(%code%)%'
    ) THEN
        ALTER TABLE seo_category ADD CONSTRAINT uq_seo_category_code UNIQUE (code);
    END IF;
END $$;
