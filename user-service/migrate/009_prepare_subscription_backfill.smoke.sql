DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = current_schema()
          AND table_name = 'legacy_membership_subscription_mappings'
          AND column_name = 'decision'
    ) THEN
        RAISE EXCEPTION 'membership mapping decision column is missing';
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_indexes
        WHERE schemaname = current_schema()
          AND indexname = 'uq_catalog_subscription_legacy_backfill_event'
    ) THEN
        RAISE EXCEPTION 'legacy backfill event index is missing';
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'chk_catalog_subscription_legacy_backfill_evidence'
    ) THEN
        RAISE EXCEPTION 'legacy backfill event evidence constraint is missing';
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_trigger
        WHERE tgname = 'trg_catalog_validate_membership_mapping_transition'
          AND NOT tgisinternal
    ) THEN
        RAISE EXCEPTION 'membership mapping transition trigger is missing';
    END IF;
END;
$$;

BEGIN;

DO $$
DECLARE
    insert_rejected BOOLEAN := FALSE;
BEGIN
    BEGIN
        INSERT INTO legacy_membership_subscription_mappings (
            legacy_membership_id,
            subscription_id,
            source_checksum,
            status,
            decision,
            reviewed_by,
            reviewed_at,
            migrated_at,
            last_run_id
        ) VALUES (
            9223372036854769999,
            9223372036854769999,
            REPEAT('a', 64),
            'cutover',
            '{"reason":"invalid direct cutover"}'::jsonb,
            1,
            NOW(),
            NOW(),
            'run_invalid'
        );
    EXCEPTION WHEN OTHERS THEN
        insert_rejected := TRUE;
    END;

    IF NOT insert_rejected THEN
        RAISE EXCEPTION 'direct cutover mapping insert was not rejected';
    END IF;
END;
$$;

DELETE FROM legacy_membership_subscription_mappings
WHERE legacy_membership_id = 9223372036854770000;

INSERT INTO legacy_membership_subscription_mappings (
    legacy_membership_id,
    source_checksum,
    status
) VALUES (
    9223372036854770000,
    REPEAT('a', 64),
    'pending_review'
);

UPDATE legacy_membership_subscription_mappings
SET
    status = 'reviewed',
    decision = '{
        "periodHandling":"paid_at_to_finished_date",
        "subscriptionKeyStrategy":"legacy_membership_id",
        "autoRenew":false,
        "reason":"smoke test"
    }'::jsonb,
    reviewed_by = 1,
    reviewed_at = NOW()
WHERE legacy_membership_id = 9223372036854770000;

DO $$
DECLARE
    transition_rejected BOOLEAN := FALSE;
BEGIN
    BEGIN
        UPDATE legacy_membership_subscription_mappings
        SET status = 'cutover'
        WHERE legacy_membership_id = 9223372036854770000;
    EXCEPTION WHEN OTHERS THEN
        transition_rejected := TRUE;
    END;

    IF NOT transition_rejected THEN
        RAISE EXCEPTION 'reviewed to cutover transition was not rejected';
    END IF;
END;
$$;

ROLLBACK;
