BEGIN;

ALTER TABLE legacy_membership_subscription_mappings
    ADD COLUMN IF NOT EXISTS decision JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS reviewed_by BIGINT,
    ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS blocked_reason TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS last_run_id VARCHAR(128) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS last_error TEXT NOT NULL DEFAULT '';

ALTER TABLE legacy_membership_subscription_mappings
    DROP CONSTRAINT IF EXISTS chk_legacy_membership_subscription_mapping_status;
ALTER TABLE legacy_membership_subscription_mappings
    ADD CONSTRAINT chk_legacy_membership_subscription_mapping_status
        CHECK (
            status IN (
                'pending_review',
                'reviewed',
                'draft_created',
                'shadow_read',
                'cutover',
                'rolled_back',
                'blocked'
            )
        );

ALTER TABLE legacy_membership_subscription_mappings
    DROP CONSTRAINT IF EXISTS chk_legacy_membership_subscription_mapping_decision,
    DROP CONSTRAINT IF EXISTS chk_legacy_membership_subscription_mapping_review,
    DROP CONSTRAINT IF EXISTS chk_legacy_membership_subscription_mapping_migrated,
    DROP CONSTRAINT IF EXISTS chk_legacy_membership_subscription_mapping_blocked,
    DROP CONSTRAINT IF EXISTS chk_legacy_membership_subscription_mapping_run_id;

ALTER TABLE legacy_membership_subscription_mappings
    ADD CONSTRAINT chk_legacy_membership_subscription_mapping_decision
        CHECK (jsonb_typeof(decision) = 'object') NOT VALID;
ALTER TABLE legacy_membership_subscription_mappings
    ADD CONSTRAINT chk_legacy_membership_subscription_mapping_review
        CHECK (
            status NOT IN ('reviewed', 'shadow_read', 'cutover')
            OR (
                reviewed_by IS NOT NULL
                AND reviewed_at IS NOT NULL
                AND decision <> '{}'::jsonb
            )
        ) NOT VALID;
ALTER TABLE legacy_membership_subscription_mappings
    ADD CONSTRAINT chk_legacy_membership_subscription_mapping_migrated
        CHECK (
            status NOT IN ('shadow_read', 'cutover')
            OR (
                subscription_id IS NOT NULL
                AND migrated_at IS NOT NULL
                AND last_run_id <> ''
            )
        ) NOT VALID;
ALTER TABLE legacy_membership_subscription_mappings
    ADD CONSTRAINT chk_legacy_membership_subscription_mapping_blocked
        CHECK (
            (status = 'blocked' AND LENGTH(BTRIM(blocked_reason)) > 0)
            OR (status <> 'blocked' AND blocked_reason = '')
        ) NOT VALID;
ALTER TABLE legacy_membership_subscription_mappings
    ADD CONSTRAINT chk_legacy_membership_subscription_mapping_run_id
        CHECK (
            last_run_id = ''
            OR last_run_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$'
        ) NOT VALID;

ALTER TABLE catalog_subscription_events
    DROP CONSTRAINT IF EXISTS chk_catalog_subscription_legacy_backfill_evidence;
ALTER TABLE catalog_subscription_events
    ADD CONSTRAINT chk_catalog_subscription_legacy_backfill_evidence
        CHECK (
            event_type <> 'migration.legacy_membership_backfilled'
            OR (
                jsonb_typeof(metadata) = 'object'
                AND metadata ? 'legacyMembershipId'
                AND metadata ? 'legacyPlanId'
                AND metadata ? 'sourceChecksum'
                AND metadata ? 'planTermsChecksum'
                AND metadata ->> 'legacyMembershipId' ~ '^[1-9][0-9]*$'
                AND metadata ->> 'legacyPlanId' ~ '^[1-9][0-9]*$'
                AND metadata ->> 'sourceChecksum' ~ '^[a-f0-9]{64}$'
                AND metadata ->> 'planTermsChecksum' ~ '^[a-f0-9]{64}$'
            )
        ) NOT VALID;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM catalog_subscription_events
        WHERE event_type = 'migration.legacy_membership_backfilled'
        GROUP BY subscription_id
        HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION
            'duplicate legacy membership backfill events must be reconciled first';
    END IF;
END;
$$;

CREATE UNIQUE INDEX IF NOT EXISTS uq_catalog_subscription_legacy_backfill_event
    ON catalog_subscription_events(subscription_id, event_type)
    WHERE event_type = 'migration.legacy_membership_backfilled';

