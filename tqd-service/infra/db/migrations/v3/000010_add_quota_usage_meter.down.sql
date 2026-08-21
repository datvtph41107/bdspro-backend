BEGIN;

DROP INDEX IF EXISTS idx_quota_usage_subject_meter_period;

ALTER TABLE quota_usage_events
    DROP CONSTRAINT IF EXISTS chk_quota_usage_meter_code;

ALTER TABLE quota_usage_events
    DROP COLUMN IF EXISTS meter_code;

COMMIT;
