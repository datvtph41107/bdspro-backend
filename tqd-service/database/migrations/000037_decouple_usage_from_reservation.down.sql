BEGIN;

-- Downgrade is safe only while every durable Usage still belongs to the old
-- Redis reservation lifecycle. Never fabricate reservation identifiers.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
          FROM quota_usage_events
         WHERE reservation_id IS NULL OR reservation_id = ''
    ) THEN
        RAISE EXCEPTION
            'cannot restore mandatory reservation_id: PostgreSQL-only usage rows exist';
    END IF;
END $$;

DROP INDEX IF EXISTS uidx_quota_usage_reservation;

ALTER TABLE quota_usage_events
    DROP CONSTRAINT IF EXISTS chk_quota_usage_commercial_evidence,
    DROP CONSTRAINT IF EXISTS chk_quota_usage_limit_snapshot,
    DROP COLUMN IF EXISTS limit_snapshot,
    DROP COLUMN IF EXISTS policy_version,
    DROP COLUMN IF EXISTS plan_version,
    DROP COLUMN IF EXISTS plan_code,
    DROP COLUMN IF EXISTS subscription_id;

ALTER TABLE quota_usage_events
    ALTER COLUMN reservation_id SET NOT NULL;

ALTER TABLE quota_usage_events
    ADD CONSTRAINT quota_usage_events_reservation_id_check
        CHECK (reservation_id <> ''),
    ADD CONSTRAINT quota_usage_events_reservation_id_key
        UNIQUE (reservation_id);

COMMIT;
