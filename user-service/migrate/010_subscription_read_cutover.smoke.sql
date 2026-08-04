DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'legacy_membership_subscription_mappings'
          AND column_name = 'cutover_evidence_checksum'
    ) THEN
        RAISE EXCEPTION 'cutover evidence column is missing';
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger
        WHERE tgname = 'trg_catalog_validate_membership_mapping_transition'
    ) THEN
        RAISE EXCEPTION 'membership mapping transition trigger is missing';
    END IF;
    IF EXISTS (
        SELECT 1
        FROM legacy_membership_subscription_mappings
        WHERE status = 'cutover'
          AND (
              cutover_evidence_checksum !~ '^[a-f0-9]{64}$'
              OR cutover_at IS NULL
              OR cutover_by = ''
          )
    ) THEN
        RAISE EXCEPTION 'cutover mapping is missing evidence';
    END IF;
END;
$$;
