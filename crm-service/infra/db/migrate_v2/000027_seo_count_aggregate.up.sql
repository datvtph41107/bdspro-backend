-- Upgrade path when 000026 previously created the old single seo_count (clicked/timestamp) table.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'seo_count'
          AND column_name = 'clicked'
    ) AND NOT EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_schema = 'public'
          AND table_name = 'seo_count_log'
    ) THEN
        ALTER TABLE seo_count RENAME TO seo_count_log;
        ALTER INDEX IF EXISTS idx_seo_count_seo_id_timestamp
            RENAME TO idx_seo_count_log_seo_id_timestamp;
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS seo_count_log (
    seo_id BIGINT NOT NULL,
    clicked BIGINT NOT NULL DEFAULT 0,
    "timestamp" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_seo_count_log_seo_id_timestamp
    ON seo_count_log(seo_id, "timestamp" DESC);

CREATE TABLE IF NOT EXISTS seo_count (
    seo_id BIGINT PRIMARY KEY,
    count BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_seo_count_updated_at
    ON seo_count(updated_at DESC);

-- Backfill aggregate from existing logs when empty.
INSERT INTO seo_count (seo_id, count, updated_at)
SELECT seo_id, SUM(clicked), MAX("timestamp")
FROM seo_count_log
GROUP BY seo_id
ON CONFLICT (seo_id) DO NOTHING;
