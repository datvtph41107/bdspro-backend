DO $$
BEGIN
    IF current_setting(
        'qhpro.allow_subscription_backfill_v2_rollback',
        true
    ) IS DISTINCT FROM 'on' THEN
        RAISE EXCEPTION
            'set qhpro.allow_subscription_backfill_v2_rollback=on before rollback';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM legacy_membership_subscription_mappings
        WHERE status IN ('reviewed', 'shadow_read', 'cutover')
           OR subscription_id IS NOT NULL
    ) THEN
        RAISE EXCEPTION
            'subscription backfill mappings must be rolled back and reconciled first';
    END IF;
END;
$$;

BEGIN;

DROP TRIGGER IF EXISTS trg_catalog_validate_membership_mapping_transition
    ON legacy_membership_subscription_mappings;
DROP FUNCTION IF EXISTS qhpro_catalog_validate_membership_mapping_transition();
DROP INDEX IF EXISTS uq_catalog_subscription_legacy_backfill_event;
ALTER TABLE catalog_subscription_events
    DROP CONSTRAINT IF EXISTS chk_catalog_subscription_legacy_backfill_evidence;

ALTER TABLE legacy_membership_subscription_mappings
    DROP CONSTRAINT IF EXISTS chk_legacy_membership_subscription_mapping_run_id,
    DROP CONSTRAINT IF EXISTS chk_legacy_membership_subscription_mapping_blocked,
    DROP CONSTRAINT IF EXISTS chk_legacy_membership_subscription_mapping_migrated,
    DROP CONSTRAINT IF EXISTS chk_legacy_membership_subscription_mapping_review,
    DROP CONSTRAINT IF EXISTS chk_legacy_membership_subscription_mapping_decision,
    DROP CONSTRAINT IF EXISTS chk_legacy_membership_subscription_mapping_status;

ALTER TABLE legacy_membership_subscription_mappings
    ADD CONSTRAINT chk_legacy_membership_subscription_mapping_status
        CHECK (
            status IN (
                'pending_review',
                'draft_created',
                'shadow_read',
                'cutover',
                'rolled_back',
                'blocked'
            )
        );

ALTER TABLE legacy_membership_subscription_mappings
    DROP COLUMN IF EXISTS last_error,
    DROP COLUMN IF EXISTS last_run_id,
    DROP COLUMN IF EXISTS blocked_reason,
    DROP COLUMN IF EXISTS reviewed_at,
    DROP COLUMN IF EXISTS reviewed_by,
    DROP COLUMN IF EXISTS decision;

COMMIT;
