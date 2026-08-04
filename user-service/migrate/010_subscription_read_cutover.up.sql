BEGIN;

ALTER TABLE legacy_membership_subscription_mappings
    ADD COLUMN IF NOT EXISTS cutover_evidence_checksum CHAR(64) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS cutover_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS cutover_by VARCHAR(128) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS rollback_evidence_checksum CHAR(64) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS rolled_back_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS rolled_back_by VARCHAR(128) NOT NULL DEFAULT '';

UPDATE catalog_subscriptions AS subscription
SET
    created_at = membership.created_at,
    updated_at = membership.updated_at
FROM legacy_membership_subscription_mappings AS mapping
JOIN memberships AS membership
    ON membership.id = mapping.legacy_membership_id
WHERE mapping.subscription_id = subscription.id
  AND mapping.status IN ('shadow_read', 'cutover')
  AND membership.created_at IS NOT NULL
  AND membership.updated_at IS NOT NULL
  AND (
      subscription.created_at IS DISTINCT FROM membership.created_at
      OR subscription.updated_at IS DISTINCT FROM membership.updated_at
  );

ALTER TABLE legacy_membership_subscription_mappings
    DROP CONSTRAINT IF EXISTS chk_legacy_membership_subscription_mapping_cutover_evidence,
    DROP CONSTRAINT IF EXISTS chk_legacy_membership_subscription_mapping_rollback_evidence;

ALTER TABLE legacy_membership_subscription_mappings
    ADD CONSTRAINT chk_legacy_membership_subscription_mapping_cutover_evidence
        CHECK (
            status NOT IN ('cutover', 'rolled_back')
            OR (
                cutover_evidence_checksum ~ '^[a-f0-9]{64}$'
                AND cutover_at IS NOT NULL
                AND cutover_by ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$'
            )
        ) NOT VALID,
    ADD CONSTRAINT chk_legacy_membership_subscription_mapping_rollback_evidence
        CHECK (
            status <> 'rolled_back'
            OR (
                rollback_evidence_checksum ~ '^[a-f0-9]{64}$'
                AND rolled_back_at IS NOT NULL
                AND rolled_back_by ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$'
            )
        ) NOT VALID;

ALTER TABLE catalog_subscription_events
    DROP CONSTRAINT IF EXISTS chk_catalog_subscription_cutover_evidence;
ALTER TABLE catalog_subscription_events
    ADD CONSTRAINT chk_catalog_subscription_cutover_evidence
        CHECK (
            event_type NOT IN (
                'migration.subscription_read_cutover_activated',
                'migration.subscription_read_cutover_rolled_back'
            )
            OR (
                jsonb_typeof(metadata) = 'object'
                AND metadata ? 'legacyMembershipId'
                AND metadata ? 'subscriptionId'
                AND metadata ? 'sourceChecksum'
                AND metadata ? 'planTermsChecksum'
                AND metadata ? 'evidenceChecksum'
                AND metadata ? 'populationChecksum'
                AND metadata ? 'batchIndex'
                AND metadata ? 'batchCount'
                AND metadata ->> 'legacyMembershipId' ~ '^[1-9][0-9]*$'
                AND metadata ->> 'subscriptionId' ~ '^[1-9][0-9]*$'
                AND metadata ->> 'sourceChecksum' ~ '^[a-f0-9]{64}$'
                AND metadata ->> 'planTermsChecksum' ~ '^[a-f0-9]{64}$'
                AND metadata ->> 'evidenceChecksum' ~ '^[a-f0-9]{64}$'
                AND metadata ->> 'populationChecksum' ~ '^[a-f0-9]{64}$'
                AND metadata ->> 'batchIndex' ~ '^[1-9][0-9]*$'
                AND metadata ->> 'batchCount' ~ '^[1-9][0-9]*$'
                AND (metadata ->> 'batchIndex')::integer <= (metadata ->> 'batchCount')::integer
            )
        ) NOT VALID;

CREATE UNIQUE INDEX IF NOT EXISTS uq_catalog_subscription_cutover_event
    ON catalog_subscription_events(subscription_id, event_type)
    WHERE event_type = 'migration.subscription_read_cutover_activated';
CREATE UNIQUE INDEX IF NOT EXISTS uq_catalog_subscription_rollback_event
    ON catalog_subscription_events(subscription_id, event_type)
    WHERE event_type = 'migration.subscription_read_cutover_rolled_back';

