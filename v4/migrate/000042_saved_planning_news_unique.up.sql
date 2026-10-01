-- user_saved_planning_news may already exist from AutoMigrate without
-- UNIQUE(user_id, seo_domain_id), so CREATE TABLE IF NOT EXISTS in 000037
-- skipped the constraint. ON CONFLICT in SavePlanningNews requires it.

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conrelid = 'user_saved_planning_news'::regclass AND contype = 'p'
    ) THEN
        -- Ensure ids are unique before adding PK (AutoMigrate tables can have dup ids).
        IF EXISTS (SELECT id FROM user_saved_planning_news GROUP BY id HAVING COUNT(*) > 1) THEN
            CREATE TEMP TABLE uspn_fix AS
            SELECT user_id, seo_domain_id, created_at, updated_at, deleted_at
            FROM user_saved_planning_news;
            TRUNCATE user_saved_planning_news RESTART IDENTITY;
            INSERT INTO user_saved_planning_news (user_id, seo_domain_id, created_at, updated_at, deleted_at)
            SELECT user_id, seo_domain_id, created_at, updated_at, deleted_at FROM uspn_fix;
            DROP TABLE uspn_fix;
        END IF;
        ALTER TABLE user_saved_planning_news
            ADD CONSTRAINT user_saved_planning_news_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conrelid = 'user_saved_planning_news'::regclass
          AND contype = 'u'
          AND pg_get_constraintdef(oid) ILIKE '%user_id%'
          AND pg_get_constraintdef(oid) ILIKE '%seo_domain_id%'
    ) AND NOT EXISTS (
        SELECT 1
        FROM pg_indexes
        WHERE tablename = 'user_saved_planning_news'
          AND indexdef ILIKE '%UNIQUE%'
          AND indexdef ILIKE '%user_id%'
          AND indexdef ILIKE '%seo_domain_id%'
    ) THEN
        -- Keep newest row per (user_id, seo_domain_id).
        DELETE FROM user_saved_planning_news a
        USING user_saved_planning_news b
        WHERE a.user_id = b.user_id
          AND a.seo_domain_id = b.seo_domain_id
          AND a.id > b.id;

        ALTER TABLE user_saved_planning_news
            ADD CONSTRAINT user_saved_planning_news_user_domain_key UNIQUE (user_id, seo_domain_id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_saved_news_user
    ON user_saved_planning_news (user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_saved_news_domain
    ON user_saved_planning_news (seo_domain_id);
