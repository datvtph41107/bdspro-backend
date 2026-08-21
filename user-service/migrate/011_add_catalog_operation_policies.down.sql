BEGIN;

DO $$
BEGIN
    IF current_setting('qhpro.allow_catalog_operation_policy_drop', TRUE) IS DISTINCT FROM 'on' THEN
        RAISE EXCEPTION
            'set qhpro.allow_catalog_operation_policy_drop=on for destructive rollback';
    END IF;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_catalog_protect_operation_policy
    ON catalog_plan_operation_policies;
DROP TRIGGER IF EXISTS trg_catalog_validate_operation_policy
    ON catalog_plan_operation_policies;
DROP FUNCTION IF EXISTS qhpro_catalog_validate_operation_policy();
DROP TABLE IF EXISTS catalog_plan_operation_policies;

COMMIT;
