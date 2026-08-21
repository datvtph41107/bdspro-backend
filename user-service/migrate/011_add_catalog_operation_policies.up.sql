BEGIN;

-- A policy row binds a stable business operation to already-versioned feature
-- and allowance terms. Runtime code must never infer this mapping with a switch.
CREATE TABLE IF NOT EXISTS catalog_plan_operation_policies (
    id BIGSERIAL PRIMARY KEY,
    plan_version_id BIGINT NOT NULL REFERENCES catalog_plan_versions(id),
    operation_code VARCHAR(128) NOT NULL,
    feature_code VARCHAR(128) NOT NULL,
    meter_code VARCHAR(128),
    units_per_action BIGINT NOT NULL DEFAULT 0,
    created_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_catalog_plan_operation_policy
        UNIQUE (plan_version_id, operation_code),
    CONSTRAINT chk_catalog_plan_operation_policy_code
        CHECK (operation_code ~ '^[a-z][a-z0-9_]*([.][a-z0-9_]+)*$'),
    CONSTRAINT chk_catalog_plan_operation_policy_feature
        CHECK (feature_code ~ '^[a-z][a-z0-9_]*([.][a-z0-9_]+)*$'),
    CONSTRAINT chk_catalog_plan_operation_policy_meter
        CHECK (
            meter_code IS NULL
            OR meter_code ~ '^[a-z][a-z0-9_]*([.][a-z0-9_]+)*$'
        ),
    CONSTRAINT chk_catalog_plan_operation_policy_units
        CHECK (
            (meter_code IS NULL AND units_per_action = 0)
            OR (meter_code IS NOT NULL AND units_per_action > 0)
        )
);

CREATE INDEX IF NOT EXISTS idx_catalog_operation_policy_lookup
    ON catalog_plan_operation_policies(plan_version_id, operation_code);

CREATE OR REPLACE FUNCTION qhpro_catalog_validate_operation_policy()
RETURNS TRIGGER AS $$
DECLARE
    parent_published_at TIMESTAMPTZ;
BEGIN
    SELECT published_at
      INTO parent_published_at
      FROM catalog_plan_versions
     WHERE id = NEW.plan_version_id
     FOR SHARE;

    IF parent_published_at IS NULL AND NOT EXISTS (
        SELECT 1
        FROM catalog_plan_versions
        WHERE id = NEW.plan_version_id
    ) THEN
        RAISE EXCEPTION 'operation policy parent plan version % does not exist', NEW.plan_version_id;
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM catalog_plan_entitlements
        WHERE plan_version_id = NEW.plan_version_id
          AND kind = 'feature_access'
          AND feature_code = NEW.feature_code
    ) THEN
        RAISE EXCEPTION
            'operation policy % feature % is not granted by plan version %',
            NEW.operation_code, NEW.feature_code, NEW.plan_version_id;
    END IF;

    IF NEW.meter_code IS NOT NULL AND NOT EXISTS (
        SELECT 1
        FROM catalog_plan_entitlements
        WHERE plan_version_id = NEW.plan_version_id
          AND kind = 'usage_allowance'
          AND meter_code = NEW.meter_code
    ) THEN
        RAISE EXCEPTION
            'operation policy % meter % has no usage allowance in plan version %',
            NEW.operation_code, NEW.meter_code, NEW.plan_version_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_catalog_validate_operation_policy
    ON catalog_plan_operation_policies;
CREATE TRIGGER trg_catalog_validate_operation_policy
BEFORE INSERT OR UPDATE OF plan_version_id, operation_code, feature_code, meter_code, units_per_action
ON catalog_plan_operation_policies
FOR EACH ROW
EXECUTE FUNCTION qhpro_catalog_validate_operation_policy();

-- Reuse the catalog child immutability guard from migration 008.
DROP TRIGGER IF EXISTS trg_catalog_protect_operation_policy
    ON catalog_plan_operation_policies;
CREATE TRIGGER trg_catalog_protect_operation_policy
BEFORE INSERT OR UPDATE OR DELETE ON catalog_plan_operation_policies
FOR EACH ROW
EXECUTE FUNCTION qhpro_catalog_protect_version_child();

COMMENT ON TABLE catalog_plan_operation_policies IS
    'Versioned operation-to-feature/meter/unit binding. Publish freezes it with commercial terms.';

COMMIT;
