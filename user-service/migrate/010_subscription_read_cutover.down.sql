DO $$
BEGIN
    IF current_setting('qhpro.allow_subscription_cutover_v2_drop', true) <> 'on' THEN
        RAISE EXCEPTION
            'set qhpro.allow_subscription_cutover_v2_drop=on before destructive rollback';
    END IF;
END;
$$;

BEGIN;

DROP INDEX IF EXISTS uq_catalog_subscription_rollback_event;
DROP INDEX IF EXISTS uq_catalog_subscription_cutover_event;
ALTER TABLE catalog_subscription_events
    DROP CONSTRAINT IF EXISTS chk_catalog_subscription_cutover_evidence;

ALTER TABLE legacy_membership_subscription_mappings
    DROP CONSTRAINT IF EXISTS chk_legacy_membership_subscription_mapping_cutover_evidence,
    DROP CONSTRAINT IF EXISTS chk_legacy_membership_subscription_mapping_rollback_evidence,
    DROP COLUMN IF EXISTS cutover_evidence_checksum,
    DROP COLUMN IF EXISTS cutover_at,
    DROP COLUMN IF EXISTS cutover_by,
    DROP COLUMN IF EXISTS rollback_evidence_checksum,
    DROP COLUMN IF EXISTS rolled_back_at,
    DROP COLUMN IF EXISTS rolled_back_by;

CREATE OR REPLACE FUNCTION qhpro_catalog_validate_membership_mapping_transition()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.status NOT IN ('pending_review', 'reviewed', 'blocked') THEN
            RAISE EXCEPTION 'new membership mapping must start pending, reviewed, or blocked';
        END IF;
        IF NEW.subscription_id IS NOT NULL OR NEW.migrated_at IS NOT NULL THEN
            RAISE EXCEPTION 'new membership mapping cannot start migrated';
        END IF;
        RETURN NEW;
    END IF;
    IF OLD.subscription_id IS NOT NULL
       AND NEW.subscription_id IS DISTINCT FROM OLD.subscription_id THEN
        RAISE EXCEPTION 'membership mapping subscription_id is immutable';
    END IF;
    IF OLD.migrated_at IS NOT NULL
       AND NEW.migrated_at IS DISTINCT FROM OLD.migrated_at THEN
        RAISE EXCEPTION 'membership mapping migrated_at is immutable';
    END IF;
    IF (OLD.status IN ('reviewed', 'shadow_read', 'cutover', 'rolled_back')
        OR OLD.subscription_id IS NOT NULL)
       AND (
            NEW.source_checksum IS DISTINCT FROM OLD.source_checksum
            OR NEW.decision IS DISTINCT FROM OLD.decision
            OR NEW.reviewed_by IS DISTINCT FROM OLD.reviewed_by
            OR NEW.reviewed_at IS DISTINCT FROM OLD.reviewed_at
       ) THEN
        RAISE EXCEPTION 'membership mapping review evidence is immutable';
    END IF;
    IF OLD.status = NEW.status THEN RETURN NEW; END IF;
    IF OLD.status = 'pending_review' AND NEW.status NOT IN ('reviewed', 'blocked') THEN
        RAISE EXCEPTION 'pending membership mapping may only become reviewed or blocked';
    END IF;
    IF OLD.status = 'reviewed' AND NEW.status NOT IN ('shadow_read', 'blocked') THEN
        RAISE EXCEPTION 'reviewed membership mapping may only become shadow_read or blocked';
    END IF;
    IF OLD.status = 'draft_created' AND NEW.status NOT IN ('shadow_read', 'blocked') THEN
        RAISE EXCEPTION 'draft membership mapping may only become shadow_read or blocked';
    END IF;
    IF OLD.status = 'shadow_read' AND NEW.status NOT IN ('cutover', 'rolled_back', 'blocked') THEN
        RAISE EXCEPTION 'shadow membership mapping transition is invalid';
    END IF;
    IF OLD.status = 'cutover' AND NEW.status <> 'rolled_back' THEN
        RAISE EXCEPTION 'cutover membership mapping may only become rolled_back';
    END IF;
    IF OLD.status = 'blocked' AND OLD.subscription_id IS NULL
       AND NEW.status NOT IN ('pending_review', 'reviewed') THEN
        RAISE EXCEPTION 'unmigrated blocked mapping requires a new review';
    END IF;
    IF OLD.status = 'blocked' AND OLD.subscription_id IS NOT NULL
       AND NEW.status <> 'rolled_back' THEN
        RAISE EXCEPTION 'migrated blocked mapping must be rolled back';
    END IF;
    IF OLD.status = 'rolled_back' THEN
        RAISE EXCEPTION 'rolled back membership mapping is terminal';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMIT;