CREATE OR REPLACE FUNCTION qhpro_catalog_validate_membership_mapping_transition()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.status NOT IN ('pending_review', 'reviewed', 'blocked') THEN
            RAISE EXCEPTION
                'new membership mapping must start pending, reviewed, or blocked';
        END IF;
        IF NEW.subscription_id IS NOT NULL OR NEW.migrated_at IS NOT NULL THEN
            RAISE EXCEPTION
                'new membership mapping cannot start migrated';
        END IF;
        RETURN NEW;
    END IF;

    IF OLD.subscription_id IS NOT NULL
       AND NEW.subscription_id IS DISTINCT FROM OLD.subscription_id THEN
        RAISE EXCEPTION
            'membership mapping subscription_id is immutable';
    END IF;

    IF OLD.migrated_at IS NOT NULL
       AND NEW.migrated_at IS DISTINCT FROM OLD.migrated_at THEN
        RAISE EXCEPTION
            'membership mapping migrated_at is immutable';
    END IF;

    IF (OLD.status IN ('reviewed', 'shadow_read', 'cutover', 'rolled_back')
        OR OLD.subscription_id IS NOT NULL)
       AND (
            NEW.source_checksum IS DISTINCT FROM OLD.source_checksum
            OR NEW.decision IS DISTINCT FROM OLD.decision
            OR NEW.reviewed_by IS DISTINCT FROM OLD.reviewed_by
            OR NEW.reviewed_at IS DISTINCT FROM OLD.reviewed_at
       ) THEN
        RAISE EXCEPTION
            'membership mapping review evidence is immutable';
    END IF;

    IF OLD.status = NEW.status THEN
        RETURN NEW;
    END IF;

    IF OLD.status = 'pending_review'
       AND NEW.status NOT IN ('reviewed', 'blocked') THEN
        RAISE EXCEPTION
            'pending membership mapping may only become reviewed or blocked';
    END IF;

    IF OLD.status = 'reviewed'
       AND NEW.status NOT IN ('shadow_read', 'blocked') THEN
        RAISE EXCEPTION
            'reviewed membership mapping may only become shadow_read or blocked';
    END IF;

    IF OLD.status = 'draft_created'
       AND NEW.status NOT IN ('shadow_read', 'blocked') THEN
        RAISE EXCEPTION
            'draft membership mapping may only become shadow_read or blocked';
    END IF;

    IF OLD.status = 'shadow_read'
       AND NEW.status NOT IN ('cutover', 'rolled_back', 'blocked') THEN
        RAISE EXCEPTION
            'shadow membership mapping may only become cutover, rolled_back, or blocked';
    END IF;

    IF OLD.status = 'cutover'
       AND NEW.status <> 'rolled_back' THEN
        RAISE EXCEPTION
            'cutover membership mapping may only become rolled_back';
    END IF;

    IF OLD.status = 'blocked'
       AND OLD.subscription_id IS NULL
       AND NEW.status NOT IN ('pending_review', 'reviewed') THEN
        RAISE EXCEPTION
            'unmigrated blocked mapping requires a new review';
    END IF;

    IF OLD.status = 'blocked'
       AND OLD.subscription_id IS NOT NULL
       AND NEW.status <> 'rolled_back' THEN
        RAISE EXCEPTION
            'migrated blocked mapping must be rolled back';
    END IF;

    IF OLD.status = 'rolled_back' THEN
        RAISE EXCEPTION
            'rolled back membership mapping is terminal';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_catalog_validate_membership_mapping_transition
    ON legacy_membership_subscription_mappings;
CREATE TRIGGER trg_catalog_validate_membership_mapping_transition
BEFORE INSERT OR UPDATE ON legacy_membership_subscription_mappings
FOR EACH ROW
EXECUTE FUNCTION qhpro_catalog_validate_membership_mapping_transition();

COMMENT ON TABLE legacy_membership_subscription_mappings IS
    'Reviewed, idempotent migration state. Legacy memberships remain the runtime source until a separate cutover.';
COMMENT ON COLUMN legacy_membership_subscription_mappings.decision IS
    'Explicit period, renewal, and subscription-key interpretation approved before backfill.';
COMMENT ON COLUMN legacy_membership_subscription_mappings.last_run_id IS
    'Latest explicit operational run that migrated or blocked the mapping.';

COMMIT;
