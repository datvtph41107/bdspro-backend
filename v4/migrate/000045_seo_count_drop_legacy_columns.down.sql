-- Legacy CRM columns are restored as nullable so the 000037 schema keeps working.
ALTER TABLE seo_count_log
    ADD COLUMN IF NOT EXISTS seo_id BIGINT,
    ADD COLUMN IF NOT EXISTS clicked BIGINT,
    ADD COLUMN IF NOT EXISTS "timestamp" TIMESTAMPTZ;

ALTER TABLE seo_count
    ADD COLUMN IF NOT EXISTS seo_id BIGINT,
    ADD COLUMN IF NOT EXISTS count BIGINT;
