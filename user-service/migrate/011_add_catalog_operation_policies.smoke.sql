DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_name = 'catalog_plan_operation_policies'
    ) THEN
        RAISE EXCEPTION 'catalog operation policy table is missing';
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger
        WHERE tgname = 'trg_catalog_protect_operation_policy'
    ) THEN
        RAISE EXCEPTION 'catalog operation policy immutability trigger is missing';
    END IF;
END;
$$;
