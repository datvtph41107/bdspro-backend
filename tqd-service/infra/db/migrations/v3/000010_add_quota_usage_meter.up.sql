BEGIN;

ALTER TABLE quota_usage_events
    ADD COLUMN IF NOT EXISTS meter_code VARCHAR(128);

-- The pre-migration ledger only contains the report reference slice. Mapping
-- is explicit so a future operation cannot be assigned to a meter by guesswork.
UPDATE quota_usage_events
   SET meter_code = 'workspace.report_generation.accepted'
 WHERE operation = 'workspace.report.generate'
   AND (meter_code IS NULL OR meter_code = '');

DO $$
DECLARE
    unsupported_operations TEXT;
BEGIN
    SELECT STRING_AGG(DISTINCT operation, ', ' ORDER BY operation)
      INTO unsupported_operations
      FROM quota_usage_events
     WHERE meter_code IS NULL OR meter_code = '';

    IF unsupported_operations IS NOT NULL THEN
        RAISE EXCEPTION
            'meter backfill is undefined for operations: %',
            unsupported_operations;
    END IF;
END $$;

ALTER TABLE quota_usage_events
    ALTER COLUMN meter_code SET NOT NULL;

ALTER TABLE quota_usage_events
    ADD CONSTRAINT chk_quota_usage_meter_code
    CHECK (meter_code <> '');

CREATE INDEX IF NOT EXISTS idx_quota_usage_subject_meter_period
    ON quota_usage_events(
        subject_type,
        subject_id,
        meter_code,
        period_start,
        period_end
    );

COMMIT;