CREATE OR REPLACE FUNCTION qhpro_catalog_validate_membership_mapping_transition()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.status NOT IN ('pending_review', 'reviewed', 'blocked') THEN
            RAISE EXCEPTION
                'new membership mapping must start pending, reviewed, or blocked';
        END IF;
        IF NEW.subscription_id IS NOT NULL OR NEW.migrated_at IS NOT NULL
           OR NEW.cutover_at IS NOT NULL OR NEW.rolled_back_at IS NOT NULL THEN
            RAISE EXCEPTION
                'new membership mapping cannot start migrated or cut over';
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
    IF OLD.cutover_evidence_checksum <> ''
       AND NEW.cutover_evidence_checksum IS DISTINCT FROM OLD.cutover_evidence_checksum THEN
        RAISE EXCEPTION 'membership mapping cutover evidence is immutable';
    END IF;
    IF OLD.cutover_at IS NOT NULL
       AND NEW.cutover_at IS DISTINCT FROM OLD.cutover_at THEN
        RAISE EXCEPTION 'membership mapping cutover_at is immutable';
    END IF;
    IF OLD.cutover_by <> '' AND NEW.cutover_by IS DISTINCT FROM OLD.cutover_by THEN
        RAISE EXCEPTION 'membership mapping cutover_by is immutable';
    END IF;
    IF OLD.rollback_evidence_checksum <> ''
       AND NEW.rollback_evidence_checksum IS DISTINCT FROM OLD.rollback_evidence_checksum THEN
        RAISE EXCEPTION 'membership mapping rollback evidence is immutable';
    END IF;
    IF OLD.rolled_back_at IS NOT NULL
       AND NEW.rolled_back_at IS DISTINCT FROM OLD.rolled_back_at THEN
        RAISE EXCEPTION 'membership mapping rolled_back_at is immutable';
    END IF;
    IF OLD.rolled_back_by <> ''
       AND NEW.rolled_back_by IS DISTINCT FROM OLD.rolled_back_by THEN
        RAISE EXCEPTION 'membership mapping rolled_back_by is immutable';
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

    IF OLD.status = NEW.status THEN
        RETURN NEW;
    END IF;
    IF OLD.status = 'pending_review'
       AND NEW.status NOT IN ('reviewed', 'blocked') THEN
        RAISE EXCEPTION 'pending membership mapping may only become reviewed or blocked';
    END IF;
    IF OLD.status = 'reviewed'
       AND NEW.status NOT IN ('shadow_read', 'blocked') THEN
        RAISE EXCEPTION 'reviewed membership mapping may only become shadow_read or blocked';
    END IF;
    IF OLD.status = 'draft_created'
       AND NEW.status NOT IN ('shadow_read', 'blocked') THEN
        RAISE EXCEPTION 'draft membership mapping may only become shadow_read or blocked';
    END IF;
    IF OLD.status = 'shadow_read'
       AND NEW.status NOT IN ('cutover', 'rolled_back', 'blocked') THEN
        RAISE EXCEPTION 'shadow membership mapping may only become cutover, rolled_back, or blocked';
    END IF;
    IF OLD.status = 'shadow_read' AND NEW.status = 'cutover'
       AND (
            NEW.cutover_evidence_checksum !~ '^[a-f0-9]{64}$'
            OR NEW.cutover_at IS NULL
            OR NEW.cutover_by = ''
       ) THEN
        RAISE EXCEPTION 'cutover requires immutable evidence';
    END IF;
    IF OLD.status = 'cutover' AND NEW.status <> 'rolled_back' THEN
        RAISE EXCEPTION 'cutover membership mapping may only become rolled_back';
    END IF;
    IF OLD.status = 'cutover' AND NEW.status = 'rolled_back'
       AND (
            NEW.rollback_evidence_checksum !~ '^[a-f0-9]{64}$'
            OR NEW.rolled_back_at IS NULL
            OR NEW.rolled_back_by = ''
       ) THEN
        RAISE EXCEPTION 'rollback requires immutable evidence';
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

COMMENT ON COLUMN legacy_membership_subscription_mappings.cutover_evidence_checksum IS
    'Checksum of the sorted concrete membership manifest approved for catalog reads.';
COMMENT ON COLUMN legacy_membership_subscription_mappings.rollback_evidence_checksum IS
    'Checksum of the explicit rollback manifest; feature-flag rollback remains immediate.';

COMMIT;
