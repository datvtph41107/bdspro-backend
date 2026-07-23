-- 000022_report_images.up.sql
ALTER TABLE reports
    ADD COLUMN IF NOT EXISTS images JSONB NOT NULL DEFAULT '[]'::jsonb;
