BEGIN;

ALTER TABLE seo_generation_log
    DROP CONSTRAINT IF EXISTS chk_seo_generation_log_status;

-- Chuẩn hóa dữ liệu cũ theo domain vocabulary hiện tại.
UPDATE seo_generation_log
SET status = 'running'
WHERE status = 'processing';

ALTER TABLE seo_generation_log
    ADD CONSTRAINT chk_seo_generation_log_status
    CHECK (
        status IN (
            'pending',
            'running',
            'success',
            'failed',
            'skipped'
        )
    ) NOT VALID;

ALTER TABLE seo_generation_log
    VALIDATE CONSTRAINT chk_seo_generation_log_status;

COMMIT;
