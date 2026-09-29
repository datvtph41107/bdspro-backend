BEGIN;

ALTER TABLE report_jobs
    DROP CONSTRAINT IF EXISTS chk_report_jobs_claim_version;

ALTER TABLE report_jobs
    DROP COLUMN IF EXISTS claim_version;

COMMIT;
