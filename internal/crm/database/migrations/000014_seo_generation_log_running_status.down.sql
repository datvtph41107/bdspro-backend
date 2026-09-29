BEGIN;

ALTER TABLE seo_generation_log
    DROP CONSTRAINT IF EXISTS chk_seo_generation_log_status;

UPDATE seo_generation_log
SET status = 'processing'
WHERE status = 'running';

ALTER TABLE seo_generation_log
    ADD CONSTRAINT chk_seo_generation_log_status
    CHECK (
        status IN (
            'pending',
            'processing',
            'success',
            'failed',
            'skipped'
        )
    ) NOT VALID;

ALTER TABLE seo_generation_log
    VALIDATE CONSTRAINT chk_seo_generation_log_status;

COMMIT;
