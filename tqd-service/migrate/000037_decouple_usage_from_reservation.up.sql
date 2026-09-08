BEGIN;

-- Durable accepted usage must not depend on the runtime admission mechanism.
-- Existing Redis-D rows keep their reservation_id; future PostgreSQL-only
-- acceptance may persist NULL here without inventing a fake reservation.
ALTER TABLE quota_usage_events
    ALTER COLUMN reservation_id DROP NOT NULL;

ALTER TABLE quota_usage_events
    DROP CONSTRAINT IF EXISTS quota_usage_events_reservation_id_key;

ALTER TABLE quota_usage_events
    DROP CONSTRAINT IF EXISTS quota_usage_events_reservation_id_check;

CREATE UNIQUE INDEX IF NOT EXISTS uidx_quota_usage_reservation
    ON quota_usage_events(reservation_id)
    WHERE reservation_id IS NOT NULL AND reservation_id <> '';

ALTER TABLE quota_usage_events
    ADD COLUMN IF NOT EXISTS subscription_id BIGINT,
    ADD COLUMN IF NOT EXISTS plan_code VARCHAR(128),
    ADD COLUMN IF NOT EXISTS plan_version VARCHAR(64),
    ADD COLUMN IF NOT EXISTS policy_version VARCHAR(64),
    ADD COLUMN IF NOT EXISTS limit_snapshot BIGINT;

ALTER TABLE quota_usage_events
    ADD CONSTRAINT chk_quota_usage_limit_snapshot
        CHECK (limit_snapshot IS NULL OR limit_snapshot >= 0),
    ADD CONSTRAINT chk_quota_usage_commercial_evidence
        CHECK (
            (subscription_id IS NULL
                AND plan_code IS NULL
                AND plan_version IS NULL
                AND policy_version IS NULL
                AND limit_snapshot IS NULL)
            OR
            (subscription_id IS NOT NULL AND subscription_id > 0
                AND plan_code IS NOT NULL AND plan_code <> ''
                AND plan_version IS NOT NULL AND plan_version <> ''
                AND policy_version IS NOT NULL AND policy_version <> ''
                AND limit_snapshot IS NOT NULL)
        );

COMMIT;
