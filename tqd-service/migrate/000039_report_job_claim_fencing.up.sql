BEGIN;

ALTER TABLE report_jobs
    ADD COLUMN IF NOT EXISTS claim_version BIGINT NOT NULL DEFAULT 0;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
          FROM pg_constraint
         WHERE conname = 'chk_report_jobs_claim_version'
           AND conrelid = 'report_jobs'::regclass
    ) THEN
        ALTER TABLE report_jobs
            ADD CONSTRAINT chk_report_jobs_claim_version
            CHECK (claim_version >= 0);
    END IF;
END $$;

COMMIT;
