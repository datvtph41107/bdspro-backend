BEGIN;

DO $$
BEGIN
    IF current_setting('qhpro.allow_catalog_v2_drop', TRUE) IS DISTINCT FROM 'on' THEN
        RAISE EXCEPTION
            'set qhpro.allow_catalog_v2_drop=on for destructive rollback';
    END IF;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_catalog_protect_subscription_event
    ON catalog_subscription_events;
DROP TRIGGER IF EXISTS trg_catalog_validate_subscription_product
    ON catalog_subscriptions;
DROP TRIGGER IF EXISTS trg_catalog_protect_price_item
    ON catalog_price_items;
DROP TRIGGER IF EXISTS trg_catalog_protect_entitlement
    ON catalog_plan_entitlements;
DROP TRIGGER IF EXISTS trg_catalog_protect_plan_version_terms
    ON catalog_plan_versions;
DROP TRIGGER IF EXISTS trg_catalog_validate_plan_version_window
    ON catalog_plan_versions;

DROP FUNCTION IF EXISTS qhpro_catalog_protect_subscription_event();
DROP FUNCTION IF EXISTS qhpro_catalog_validate_subscription_product();
DROP FUNCTION IF EXISTS qhpro_catalog_protect_version_child();
DROP FUNCTION IF EXISTS qhpro_catalog_protect_plan_version_terms();
DROP FUNCTION IF EXISTS qhpro_catalog_validate_plan_version_window();

DROP TABLE IF EXISTS legacy_membership_subscription_mappings;
DROP TABLE IF EXISTS legacy_plan_catalog_mappings;
DROP TABLE IF EXISTS catalog_subscription_events;
DROP TABLE IF EXISTS catalog_subscriptions;
DROP TABLE IF EXISTS catalog_price_items;
DROP TABLE IF EXISTS catalog_plan_entitlements;
DROP TABLE IF EXISTS catalog_plan_versions;
DROP TABLE IF EXISTS catalog_plans;
DROP TABLE IF EXISTS catalog_products;

COMMIT;
