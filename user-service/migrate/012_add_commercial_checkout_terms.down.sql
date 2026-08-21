BEGIN;

DROP TABLE IF EXISTS commercial_checkout_commands;

DROP TRIGGER IF EXISTS trg_catalog_protect_subscription_term
    ON catalog_plan_versions;
DROP FUNCTION IF EXISTS qhpro_catalog_protect_subscription_term();

DROP INDEX IF EXISTS idx_catalog_subscriptions_pending_plan;
ALTER TABLE catalog_subscriptions
    DROP CONSTRAINT IF EXISTS chk_catalog_subscriptions_pending_change;
ALTER TABLE catalog_subscriptions
    DROP COLUMN IF EXISTS pending_effective_at;
ALTER TABLE catalog_subscriptions
    DROP COLUMN IF EXISTS pending_plan_version_id;

ALTER TABLE catalog_plan_versions
    DROP CONSTRAINT IF EXISTS chk_catalog_plan_versions_subscription_term;
ALTER TABLE catalog_plan_versions
    DROP COLUMN IF EXISTS subscription_term_days;

ALTER TABLE catalog_plans
    DROP CONSTRAINT IF EXISTS chk_catalog_plans_tier_rank;
ALTER TABLE catalog_plans
    DROP COLUMN IF EXISTS tier_rank;

COMMIT;
